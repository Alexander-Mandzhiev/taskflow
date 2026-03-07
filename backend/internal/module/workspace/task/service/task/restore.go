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
	var restored *model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		task, errTx := s.taskRepo.GetByIDIncludeDeleted(ctx, tx, taskID)
		if errTx != nil {
			return errTx
		}
		if _, errTx := s.teamRepo.GetMember(ctx, tx, task.TeamID, userID); errTx != nil {
			if errors.Is(errTx, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return errTx
		}
		if errTx := s.taskRepo.Restore(ctx, tx, taskID); errTx != nil {
			return errTx
		}
		restored, errTx = s.taskRepo.GetByID(ctx, tx, taskID)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		logger.Error(ctx, "Restore task failed", zap.Error(err))
		return nil, err
	}
	return restored, nil
}
