package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToDomainTaskHistory преобразует строку БД (TaskHistoryRow) в доменную модель TaskHistory.
func ToDomainTaskHistory(r resources.TaskHistoryRow) (model.TaskHistory, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.TaskHistory{}, err
	}
	taskID, err := uuid.Parse(r.TaskID)
	if err != nil {
		return model.TaskHistory{}, err
	}
	changedBy, err := uuid.Parse(r.ChangedBy)
	if err != nil {
		return model.TaskHistory{}, err
	}
	return model.TaskHistory{
		ID:        id,
		TaskID:    taskID,
		ChangedBy: changedBy,
		FieldName: r.FieldName,
		OldValue:  r.OldValue,
		NewValue:  r.NewValue,
		ChangedAt: r.ChangedAt,
	}, nil
}

// ToRepoTaskHistory преобразует доменную запись TaskHistory в ресурс для INSERT в task_history.
func ToRepoTaskHistory(entry *model.TaskHistory) resources.TaskHistoryRow {
	if entry == nil {
		return resources.TaskHistoryRow{}
	}
	return resources.TaskHistoryRow{
		ID:        entry.ID.String(),
		TaskID:    entry.TaskID.String(),
		ChangedBy: entry.ChangedBy.String(),
		FieldName: entry.FieldName,
		OldValue:  entry.OldValue,
		NewValue:  entry.NewValue,
		ChangedAt: entry.ChangedAt,
	}
}
