package report

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskReportService) TeamTaskStats(ctx context.Context, userID uuid.UUID, since time.Time) ([]model.TeamTaskStats, error) {
	userTeams, err := s.teamSvc.ListByUserID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "ListByUserID failed", zap.Error(err))
		return nil, err
	}
	teamIDs := make(map[uuid.UUID]struct{})
	for _, t := range userTeams {
		teamIDs[t.ID] = struct{}{}
	}

	all, err := s.reportRepo.TeamTaskStats(ctx, nil, since)
	if err != nil {
		logger.Error(ctx, "TeamTaskStats failed", zap.Error(err))
		return nil, err
	}

	out := make([]model.TeamTaskStats, 0, len(all))
	for _, st := range all {
		if _, ok := teamIDs[st.TeamID]; ok {
			out = append(out, st)
		}
	}
	return out, nil
}
