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

func (s *taskService) Create(ctx context.Context, userID, teamID uuid.UUID, input *model.TaskInput) (*model.Task, error) {
	if input == nil {
		logger.Warn(ctx, "Create task: nil input")
		return nil, model.ErrNilInput
	}
	if err := model.ValidateTaskInput(input); err != nil {
		return nil, err
	}

	if _, err := s.teamSvc.GetMember(ctx, teamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return nil, model.ErrTaskNotFound
		}
		return nil, err
	}

	if input.AssigneeID != nil {
		if _, err := s.teamSvc.GetMember(ctx, teamID, *input.AssigneeID); err != nil {
			if errors.Is(err, teamModel.ErrMemberNotFound) {
				return nil, model.ErrAssigneeNotInTeam
			}
			return nil, err
		}
	}

	prepared := *input
	if prepared.Status == "" {
		prepared.Status = model.TaskStatusTodo
	}

	var created *model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		created, errTx = s.taskRepo.Create(ctx, tx, teamID, &prepared, userID)
		return errTx
	}); err != nil {
		logger.Error(ctx, "Create task failed", zap.Error(err))
		return nil, err
	}
	return created, nil
}
