package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskService) Restore(ctx context.Context, userID, taskID uuid.UUID) (*model.Task, error) {
	task, err := s.taskRepo.GetByIDIncludeDeleted(ctx, nil, taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		return nil, err
	}
	if _, err := s.teamSvc.GetMember(ctx, task.TeamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return nil, model.ErrForbidden
		}
		return nil, err
	}

	var restored *model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if errTx := s.taskRepo.Restore(ctx, tx, taskID); errTx != nil {
			return errTx
		}
		restored, _ = s.taskRepo.GetByID(ctx, tx, taskID)
		return nil
	}); err != nil {
		logger.Error(ctx, "Restore task failed", zap.Error(err))
		return nil, err
	}
	return restored, nil
}
