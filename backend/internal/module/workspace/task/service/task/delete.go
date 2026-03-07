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
	task, err := s.taskRepo.GetByID(ctx, nil, taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return err
		}
		return err
	}
	if _, err := s.teamSvc.GetMember(ctx, task.TeamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return model.ErrForbidden
		}
		return err
	}

	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.taskRepo.SoftDelete(ctx, tx, taskID)
	}); err != nil {
		logger.Error(ctx, "Delete task failed", zap.Error(err))
		return err
	}
	return nil
}
