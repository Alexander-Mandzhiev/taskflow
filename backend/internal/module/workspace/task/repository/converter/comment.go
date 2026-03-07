package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToDomainTaskComment преобразует строку БД (TaskCommentRow) в доменную модель TaskComment.
func ToDomainTaskComment(r resources.TaskCommentRow) (model.TaskComment, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.TaskComment{}, err
	}
	taskID, err := uuid.Parse(r.TaskID)
	if err != nil {
		return model.TaskComment{}, err
	}
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.TaskComment{}, err
	}
	return model.TaskComment{
		ID:        id,
		TaskID:    taskID,
		UserID:    userID,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}, nil
}

// ToRepoTaskCommentCreateInput преобразует доменные данные в ресурс репозитория для INSERT в task_comments.
func ToRepoTaskCommentCreateInput(taskID, userID uuid.UUID, content string) resources.TaskCommentCreateInput {
	return resources.TaskCommentCreateInput{
		TaskID:  taskID.String(),
		UserID:  userID.String(),
		Content: content,
	}
}
