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
		if _, errTx := s.teamRepo.GetMember(ctx, tx, current.TeamID, userID); errTx != nil {
			if errors.Is(errTx, teamModel.ErrMemberNotFound) {
				return model.ErrTaskNotFound
			}
			return errTx
		}
		if input.AssigneeID != nil {
			if _, errTx := s.teamRepo.GetMember(ctx, tx, current.TeamID, *input.AssigneeID); errTx != nil {
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
		if current.Title != input.Title {
			if errTx := s.historyRepo.CreateHistoryEntry(ctx, tx, &model.TaskHistory{TaskID: taskID, ChangedBy: userID, FieldName: "title", OldValue: current.Title, NewValue: input.Title, ChangedAt: now}); errTx != nil {
				return errTx
			}
		}
		if current.Description != input.Description {
			if errTx := s.historyRepo.CreateHistoryEntry(ctx, tx, &model.TaskHistory{TaskID: taskID, ChangedBy: userID, FieldName: "description", OldValue: current.Description, NewValue: input.Description, ChangedAt: now}); errTx != nil {
				return errTx
			}
		}
		if current.Status != input.Status {
			if errTx := s.historyRepo.CreateHistoryEntry(ctx, tx, &model.TaskHistory{TaskID: taskID, ChangedBy: userID, FieldName: "status", OldValue: current.Status, NewValue: input.Status, ChangedAt: now}); errTx != nil {
				return errTx
			}
		}
		oldAssignee := ""
		if current.AssigneeID != nil {
			oldAssignee = current.AssigneeID.String()
		}
		newAssignee := ""
		if input.AssigneeID != nil {
			newAssignee = input.AssigneeID.String()
		}
		if oldAssignee != newAssignee {
			if errTx := s.historyRepo.CreateHistoryEntry(ctx, tx, &model.TaskHistory{TaskID: taskID, ChangedBy: userID, FieldName: "assignee_id", OldValue: oldAssignee, NewValue: newAssignee, ChangedAt: now}); errTx != nil {
				return errTx
			}
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
