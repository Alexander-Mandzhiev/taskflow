package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"mkk/pkg/closer"
	"mkk/pkg/logger"
)

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	_ = appCtx
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "Shutdown error", zap.Error(err))
	}
}
