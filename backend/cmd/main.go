package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/app"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/closer"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func main() {
	if err := logger.InitDefault(); err != nil {
		fmt.Fprintf(os.Stderr, "logger init: %v\n", err)
		os.Exit(1)
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	cfg, err := config.Load(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось загрузить конфигурацию", zap.Error(err))
		return
	}

	cl := closer.NewWithLogger(logger.Logger(), syscall.SIGINT, syscall.SIGTERM)
	defer gracefulShutdown(cl)

	a, err := app.New(appCtx, cfg, cl)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	if err = a.Start(appCtx); err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown(cl *closer.Closer) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := cl.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
