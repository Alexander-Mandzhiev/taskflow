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

func (s *taskService) Delete(ctx context.Context, userID, taskID uuid.UUID) error {
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		task, errTx := s.taskRepo.GetByID(ctx, tx, taskID)
		if errTx != nil {
			return errTx
		}
		if _, errTx := s.teamRepo.GetMember(ctx, tx, task.TeamID, userID); errTx != nil {
			if errors.Is(errTx, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return errTx
		}
		return s.taskRepo.SoftDelete(ctx, tx, taskID)
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return err
		}
		logger.Error(ctx, "Delete task failed", zap.Error(err))
		return err
	}
	return nil
}
