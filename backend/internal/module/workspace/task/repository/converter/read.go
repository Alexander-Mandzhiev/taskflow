package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToDomainTask преобразует строку БД (TaskRow) в доменную модель Task.
func ToDomainTask(r resources.TaskRow) (model.Task, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.Task{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.Task{}, err
	}
	createdBy, err := uuid.Parse(r.CreatedBy)
	if err != nil {
		return model.Task{}, err
	}
	var assigneeID *uuid.UUID
	if r.AssigneeID != nil && *r.AssigneeID != "" {
		a, err := uuid.Parse(*r.AssigneeID)
		if err != nil {
			return model.Task{}, err
		}
		assigneeID = &a
	}
	return model.Task{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
		AssigneeID:  assigneeID,
		TeamID:      teamID,
		CreatedBy:   createdBy,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}, nil
}

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

// ToDomainTeamTaskStats преобразует строку отчёта в model.TeamTaskStats.
func ToDomainTeamTaskStats(r resources.TeamTaskStatsRow) (model.TeamTaskStats, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamTaskStats{}, err
	}
	return model.TeamTaskStats{
		TeamID:         teamID,
		TeamName:       r.TeamName,
		MemberCount:    r.MemberCount,
		DoneTasksCount: r.DoneTasksCount,
	}, nil
}

// ToDomainTeamTopCreator преобразует строку отчёта в model.TeamTopCreator.
func ToDomainTeamTopCreator(r resources.TeamTopCreatorRow) (model.TeamTopCreator, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamTopCreator{}, err
	}
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.TeamTopCreator{}, err
	}
	return model.TeamTopCreator{
		TeamID:       teamID,
		UserID:       userID,
		Rank:         r.Rank,
		CreatedCount: r.CreatedCount,
	}, nil
}
