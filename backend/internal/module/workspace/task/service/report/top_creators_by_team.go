package report

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskReportService) TopCreatorsByTeam(ctx context.Context, userID uuid.UUID, since time.Time, limit int) ([]*model.TeamTopCreator, error) {
	if limit <= 0 {
		limit = 3
	}
	userTeams, err := s.teamSvc.ListByUserID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "ListByUserID failed", zap.Error(err))
		return nil, err
	}
	teamIDs := make(map[uuid.UUID]struct{})
	for _, t := range userTeams {
		teamIDs[t.ID] = struct{}{}
	}

	all, err := s.reportRepo.TopCreatorsByTeam(ctx, nil, since, limit)
	if err != nil {
		logger.Error(ctx, "TopCreatorsByTeam failed", zap.Error(err))
		return nil, err
	}

	out := make([]*model.TeamTopCreator, 0, len(all))
	for _, c := range all {
		if _, ok := teamIDs[c.TeamID]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}
