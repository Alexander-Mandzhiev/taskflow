package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// LoggingMiddleware создает HTTP middleware для логирования запросов.
// Логирует только ошибки и медленные запросы (>500ms) — успешные быстрые запросы покрыты метриками.
const slowRequestThreshold = 500 * time.Millisecond

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := middleware.GetReqID(ctx)
		startTime := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, r)
		duration := time.Since(startTime)

		// Логируем только ошибки и медленные запросы; успешные быстрые — без лога (метрики).
		switch {
		case sr.status >= 500:
			logger.Error(ctx, "[HTTP] Server error",
				zap.String("req_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		case sr.status >= 400:
			if sr.status == 404 && (strings.Contains(r.URL.Path, "/statistic/") || strings.Contains(r.URL.Path, "/availability/")) {
				logger.Debug(ctx, "[HTTP] Resource not found (statistic/availability)",
					zap.String("req_id", reqID),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", sr.status),
					zap.Duration("duration", duration),
				)
			} else {
				logger.Warn(ctx, "[HTTP] Client error",
					zap.String("req_id", reqID),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", sr.status),
					zap.Duration("duration", duration),
				)
			}
		case duration > slowRequestThreshold:
			logger.Info(ctx, "[HTTP] Slow request",
				zap.String("req_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		}
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
