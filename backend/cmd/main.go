package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"mkk/pkg/closer"
	"mkk/pkg/config"
	"mkk/internal/app"
	"mkk/pkg/logger"
)

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	cfg, err := config.Load(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось загрузить конфигурацию", zap.Error(err))
		return
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx, cfg)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	if err = a.Start(appCtx); err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
