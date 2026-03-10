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

// Create создаёт комментарий к задаче. userID должен быть участником команды задачи.
func (s *commentService) Create(ctx context.Context, taskID, userID uuid.UUID, content string) (model.TaskComment, error) {
	var comment model.TaskComment
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
		comment, errTx = s.commentRepo.CreateComment(ctx, tx, taskID, userID, content)
		return errTx
	}); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return model.TaskComment{}, err
		}
		logger.Error(ctx, "Create comment failed", zap.Error(err))
		return model.TaskComment{}, err
	}
	return comment, nil
}
