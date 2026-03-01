package metric

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	doubleSlashRe = regexp.MustCompile(`/+`)
	uuidSegmentRe = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	intSegmentRe  = regexp.MustCompile(`^\d+$`)
	tempSegmentRe = regexp.MustCompile(`(?i)^temp-\d+$`)
	hex24Segment  = regexp.MustCompile(`(?i)^[0-9a-f]{24}$`)

	// ErrNotHijacker возвращается из Hijack(), если обёрнутый ResponseWriter не реализует http.Hijacker.
	ErrNotHijacker = errors.New("response writer does not support hijacking")
	// ErrNotPusher возвращается из Push(), если обёрнутый ResponseWriter не реализует http.Pusher.
	ErrNotPusher = errors.New("response writer does not support push")
)

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Схлопываем повторяющиеся слэши одной операцией.
	path = doubleSlashRe.ReplaceAllString(path, "/")

	// Убираем trailing slash (кроме корня).
	if len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}

	parts := strings.Split(path, "/")
	for i := range parts {
		seg := parts[i]
		if seg == "" {
			continue
		}

		switch {
		case hex24Segment.MatchString(seg):
			// Проверяем hex24 ПЕРЕД int, так как hex24 может быть числом
			parts[i] = "{id}"
		case uuidSegmentRe.MatchString(seg):
			parts[i] = "{id}"
		case tempSegmentRe.MatchString(seg):
			parts[i] = "temp-{id}"
		case intSegmentRe.MatchString(seg):
			// Самое общее - в конец
			parts[i] = "{id}"
		}
	}

	out := strings.Join(parts, "/")
	if out == "" {
		return "/"
	}
	if !strings.HasPrefix(out, "/") {
		out = "/" + out
	}
	return out
}

func shouldRecordMetrics(path string) bool {
	if path == "" {
		return false
	}
	// Метрики только для API и health, чтобы не шуметь на сканерах.
	if strings.HasPrefix(path, "/api") {
		return true
	}
	switch path {
	case "/health", "/healthz", "/live", "/ready", "/start":
		return true
	default:
		return false
	}
}

type httpInstruments struct {
	requestCounter        metric.Int64Counter
	responseCounter       metric.Int64Counter
	responseTimeHistogram metric.Float64Histogram
}

func noOpMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { return next }
}

func (m *Metrics) createHTTPInstruments(meter metric.Meter, bucketBoundaries []float64) (*httpInstruments, error) {
	// Вычисляем имена метрик один раз при инициализации
	requestCounterName := m.getMetricName("http_requests_total")
	responseCounterName := m.getMetricName("http_responses_total")
	responseTimeHistName := m.getMetricName("http_response_time_seconds")

	requestCounter, err := meter.Int64Counter(
		requestCounterName,
		metric.WithDescription("Количество HTTP запросов"),
	)
	if err != nil {
		return nil, fmt.Errorf("create request counter: %w", err)
	}

	responseCounter, err := meter.Int64Counter(
		responseCounterName,
		metric.WithDescription("Количество HTTP ответов"),
	)
	if err != nil {
		return nil, fmt.Errorf("create response counter: %w", err)
	}

	responseTimeHistogram, err := meter.Float64Histogram(
		responseTimeHistName,
		metric.WithDescription("Время выполнения HTTP запросов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(bucketBoundaries...),
	)
	if err != nil {
		return nil, fmt.Errorf("create response time histogram: %w", err)
	}

	return &httpInstruments{
		requestCounter:        requestCounter,
		responseCounter:       responseCounter,
		responseTimeHistogram: responseTimeHistogram,
	}, nil
}

func (m *Metrics) httpMiddlewareHandler(instruments *httpInstruments) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// СРАЗУ нормализуем путь - до всех операций с метриками
			requestPath := normalizePath(r.URL.Path)

			if !shouldRecordMetrics(requestPath) {
				next.ServeHTTP(w, r)
				return
			}

			// Используем шаблон маршрута, если он доступен, чтобы избежать взрыва кардинальности.
			pathLabel := requestPath
			if rctx := chi.RouteContext(r.Context()); rctx != nil {
				if routePattern := strings.TrimSpace(rctx.RoutePattern()); routePattern != "" {
					pathLabel = normalizePath(routePattern)
				} else if strings.HasPrefix(requestPath, "/api") {
					// Все нематченные API-сканеры сводим к одному лейблу.
					pathLabel = "/api/_unmatched"
				}
			} else if strings.HasPrefix(requestPath, "/api") {
				pathLabel = "/api/_unmatched"
			}

			// Увеличиваем счетчик входящих запросов с атрибутами
			// НЕ добавляем client_ip чтобы избежать взрыва кардинальности метрик
			instruments.requestCounter.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", pathLabel),
				),
			)

			// Засекаем время начала обработки
			startTime := time.Now()

			// Создаем wrapper для ResponseWriter для отслеживания статуса
			sw := &statusWriter{
				ResponseWriter: w,
				status:         http.StatusOK, // Явная инициализация статуса
				written:        false,
			}

			// Выполняем обработчик
			next.ServeHTTP(sw, r)
			duration := time.Since(startTime)

			// Определяем статус ответа
			// Разделяем типы ошибок для security/ops сигналов (сканеры vs auth vs реальные 5xx)
			status := "success"
			switch {
			case sw.status >= 500:
				status = "server_error"
			case sw.status == http.StatusNotFound:
				status = "not_found"
			case sw.status == http.StatusUnauthorized || sw.status == http.StatusForbidden:
				status = "auth_error"
			case sw.status == http.StatusTooManyRequests:
				status = "rate_limit"
			case sw.status >= 400:
				status = "client_error"
			}

			// Увеличиваем счетчик ответов с атрибутами
			// НЕ добавляем client_ip чтобы избежать взрыва кардинальности метрик
			instruments.responseCounter.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("status", status),
					attribute.String("status_code", fmt.Sprintf("%d", sw.status)),
					attribute.String("method", r.Method),
					attribute.String("path", pathLabel),
				),
			)

			// Записываем время выполнения в гистограмму
			// НЕ добавляем client_ip чтобы избежать взрыва кардинальности метрик
			instruments.responseTimeHistogram.Record(r.Context(), duration.Seconds(),
				metric.WithAttributes(
					attribute.String("status", status),
					attribute.String("path", pathLabel),
					attribute.String("method", r.Method),
				),
			)
		})
	}
}

// HTTPMiddleware создает HTTP middleware для сбора метрик (глобальная функция для обратной совместимости).
// Middleware собирает метрики запросов, ответов и времени выполнения.
// bucketBoundaries - настраиваемые границы для гистограммы времени ответа
func HTTPMiddleware(ctx context.Context, bucketBoundaries []float64) func(http.Handler) http.Handler {
	return globalMetrics.HTTPMiddleware(ctx, bucketBoundaries)
}

// HTTPMiddleware создает HTTP middleware для сбора метрик для конкретного экземпляра Metrics.
// Middleware собирает метрики запросов, ответов и времени выполнения.
// bucketBoundaries - настраиваемые границы для гистограммы времени ответа
func (m *Metrics) HTTPMiddleware(ctx context.Context, bucketBoundaries []float64) func(http.Handler) http.Handler {
	// Проверяем, инициализирован ли MeterProvider
	meterProvider := otel.GetMeterProvider()
	if meterProvider == nil {
		// Метрики не инициализированы - возвращаем простой middleware без метрик
		return noOpMiddleware()
	}

	// Создаем инструменты метрик один раз при инициализации middleware.
	// Имя meter включает имя сервиса, чтобы в OpenTelemetry было понятно, какой компонент шлёт данные.
	meterName := "http-server"
	if m.config != nil && m.config.name != "" {
		meterName = m.config.name + "/http-server"
	}
	meter := otel.Meter(meterName)

	instruments, err := m.createHTTPInstruments(meter, bucketBoundaries)
	if err != nil {
		m.logError(ctx, "failed to create HTTP instruments", err)
		return noOpMiddleware()
	}

	return m.httpMiddlewareHandler(instruments)
}

// statusWriter оборачивает http.ResponseWriter для отслеживания статуса ответа.
// Пробрасывает опциональные интерфейсы (Flusher, Hijacker, Pusher), чтобы не ломать
// SSE/streaming, WebSockets и HTTP/2 server push.
type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (sw *statusWriter) WriteHeader(code int) {
	if sw.written {
		// WriteHeader уже был вызван, игнорируем повторный вызов
		return
	}
	sw.status = code
	sw.written = true
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if !sw.written {
		// Если Write вызван без WriteHeader, статус по умолчанию 200
		sw.status = http.StatusOK
		sw.written = true
	}
	return sw.ResponseWriter.Write(b)
}

// Flush реализует http.Flusher для поддержки SSE/streaming.
func (sw *statusWriter) Flush() {
	if f, ok := sw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack реализует http.Hijacker для поддержки WebSockets и т.п.
func (sw *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := sw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, ErrNotHijacker
}

// Push реализует http.Pusher для HTTP/2 server push.
func (sw *statusWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := sw.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return ErrNotPusher
}
