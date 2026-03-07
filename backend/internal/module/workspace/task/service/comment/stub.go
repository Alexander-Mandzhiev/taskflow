package comment

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/service"
)

var _ service.TaskCommentService = (*Stub)(nil)

// Stub — заглушка TaskCommentService. Возвращает model.ErrCommentNotImplemented.
// Используется пока не реализован репозиторий комментариев.
type Stub struct{}

// NewStub создаёт заглушку сервиса комментариев.
func NewStub() *Stub {
	return &Stub{}
}

// ListByTaskID возвращает ErrCommentNotImplemented.
func (s *Stub) ListByTaskID(ctx context.Context, taskID, userID uuid.UUID) ([]*model.TaskComment, error) {
	_ = ctx
	_ = taskID
	_ = userID
	return nil, model.ErrCommentNotImplemented
}

// Create возвращает ErrCommentNotImplemented.
func (s *Stub) Create(ctx context.Context, taskID, userID uuid.UUID, content string) (*model.TaskComment, error) {
	_ = ctx
	_ = taskID
	_ = userID
	_ = content
	return nil, model.ErrCommentNotImplemented
}
