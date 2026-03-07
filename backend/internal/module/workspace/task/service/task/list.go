package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

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

	var items []*model.Task
	var total int
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if _, err := s.teamRepo.GetMember(ctx, tx, *filter.TeamID, userID); err != nil {
			if errors.Is(err, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return err
		}
		var errTx error
		items, total, errTx = s.taskRepo.List(ctx, tx, filter)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, 0, err
		}
		return nil, 0, err
	}
	return items, total, nil
}
