package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

type userLimiterEntry struct {
	count    int
	windowTo time.Time
	lastSeen time.Time
}

type userRateLimiter struct {
	mu      sync.Mutex
	entries map[uuid.UUID]*userLimiterEntry
	cancel  context.CancelFunc
}

func newUserRateLimiter(ctx context.Context) *userRateLimiter {
	ctx, cancel := context.WithCancel(ctx)
	l := &userRateLimiter{
		entries: make(map[uuid.UUID]*userLimiterEntry),
		cancel:  cancel,
	}
	go l.cleanupLoop(ctx)
	return l
}

func (l *userRateLimiter) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(pkghttp.UserLimiterCleanupInt)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.cleanup()
		}
	}
}

func (l *userRateLimiter) cleanup() {
	cutoff := time.Now().Add(-pkghttp.RateLimitWindow5Min)
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, v := range l.entries {
		if v.lastSeen.Before(cutoff) {
			delete(l.entries, k)
		}
	}
}

// Stop останавливает фоновую очистку записей. Безопасно вызывать несколько раз.
func (l *userRateLimiter) Stop() {
	l.cancel()
}

// allow проверяет лимит для пользователя. При превышении возвращает (false, retryAfter).
// retryAfter — время до сброса окна (для заголовка Retry-After).
func (l *userRateLimiter) allow(userID uuid.UUID) (allowed bool, retryAfter time.Duration) {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if e, ok := l.entries[userID]; ok {
		e.lastSeen = now
		if now.After(e.windowTo) {
			e.count = 1
			e.windowTo = now.Add(pkghttp.UserRateLimitWindow)
			return true, 0
		}
		if e.count >= pkghttp.UserRateLimitMax {
			d := max(0, time.Until(e.windowTo))
			return false, d
		}
		e.count++
		return true, 0
	}

	// Fail-open: при переполнении карты пропускаем запрос, чтобы не блокировать новых пользователей.
	if len(l.entries) >= pkghttp.RateLimiterMaxSize {
		return true, 0
	}

	l.entries[userID] = &userLimiterEntry{
		count:    1,
		windowTo: now.Add(pkghttp.UserRateLimitWindow),
		lastSeen: now,
	}
	return true, 0
}

// UserRateLimitMiddleware ограничивает количество запросов на аутентифицированного пользователя.
// Применяется после проверки аутентификации (JWT или сессия) — использует user_id из контекста.
// Неаутентифицированные запросы пропускаются (их покрывает IP rate limiter).
// Возвращает middleware и функцию stop для graceful shutdown.
func UserRateLimitMiddleware(ctx context.Context) (func(http.Handler) http.Handler, func()) {
	limiter := newUserRateLimiter(ctx)

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := metadata.UserID(r.Context())
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			allowed, retryAfter := limiter.allow(userID)
			if !allowed {
				seconds := int(retryAfter.Seconds())
				if seconds < pkghttp.RateLimitRetryAfterSeconds {
					seconds = pkghttp.RateLimitRetryAfterSeconds
				}
				w.Header().Set("Retry-After", strconv.Itoa(seconds))
				pkghttp.WriteJSON(r.Context(), w, http.StatusTooManyRequests, pkghttp.ErrorBody{
					Code:    http.StatusTooManyRequests,
					Message: "Превышен лимит запросов. Повторите позже",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return mw, limiter.Stop
}
