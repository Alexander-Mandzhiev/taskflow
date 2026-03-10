package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *userService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := s.repo.GetByEmail(ctx, nil, email)
	if err != nil {
		logger.Error(ctx, "GetByEmail failed", zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}
