package dto

// TaskHistoryEntryResponse — одна запись истории изменений задачи.
type TaskHistoryEntryResponse struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	ChangedBy string `json:"changed_by"`
	FieldName string `json:"field_name"`
	OldValue  string `json:"old_value,omitempty"`
	NewValue  string `json:"new_value,omitempty"`
	ChangedAt string `json:"changed_at"`
}

// TaskHistoryResponse — история изменений задачи.
type TaskHistoryResponse struct {
	TaskID  string                     `json:"task_id"`
	Entries []TaskHistoryEntryResponse `json:"entries"`
}
