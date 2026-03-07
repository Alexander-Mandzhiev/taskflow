package team

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TeamWithRole, error) {
	teams, err := s.repo.ListByUserID(ctx, nil, userID)
	if err != nil {
		logger.Error(ctx, "ListByUserID failed", zap.Error(err))
		return nil, err
	}
	return teams, nil
}
