// Package http_router предоставляет функции для создания и настройки HTTP роутера.
//
// Основные возможности:
// - Создание роутера на базе chi с предустановленными middleware
// - Встроенные middleware для recovery, логирования, таймаутов
// - Автоматическое подключение health endpoints
//
// Пример использования:
//
//	router, stopLimiter := http_router.NewRouter(ctx, timeout, bucketBoundaries)
//	defer stopLimiter()
//	router.Get("/api/health", healthHandler)
//	server := http_server.NewServer(router, address, ...)
package http_router

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http/middleware"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metric"
)

// NewRouter создаёт роутер и настраивает middleware в порядке:
// 1. Идентификация (RequestID, RealIP)
// 2. Быстрая защита (security headers — reverse proxy; RequestFirewall, IP RateLimit) — до тяжёлых логеров
// 3. Стабильность (Recoverer, Timeout)
// 4. Наблюдаемость (Logging, метрики)
// 5. CORS (если настроено)
// Возвращает роутер и stopFunc для graceful shutdown (остановка IP rate limiter).
func NewRouter(
	ctx context.Context,
	timeout time.Duration,
	allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders []string,
	allowCredentials bool,
	maxAge int,
	httpBucketBoundaries []float64,
) (*chi.Mux, func()) {
	r := chi.NewRouter()

	// 1. Идентификация
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	// 2. Быстрая защита (до логеров и метрик); security headers — reverse proxy
	r.Use(middleware.RequestFirewallMiddleware)
	ipRateLimitMw, stopIPLimiter := middleware.RateLimitMiddleware()
	r.Use(ipRateLimitMw)

	// 3. Стабильность
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(timeout))

	// 4. Наблюдаемость
	r.Use(middleware.LoggingMiddleware)
	r.Use(metric.HTTPMiddleware(ctx, httpBucketBoundaries))

	if len(allowedOrigins) > 0 {
		r.Use(middleware.CORSMiddleware(
			allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders,
			allowCredentials, maxAge,
		))
	}

	return r, stopIPLimiter
}
