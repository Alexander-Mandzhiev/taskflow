package service

import (
	"context"

	"go.uber.org/zap"

	"mkk/internal/module/identity/user/model"
	"mkk/pkg/logger"
)

func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, nil, email)
	if err != nil {
		logger.Error(ctx, "GetByEmail failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}
