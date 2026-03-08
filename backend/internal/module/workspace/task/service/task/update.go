package task

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskService) Update(ctx context.Context, userID, taskID uuid.UUID, input *model.TaskInput) (*model.Task, error) {
	if input == nil {
		return nil, model.ErrNilInput
	}
	if err := model.ValidateTaskInput(input); err != nil {
		return nil, err
	}

	var updated *model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		current, errTx := s.taskRepo.GetByID(ctx, tx, taskID)
		if errTx != nil {
			return errTx
		}
		if _, errTx := s.memberRepo.GetMember(ctx, tx, current.TeamID, userID); errTx != nil {
			if errors.Is(errTx, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return errTx
		}
		if input.AssigneeID != nil {
			if _, errTx := s.memberRepo.GetMember(ctx, tx, current.TeamID, *input.AssigneeID); errTx != nil {
				if errors.Is(errTx, teamModel.ErrMemberNotFound) {
					return model.ErrAssigneeNotInTeam
				}
				return errTx
			}
		}
		if errTx := s.taskRepo.Update(ctx, tx, taskID, input); errTx != nil {
			return errTx
		}
		now := time.Now()
		if errTx := s.recordUpdateHistory(ctx, tx, taskID, userID, current, input, now); errTx != nil {
			return errTx
		}
		updated, errTx = s.taskRepo.GetByID(ctx, tx, taskID)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) || errors.Is(err, model.ErrAssigneeNotInTeam) {
			return nil, err
		}
		logger.Error(ctx, "Update task failed", zap.Error(err))
		return nil, err
	}
	return updated, nil
}

// recordUpdateHistory пишет в task_history записи только для изменившихся полей.
func (s *taskService) recordUpdateHistory(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, current *model.Task, input *model.TaskInput, now time.Time) error {
	if err := s.recordHistoryIfChanged(ctx, tx, taskID, userID, "title", current.Title, input.Title, now); err != nil {
		return err
	}
	if err := s.recordHistoryIfChanged(ctx, tx, taskID, userID, "description", current.Description, input.Description, now); err != nil {
		return err
	}
	if err := s.recordHistoryIfChanged(ctx, tx, taskID, userID, "status", current.Status, input.Status, now); err != nil {
		return err
	}
	oldAssignee := ""
	if current.AssigneeID != nil {
		oldAssignee = current.AssigneeID.String()
	}
	newAssignee := ""
	if input.AssigneeID != nil {
		newAssignee = input.AssigneeID.String()
	}
	return s.recordHistoryIfChanged(ctx, tx, taskID, userID, "assignee_id", oldAssignee, newAssignee, now)
}

func (s *taskService) recordHistoryIfChanged(ctx context.Context, tx *sqlx.Tx, taskID, userID uuid.UUID, fieldName, oldVal, newVal string, now time.Time) error {
	if oldVal == newVal {
		return nil
	}
	return s.historyRepo.CreateHistoryEntry(ctx, tx, &model.TaskHistory{
		TaskID: taskID, ChangedBy: userID, FieldName: fieldName, OldValue: oldVal, NewValue: newVal, ChangedAt: now,
	})
}
