package apiclient

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

const (
	defaultPasswordLength = 12
)

// FakeRegister возвращает запрос регистрации со случайными данными (email, пароль ≥8 символов, имя).
func FakeRegister() RegisterRequest {
	return RegisterRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, defaultPasswordLength),
		Name:     gofakeit.Name(),
	}
}

// FakeLogin возвращает запрос входа (те же email и пароль, что переданы — для использования после FakeRegister).
func FakeLogin(email, password string) LoginRequest {
	return LoginRequest{Email: email, Password: password}
}

// FakeCreateTeam возвращает запрос создания команды со случайным названием.
func FakeCreateTeam() CreateTeamRequest {
	return CreateTeamRequest{
		Name: gofakeit.Company() + " Team",
	}
}

// FakeCreateTask возвращает запрос создания задачи; teamID нужно передать после создания команды.
func FakeCreateTask(teamID string) CreateTaskRequest {
	statuses := []string{"todo", "in_progress", "done"}
	return CreateTaskRequest{
		TeamID:      teamID,
		Title:       gofakeit.Sentence(4),
		Description: gofakeit.Paragraph(1, 2, 5, " "),
		Status:      statuses[gofakeit.Number(0, len(statuses)-1)],
	}
}

// FakeCreateTaskWithAssignee — как FakeCreateTask, но с assignee_id (uuid).
func FakeCreateTaskWithAssignee(teamID, assigneeID string) CreateTaskRequest {
	t := FakeCreateTask(teamID)
	t.AssigneeID = &assigneeID
	return t
}

// FakeInvite возвращает запрос приглашения в команду (email, роль admin или member).
func FakeInvite(email, role string) InviteRequest {
	if role != "admin" && role != "member" {
		role = "member"
	}
	return InviteRequest{Email: email, Role: role}
}

// FakeCreateComment возвращает запрос создания комментария со случайным текстом.
func FakeCreateComment() CreateCommentRequest {
	return CreateCommentRequest{
		Content: gofakeit.Sentence(6),
	}
}

// FakeUpdateTask возвращает запрос обновления задачи (title, description, status).
func FakeUpdateTask() UpdateTaskRequest {
	statuses := []string{"todo", "in_progress", "done"}
	return UpdateTaskRequest{
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Paragraph(1, 1, 4, " "),
		Status:      statuses[gofakeit.Number(0, len(statuses)-1)],
	}
}

// Seed устанавливает seed для gofakeit (для воспроизводимых прогонов).
func Seed(seed int64) {
	_ = gofakeit.Seed(seed)
}

// RandomTeamID возвращает случайный UUID (для тестов, когда нужен несуществующий team_id).
func RandomTeamID() string {
	return uuid.New().String()
}
