package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"mkk/pkg/logger"
)

func (s *userService) ChangePassword(ctx context.Context, id, passwordHash string) error {
	err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.repo.UpdatePasswordHash(ctx, tx, id, passwordHash)
	})
	if err != nil {
		logger.Error(ctx, "ChangePassword failed", zap.Error(err))
		return err
	}
	return nil
}
