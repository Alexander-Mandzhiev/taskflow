package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"mkk/internal/module/identity/user/model"
	"mkk/pkg/logger"
)

func (s *userService) Update(ctx context.Context, id string, input *model.UserInput) (*model.User, error) {
	if input == nil {
		logger.Warn(ctx, "Update user: nil input")
		return nil, model.ErrNilInput
	}
	var user *model.User
	err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		user, errTx = s.repo.Update(ctx, tx, id, input)
		return errTx
	})
	if err != nil {
		logger.Error(ctx, "Update user: transaction failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}
