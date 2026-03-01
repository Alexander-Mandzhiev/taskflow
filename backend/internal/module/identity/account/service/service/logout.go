package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"mkk/pkg/logger"
)

// Logout удаляет сессию из кеша.
func (s *accountService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	err := s.sessionRepo.Delete(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "Logout: delete session failed", zap.Error(err))
		return err
	}
	return nil
}
