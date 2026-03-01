package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init инициализирует глобальный логгер с функциональными опциями.
// Поддерживает одновременную запись в stdout и OTLP коллектор.
// Повторные вызовы не переинициализируют логгер (защищено sync.Once).
func Init(ctx context.Context, opts ...Option) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	initOnce.Do(func() {
		initLogger(ctx, cfg)
	})

	if globalLogger == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// InitDefault инициализирует логгер с дефолтными настройками для первичной загрузки.
// Используется в main.go до загрузки конфигурации.
func InitDefault() error {
	initOnce.Do(func() {
		initLogger(context.Background(), defaultConfig())
	})

	if globalLogger == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// initLogger выполняет основную логику инициализации логгера
func initLogger(ctx context.Context, cfg *Config) {
	level = zap.NewAtomicLevelAt(parseLevel(cfg.level))

	cores := buildCores(ctx, cfg)
	zapLogger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))

	globalLogger = &logger{zapLogger: zapLogger}
}

// Reinit сбрасывает состояние и инициализирует логгер заново с функциональными опциями.
func Reinit(ctx context.Context, opts ...Option) error {
	resetGlobalState() //nolint:contextcheck // internal cleanup function
	return Init(ctx, opts...)
}

// Shutdown корректно завершает работу логгера.
// Останавливает OTLP provider с таймаутом для отправки оставшихся логов.
func Shutdown(ctx context.Context, timeout time.Duration) error {
	if otelProvider != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return otelProvider.Shutdown(shutdownCtx)
	}

	return nil
}

// resetGlobalState сбрасывает глобальное состояние логгера
func resetGlobalState() {
	globalLogger = nil
	initOnce = sync.Once{}
	level = zap.AtomicLevel{}
	if otelProvider != nil {
		_ = otelProvider.Shutdown(context.Background()) //nolint:gosec // best effort cleanup
		otelProvider = nil
	}
}

// InitForBenchmark настраивает NOP-логгер для бенчмарков/тестов.
func InitForBenchmark() {
	core := zapcore.NewNopCore()

	globalLogger = &logger{
		zapLogger: zap.New(core),
	}
}

// SetNopLogger устанавливает глобальный NOP-логгер (полезно в тестах).
func SetNopLogger() {
	globalLogger = &logger{zapLogger: zap.NewNop()}
}
