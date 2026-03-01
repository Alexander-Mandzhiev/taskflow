package middleware

import (
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	pkghttp "mkk/pkg/http"
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
		// Нормализуем путь: //api//v1 -> /api/v1 (API терпимее к мелким ошибкам клиентов)
		p = path.Clean(p)
		if p != r.URL.Path {
			r.URL.Path = p
		}

		// Быстрые блоки на самые частые сканер-пути
		lp := strings.ToLower(p)
		// На backend не разрешаем никаких dot-paths (включая /.well-known/*) — ACME пусть обслуживает Traefik.
		if strings.HasPrefix(lp, "/.") {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(lp, "..") || strings.Contains(lp, "%2e") {
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
	l := &ipRateLimiter{
		entries:  make(map[string]*ipLimiterEntry),
		stopCh:   make(chan struct{}),
		ttl:      ttl,
		interval: interval,
		maxSize:  10000,
	}
	l.startCleanup()
	return l
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

// Чувствительные endpoints, требующие строгого rate limit (только для POST).
// Пути должны совпадать с роутером: routes.RegisterAPIs → /api/v1 + public.Register → /register, /login.
var sensitiveEndpoints = map[string]struct{}{
	"/api/v1/login":           {},
	"/api/v1/register":        {},
	"/api/v1/forgot-password": {},
	"/api/v1/reset-password":  {},
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

			ip := pkghttp.ClientIP(r)

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
