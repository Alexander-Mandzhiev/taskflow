package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"mkk/pkg/logger"
)

func (s *userService) Delete(ctx context.Context, id string) error {
	err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.repo.Delete(ctx, tx, id)
	})
	if err != nil {
		logger.Error(ctx, "Delete user failed", zap.Error(err))
		return err
	}
	return nil
}
