package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, nil, id)
	if err != nil {
		logger.Error(ctx, "GetByID failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}
