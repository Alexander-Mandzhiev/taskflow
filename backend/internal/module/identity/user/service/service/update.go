package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *userService) Update(ctx context.Context, id string, input model.UserInput) (model.User, error) {
	if input.Email == "" && input.Name == "" {
		logger.Warn(ctx, "Update user: nil input")
		return model.User{}, model.ErrNilInput
	}

	var user model.User
	err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		user, errTx = s.repo.Update(ctx, tx, id, input)
		return errTx
	})
	if err != nil {
		logger.Error(ctx, "Update user failed", zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}
