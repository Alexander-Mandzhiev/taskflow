package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *taskService) GetByID(ctx context.Context, taskID, userID uuid.UUID) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, nil, taskID)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		logger.Error(ctx, "GetByID task failed", zap.Error(err))
		return nil, err
	}

	if _, err := s.teamSvc.GetMember(ctx, task.TeamID, userID); err != nil {
		if errors.Is(err, teamModel.ErrMemberNotFound) {
			return nil, model.ErrForbidden
		}
		return nil, err
	}
	return task, nil
}
