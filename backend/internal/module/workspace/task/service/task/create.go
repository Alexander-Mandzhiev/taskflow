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

func (s *taskService) Create(ctx context.Context, userID, teamID uuid.UUID, input model.TaskInput) (model.Task, error) {
	prepared := input
	if prepared.Status == "" {
		prepared.Status = model.TaskStatusTodo
	}

	var created model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if _, err := s.memberRepo.GetMember(ctx, tx, teamID, userID); err != nil {
			if errors.Is(err, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return err
		}
		if input.AssigneeID != nil {
			if _, err := s.memberRepo.GetMember(ctx, tx, teamID, *input.AssigneeID); err != nil {
				if errors.Is(err, teamModel.ErrMemberNotFound) {
					return model.ErrAssigneeNotInTeam
				}
				return err
			}
		}
		var errTx error
		created, errTx = s.taskRepo.Create(ctx, tx, teamID, prepared, userID)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) || errors.Is(err, model.ErrAssigneeNotInTeam) {
			return model.Task{}, err
		}
		logger.Error(ctx, "Create task failed", zap.Error(err))
		return model.Task{}, err
	}
	return created, nil
}
