package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *userService) Create(ctx context.Context, input model.UserInput, passwordHash string) (model.User, error) {
	if input.Email == "" && input.Name == "" {
		logger.Warn(ctx, "Create user: nil input")
		return model.User{}, model.ErrNilInput
	}

	var user model.User
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		user, errTx = s.repo.Create(ctx, tx, input, passwordHash)
		return errTx
	}); err != nil {
		logger.Error(ctx, "Create user failed", zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}
