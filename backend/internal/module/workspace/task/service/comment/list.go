package comment

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

// ListByTaskID возвращает комментарии к задаче. userID должен быть участником команды задачи.
func (s *Service) ListByTaskID(ctx context.Context, taskID, userID uuid.UUID) ([]*model.TaskComment, error) {
	var comments []*model.TaskComment
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
		comments, errTx = s.commentRepo.ListCommentsByTaskID(ctx, tx, taskID)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return nil, err
		}
		logger.Error(ctx, "ListByTaskID comments failed", zap.Error(err))
		return nil, err
	}
	return comments, nil
}
