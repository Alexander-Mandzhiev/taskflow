package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SecurityHeadersMiddleware добавляет базовые security headers.
// KISS: безопасные значения по умолчанию, без жёсткой CSP (чтобы не ломать фронт).
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Не даём браузерам пытаться угадать тип контента.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Запрещаем встраивание в iframe (защита от clickjacking).
		w.Header().Set("X-Frame-Options", "DENY")
		// Минимизируем утечки реферера.
		w.Header().Set("Referrer-Policy", "no-referrer")
		// Отключаем “фичи” браузера, которыми API обычно не пользуется.
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS только если запрос пришёл по HTTPS (или прокси сообщает об этом).
		// Иначе можно случайно “прибить” доступ к dev по http.
		if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		next.ServeHTTP(w, r)
	})
}

// RequestFirewallMiddleware делает KISS-защиту от типичных scanner-запросов:
// - всё, что не похоже на API/health, отвечает 404 (быстро)
// - отдельные “опасные” паттерны (/.env, /.git, wp-*, etc) тоже 404
func RequestFirewallMiddleware(next http.Handler) http.Handler {
	allowedExact := map[string]struct{}{
		"/health":  {},
		"/healthz": {},
		"/live":    {},
		"/ready":   {},
		"/start":   {},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "" {
			p = "/"
		}

		// Быстрые блоки на самые частые сканер-пути
		lp := strings.ToLower(p)
		// На backend не разрешаем никаких dot-paths (включая /.well-known/*) — ACME пусть обслуживает Traefik.
		// Если когда-нибудь понадобится — добавим явный allowlist.
		if strings.HasPrefix(lp, "/.") {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(lp, "..") || strings.Contains(lp, "%2e") {
			http.NotFound(w, r)
			return
		}
		// Блокируем попытки обхода роутера через двойные слэши (//api -> /api)
		if strings.Contains(p, "//") {
			http.NotFound(w, r)
			return
		}
		// Паттерны сканеров (без дублей - /.git, /.env уже заблокированы выше через /.)
		if hasAnyPrefix(lp,
			"/wp-", "/wordpress", "/xmlrpc.php",
			"/_next", "/_nuxt", "/_astro", "/@vite",
			"/actuator", "/server-status", "/phpinfo",
		) {
			http.NotFound(w, r)
			return
		}

		// Разрешаем health эндпоинты
		if _, ok := allowedExact[p]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// Разрешаем только API.
		// В текущем приложении все роуты живут под /api/v1.
		if strings.HasPrefix(p, "/api/") || p == "/api" {
			next.ServeHTTP(w, r)
			return
		}

		http.NotFound(w, r)
	})
}

func hasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

type ipLimiterEntry struct {
	count    int
	windowTo time.Time
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu       sync.Mutex
	entries  map[string]*ipLimiterEntry
	once     sync.Once
	stopCh   chan struct{}
	ttl      time.Duration
	interval time.Duration
	maxSize  int
}

func newIPRateLimiter(ttl, interval time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		entries:  make(map[string]*ipLimiterEntry),
		stopCh:   make(chan struct{}),
		ttl:      ttl,
		interval: interval,
		maxSize:  10000,
	}
}

func (l *ipRateLimiter) startCleanup() {
	l.once.Do(func() {
		go func() {
			t := time.NewTicker(l.interval)
			defer t.Stop()
			for {
				select {
				case <-l.stopCh:
					return
				case <-t.C:
					cutoff := time.Now().Add(-l.ttl)
					l.mu.Lock()
					for k, v := range l.entries {
						if v.lastSeen.Before(cutoff) {
							delete(l.entries, k)
						}
					}
					l.mu.Unlock()
				}
			}
		}()
	})
}

// Stop завершает фоновую горутину очистки.
func (l *ipRateLimiter) Stop() {
	select {
	case <-l.stopCh:
	default:
		close(l.stopCh)
	}
}

func (l *ipRateLimiter) allow(ip string, maxPerWindow int, window time.Duration) bool {
	l.startCleanup()
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if e, ok := l.entries[ip]; ok {
		e.lastSeen = now
		if now.After(e.windowTo) {
			e.count = 0
			e.windowTo = now.Add(window)
		}
		if e.count >= maxPerWindow {
			return false
		}
		e.count++
		return true
	}

	// Очистка только при переполнении — O(N) не на каждый запрос, а только при достижении лимита.
	if len(l.entries) >= l.maxSize {
		cutoff := now.Add(-l.ttl)
		for k, v := range l.entries {
			if v.lastSeen.Before(cutoff) {
				delete(l.entries, k)
			}
		}
		if len(l.entries) >= l.maxSize {
			return false
		}
	}

	l.entries[ip] = &ipLimiterEntry{
		count:    1,
		windowTo: now.Add(window),
		lastSeen: now,
	}
	return true
}

// Чувствительные endpoints, требующие строгого rate limit (только для POST)
var sensitiveEndpoints = map[string]struct{}{
	"/api/v1/session/login":    {},
	"/api/v1/session/register": {},
	"/api/v1/forgot-password":  {},
	"/api/v1/reset-password":   {},
}

// RateLimitMiddleware — простой in-memory rate limit по IP.
// Отдельно ужесточаем лимиты на чувствительные endpoints (login, register, password reset).
// Возвращает middleware и функцию stop для graceful shutdown (остановка фоновой очистки).
func RateLimitMiddleware() (func(http.Handler) http.Handler, func()) {
	limiterStore := newIPRateLimiter(15*time.Minute, 5*time.Minute)

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/health", "/healthz", "/live", "/ready", "/start":
				next.ServeHTTP(w, r)
				return
			}

			ip := clientIP(r)

			maxPerSecond := 30

			if r.Method == http.MethodPost {
				if _, ok := sensitiveEndpoints[r.URL.Path]; ok {
					maxPerSecond = 7
				}
			}

			if !limiterStore.allow(ip, maxPerSecond, time.Second) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return mw, limiterStore.Stop
}

func clientIP(r *http.Request) string {
	// После chi/middleware.RealIP обычно уже нормализован RemoteAddr,
	// но на всякий случай парсим host:port.
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	if ip := net.ParseIP(strings.TrimSpace(r.RemoteAddr)); ip != nil {
		return ip.String()
	}
	// Последний шанс — X-Real-IP (могут выставлять прокси).
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); net.ParseIP(xrip) != nil {
		return xrip
	}
	return r.RemoteAddr
}
