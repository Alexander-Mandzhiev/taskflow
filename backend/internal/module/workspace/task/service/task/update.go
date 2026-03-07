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

	current, err := s.taskRepo.GetByID(ctx, nil, taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		return nil, err
	}
	if _, err := s.teamSvc.GetMember(ctx, current.TeamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return nil, model.ErrTaskNotFound
		}
		return nil, err
	}
	if input.AssigneeID != nil {
		if _, err := s.teamSvc.GetMember(ctx, current.TeamID, *input.AssigneeID); err != nil {
			if errors.Is(err, teamModel.ErrMemberNotFound) {
				return nil, model.ErrAssigneeNotInTeam
			}
			return nil, err
		}
	}

	var updated *model.Task
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if errTx := s.taskRepo.Update(ctx, tx, taskID, input); errTx != nil {
			return errTx
		}
		now := time.Now()
		if errTx := writeHistory(ctx, tx, s.historyRepo, taskID, userID, now, current, input); errTx != nil {
			return errTx
		}
		var errTx error
		updated, errTx = s.taskRepo.GetByID(ctx, tx, taskID)
		return errTx
	}); err != nil {
		logger.Error(ctx, "Update task failed", zap.Error(err))
		return nil, err
	}
	return updated, nil
}

type historyWriter interface {
	CreateHistoryEntry(ctx context.Context, tx *sqlx.Tx, entry *model.TaskHistory) error
}

func writeHistory(ctx context.Context, tx *sqlx.Tx, repo historyWriter, taskID, changedBy uuid.UUID, now time.Time, current *model.Task, input *model.TaskInput) error {
	assigneeStr := func(p *uuid.UUID) string {
		if p == nil {
			return ""
		}
		return p.String()
	}
	entries := []struct {
		field    string
		old, new string
	}{
		{"title", current.Title, input.Title},
		{"description", current.Description, input.Description},
		{"status", current.Status, input.Status},
		{"assignee_id", assigneeStr(current.AssigneeID), assigneeStr(input.AssigneeID)},
	}
	for _, e := range entries {
		if e.old == e.new {
			continue
		}
		entry := &model.TaskHistory{
			TaskID:    taskID,
			ChangedBy: changedBy,
			FieldName: e.field,
			OldValue:  e.old,
			NewValue:  e.new,
			ChangedAt: now,
		}
		if err := repo.CreateHistoryEntry(ctx, tx, entry); err != nil {
			return err
		}
	}
	return nil
}
