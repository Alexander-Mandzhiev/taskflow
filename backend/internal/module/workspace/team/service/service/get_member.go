package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) GetMember(ctx context.Context, teamID, userID uuid.UUID) (*model.TeamMember, error) {
	member, err := s.repo.GetMember(ctx, nil, teamID, userID)
	if err != nil {
		logger.Error(ctx, "GetMember failed", zap.Error(err))
		return nil, err
	}
	return member, nil
}
