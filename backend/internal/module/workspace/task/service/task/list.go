package task

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *taskService) List(ctx context.Context, userID uuid.UUID, filter *model.TaskListFilter) ([]*model.Task, int, error) {
	if filter == nil {
		return nil, 0, model.ErrPaginationRequired
	}
	if filter.Limit <= 0 || filter.Offset < 0 {
		return nil, 0, model.ErrPaginationRequired
	}
	if filter.TeamID == nil {
		return nil, 0, model.ErrForbidden
	}
	if _, err := s.teamSvc.GetMember(ctx, *filter.TeamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return nil, 0, model.ErrForbidden
		}
		return nil, 0, err
	}
	return s.taskRepo.List(ctx, nil, filter)
}
