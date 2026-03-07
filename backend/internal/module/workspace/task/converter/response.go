package converter

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/task/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

const timeFormat = time.RFC3339

// TaskToResponse конвертирует доменную задачу в DTO ответа.
func TaskToResponse(t *model.Task) dto.TaskResponse {
	if t == nil {
		return dto.TaskResponse{}
	}
	resp := dto.TaskResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		TeamID:      t.TeamID.String(),
		CreatedBy:   t.CreatedBy.String(),
		CreatedAt:   t.CreatedAt.Format(timeFormat),
		UpdatedAt:   t.UpdatedAt.Format(timeFormat),
	}
	if t.AssigneeID != nil {
		s := t.AssigneeID.String()
		resp.AssigneeID = &s
	}
	if t.CompletedAt != nil {
		s := t.CompletedAt.Format(timeFormat)
		resp.CompletedAt = &s
	}
	return resp
}

// TasksToResponse конвертирует список задач в DTO ответа.
func TasksToResponse(tasks []*model.Task) []dto.TaskResponse {
	if len(tasks) == 0 {
		return nil
	}
	out := make([]dto.TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, TaskToResponse(t))
	}
	return out
}

// TaskHistoryToResponse конвертирует историю изменений задачи в DTO ответа.
func TaskHistoryToResponse(taskID string, entries []*model.TaskHistory) dto.TaskHistoryResponse {
	if entries == nil {
		entries = []*model.TaskHistory{}
	}
	items := make([]dto.TaskHistoryEntryResponse, 0, len(entries))
	for _, e := range entries {
		items = append(items, dto.TaskHistoryEntryResponse{
			ID:        e.ID.String(),
			TaskID:    e.TaskID.String(),
			ChangedBy: e.ChangedBy.String(),
			FieldName: e.FieldName,
			OldValue:  e.OldValue,
			NewValue:  e.NewValue,
			ChangedAt: e.ChangedAt.Format(timeFormat),
		})
	}
	return dto.TaskHistoryResponse{TaskID: taskID, Entries: items}
}

// TeamTaskStatsToResponse конвертирует статистику по команде в DTO ответа.
func TeamTaskStatsToResponse(s *model.TeamTaskStats) dto.TeamTaskStatsResponse {
	if s == nil {
		return dto.TeamTaskStatsResponse{}
	}
	return dto.TeamTaskStatsResponse{
		TeamID:         s.TeamID.String(),
		TeamName:       s.TeamName,
		MemberCount:    s.MemberCount,
		DoneTasksCount: s.DoneTasksCount,
	}
}

// TeamTaskStatsListToResponse конвертирует список статистики в DTO ответа.
func TeamTaskStatsListToResponse(items []*model.TeamTaskStats) dto.TeamTaskStatsListResponse {
	if len(items) == 0 {
		return dto.TeamTaskStatsListResponse{Items: []dto.TeamTaskStatsResponse{}}
	}
	out := make([]dto.TeamTaskStatsResponse, 0, len(items))
	for _, s := range items {
		out = append(out, TeamTaskStatsToResponse(s))
	}
	return dto.TeamTaskStatsListResponse{Items: out}
}

// TopCreatorToResponse конвертирует топ создателя в DTO ответа.
func TopCreatorToResponse(c *model.TeamTopCreator) dto.TopCreatorResponse {
	if c == nil {
		return dto.TopCreatorResponse{}
	}
	return dto.TopCreatorResponse{
		TeamID:       c.TeamID.String(),
		UserID:       c.UserID.String(),
		Rank:         c.Rank,
		CreatedCount: c.CreatedCount,
	}
}

// TopCreatorsListToResponse конвертирует список топ создателей в DTO ответа.
func TopCreatorsListToResponse(items []*model.TeamTopCreator) dto.TopCreatorsListResponse {
	if len(items) == 0 {
		return dto.TopCreatorsListResponse{Items: []dto.TopCreatorResponse{}}
	}
	out := make([]dto.TopCreatorResponse, 0, len(items))
	for _, c := range items {
		out = append(out, TopCreatorToResponse(c))
	}
	return dto.TopCreatorsListResponse{Items: out}
}

// CommentToResponse конвертирует доменный комментарий в DTO ответа.
func CommentToResponse(c *model.TaskComment) dto.CommentResponse {
	if c == nil {
		return dto.CommentResponse{}
	}
	return dto.CommentResponse{
		ID:        c.ID.String(),
		TaskID:    c.TaskID.String(),
		UserID:    c.UserID.String(),
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(timeFormat),
		UpdatedAt: c.UpdatedAt.Format(timeFormat),
	}
}

// CommentsToResponse конвертирует список комментариев в DTO ответа.
func CommentsToResponse(comments []*model.TaskComment) dto.CommentListResponse {
	if len(comments) == 0 {
		return dto.CommentListResponse{Items: []dto.CommentResponse{}}
	}
	out := make([]dto.CommentResponse, 0, len(comments))
	for _, c := range comments {
		out = append(out, CommentToResponse(c))
	}
	return dto.CommentListResponse{Items: out}
}
