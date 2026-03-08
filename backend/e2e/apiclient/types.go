package apiclient

// Request/response types для e2e-клиента API (соответствуют контрактам backend).

// RegisterRequest — запрос регистрации.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"` //nolint:gosec // G117: e2e DTO, not a secret storage
	Name     string `json:"name,omitempty"`
}

// LoginRequest — запрос входа.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"` //nolint:gosec // G117: e2e DTO, not a secret storage
}

// CreateTeamRequest — запрос создания команды.
type CreateTeamRequest struct {
	Name string `json:"name"`
}

// InviteRequest — запрос приглашения в команду.
type InviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"` // admin | member
}

// InvitationResponse — данные приглашения в ответе Invite.
type InvitationResponse struct {
	ID        string `json:"id"`
	TeamID    string `json:"team_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ExpiresAt string `json:"expires_at"`
}

// InviteResponse — ответ POST /teams/{id}/invite.
type InviteResponse struct {
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
	Invitation InvitationResponse `json:"invitation,omitempty"`
}

// CreateTaskRequest — запрос создания задачи.
type CreateTaskRequest struct {
	TeamID      string  `json:"team_id"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status,omitempty"` // todo | in_progress | done
	AssigneeID  *string `json:"assignee_id,omitempty"`
}

// UpdateTaskRequest — запрос обновления задачи.
type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Status      string  `json:"status"` // todo | in_progress | done
	AssigneeID  *string `json:"assignee_id,omitempty"`
}

// CreateCommentRequest — запрос создания комментария.
type CreateCommentRequest struct {
	Content string `json:"content"`
}

// Team — команда в ответе API.
type Team struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TeamWithRole — команда с ролью пользователя (список).
type TeamWithRole struct {
	Team
	Role string `json:"role"`
}

// TeamWithMembersResponse — ответ GET /api/v1/teams/{id} (команда с участниками).
type TeamWithMembersResponse struct {
	Team    Team `json:"team"`
	Members []struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	} `json:"members,omitempty"`
}

// Task — задача в ответе API.
type Task struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssigneeID  *string `json:"assignee_id,omitempty"`
	TeamID      string  `json:"team_id"`
	CreatedBy   string  `json:"created_by"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

// TaskListResponse — список задач с пагинацией.
type TaskListResponse struct {
	Items  []Task `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// TaskHistoryResponse — история изменений задачи (GET /tasks/{id}/history).
type TaskHistoryResponse struct {
	TaskID  string         `json:"task_id"`
	Entries []HistoryEntry `json:"entries"`
}

// HistoryEntry — запись истории задачи.
type HistoryEntry struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	ChangedBy string `json:"changed_by"`
	FieldName string `json:"field_name"`
	OldValue  string `json:"old_value,omitempty"`
	NewValue  string `json:"new_value,omitempty"`
	ChangedAt string `json:"changed_at"`
}

// Comment — комментарий к задаче.
type Comment struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CommentListResponse — список комментариев к задаче (GET /tasks/{id}/comments).
type CommentListResponse struct {
	Items []Comment `json:"items"`
}

// TeamTaskStats — статистика по команде (отчёт).
type TeamTaskStats struct {
	TeamID         string `json:"team_id"`
	TeamName       string `json:"team_name"`
	MemberCount    int    `json:"member_count"`
	DoneTasksCount int    `json:"done_tasks_count"`
}

// ListTasksOpts — опции для списка задач.
type ListTasksOpts struct {
	TeamID     string
	Status     string
	AssigneeID string
	Limit      int
	Offset     int
}
