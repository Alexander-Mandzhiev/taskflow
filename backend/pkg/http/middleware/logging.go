package middleware

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"mkk/pkg/logger"
)

// LoggingMiddleware создает HTTP middleware для логирования запросов.
// Логирует только ошибки и медленные запросы (>500ms) — успешные быстрые запросы покрыты метриками.
const slowRequestThreshold = 500 * time.Millisecond

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startTime := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, r)
		duration := time.Since(startTime)

		// Логируем только ошибки и медленные запросы
		switch {
		case sr.status >= 500:
			logger.Error(ctx, "[HTTP] Server error",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		case sr.status >= 400:
			// 4xx — ошибка клиента
			// 404 для статистики — это нормально (сущность может быть удалена), логируем как DEBUG
			if sr.status == 404 && (strings.Contains(r.URL.Path, "/statistic/") || strings.Contains(r.URL.Path, "/availability/")) {
				logger.Debug(ctx, "[HTTP] Resource not found (statistic/availability)",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", sr.status),
					zap.Duration("duration", duration),
				)
			} else {
				// Другие 4xx — логируем как WARN
				logger.Warn(ctx, "[HTTP] Client error",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", sr.status),
					zap.Duration("duration", duration),
				)
			}
		case duration > slowRequestThreshold:
			// Медленные успешные запросы — логируем для анализа производительности
			logger.Info(ctx, "[HTTP] Slow request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		}
		// Успешные быстрые запросы не логируем — покрыто метриками
	})
}

// statusRecorder оборачивает http.ResponseWriter для отслеживания статуса и размера ответа
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	n, err := sr.ResponseWriter.Write(b)
	sr.bytes += n
	return n, err
}
