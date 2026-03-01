package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// StartSpan создает новый спан и возвращает его вместе с новым контекстом.
// Это удобная обертка над trace.Tracer.Start, которая использует глобальный трейсер.
//
// Разница между Tracer и TracerProvider:
// 1. TracerProvider - это фабрика трейсеров, которая:
//   - Управляет жизненным циклом трейсеров
//   - Настраивает экспорт трейсов
//   - Контролирует семплирование
//   - Хранит глобальные настройки
//
// 2. Tracer - это конкретный инструмент для создания спанов:
//   - Создает спаны для определенного сервиса/компонента
//   - Управляет связями между спанами
//   - Добавляет атрибуты к спанам
//   - Отслеживает контекст выполнения
//
// Имя трейсера (serviceName):
// - Используется для идентификации источника спанов в системе трассировки
// - Позволяет группировать спаны по сервисам в UI (например, в Jaeger)
// - Если трейсер с таким именем уже существует - он возвращается
// - Если нет - создается новый трейсер с этим именем
// - В нашем случае используется имя сервиса из конфигурации
//
// Создание спана:
// - При создании первого (корневого) спана генерируется новый trace ID
// - Все последующие спаны в цепочке наследуют этот trace ID
// - Trace ID связывает все спаны одного запроса/операции
// - Если в контексте уже есть trace ID, новый не генерируется
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	// Получаем трейсер из нашего провайдера
	return GetTracerProvider().Tracer("tracing-utils").Start(ctx, name, opts...)
}

// SpanFromContext возвращает текущий активный спан из контекста.
// Если спан не существует, возвращается NoopSpan.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext извлекает trace ID из контекста.
// Возвращает строку с ID трейса или пустую строку, если трейс не найден.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
