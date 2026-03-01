package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"mkk/internal/module/identity/user/model"
	"mkk/pkg/logger"
)

func (s *userService) Create(ctx context.Context, input *model.UserInput, passwordHash string) (*model.User, error) {
	if input == nil {
		logger.Warn(ctx, "Create user: nil input")
		return nil, model.ErrNilInput
	}

	var user *model.User

	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		user, errTx = s.repo.Create(ctx, tx, input, passwordHash)
		return errTx
	}); err != nil {
		logger.Error(ctx, "Create user: transaction failed", zap.Error(err))
		return nil, err
	}

	return user, nil
}
