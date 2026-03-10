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

func (s *taskService) GetHistory(ctx context.Context, taskID, userID uuid.UUID) ([]model.TaskHistory, error) {
	var history []model.TaskHistory
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		task, errTx := s.taskRepo.GetByID(ctx, tx, taskID)
		if errTx != nil {
			return errTx
		}
		if _, errTx := s.memberRepo.GetMember(ctx, tx, task.TeamID, userID); errTx != nil {
			if errors.Is(errTx, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return errTx
		}
		history, errTx = s.historyRepo.ListHistoryByTaskID(ctx, tx, taskID)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		logger.Error(ctx, "GetHistory failed", zap.Error(err))
		return nil, err
	}
	return history, nil
}
