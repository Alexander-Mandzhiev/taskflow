// Package http_router предоставляет функции для создания и настройки HTTP роутера.
//
// Основные возможности:
// - Создание роутера на базе chi с предустановленными middleware
// - Встроенные middleware для recovery, логирования, таймаутов
// - Автоматическое подключение health endpoints
//
// Пример использования:
//
//	router := http_router.NewRouter(ctx, timeout, bucketBoundaries)
//	router.Get("/api/health", healthHandler)
//	server := http_server.NewServer(router, address, ...)
package http_router

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"mkk/pkg/http/middleware"
	"mkk/pkg/metric"
)

// NewRouter создаёт роутер и настраивает его с базовыми middleware.
// Автоматически подключает стандартные middleware в оптимальном порядке:
// 1. RequestID - добавление уникального ID для каждого запроса
// 2. RealIP - извлечение реального IP адреса клиента
// 3. Recovery - обработка паник
// 4. Timeout - ограничение времени выполнения запроса
// 5. LoggingMiddleware - логирование входящих запросов
// 6. HTTPMiddleware - сбор метрик запросов, ответов и времени выполнения
// 7. CORS - настройка политик межсайтовых запросов (если настроено)
func NewRouter(
	ctx context.Context,
	timeout time.Duration,
	allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders []string,
	allowCredentials bool,
	maxAge int,
	httpBucketBoundaries []float64,
) *chi.Mux {
	r := chi.NewRouter()

	// Базовые middleware от chi
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(timeout))

	// Логирование и метрики
	r.Use(middleware.LoggingMiddleware)
	r.Use(metric.HTTPMiddleware(ctx, httpBucketBoundaries))

	// CORS middleware (если настроено)
	if len(allowedOrigins) > 0 {
		r.Use(middleware.CORSMiddleware(
			allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders,
			allowCredentials, maxAge,
		))
	}

	return r
}
