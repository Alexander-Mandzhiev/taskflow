package task_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

func (s *APISuite) TestCreate_Success() {
	teamID := uuid.New().String()
	body, _ := json.Marshal(map[string]interface{}{
		"team_id":      teamID,
		"title":        "New Task",
		"description":  "Description",
		"status":       model.TaskStatusTodo,
	})
	userID := uuid.MustParse(testUserID)
	created := &model.Task{
		ID:          uuid.New(),
		Title:       "New Task",
		Description: "Description",
		Status:      model.TaskStatusTodo,
		TeamID:      uuid.MustParse(teamID),
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.taskService.On("Create", mock.Anything, userID, uuid.MustParse(teamID), mock.MatchedBy(func(in *model.TaskInput) bool {
		return in != nil && in.Title == "New Task" && in.Description == "Description" && in.Status == model.TaskStatusTodo
	})).Return(created, nil).Once()

	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), "New Task", resp.Title)
	assert.Equal(s.T(), model.TaskStatusTodo, resp.Status)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreate_NoAuth() {
	body, _ := json.Marshal(map[string]string{"team_id": uuid.New().String(), "title": "Task"})
	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestCreate_InvalidJSON() {
	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreate_ValidationError_EmptyTitle() {
	body, _ := json.Marshal(map[string]string{"team_id": uuid.New().String(), "title": ""})
	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreate_ValidationError_InvalidTeamID() {
	body, _ := json.Marshal(map[string]string{"team_id": "not-uuid", "title": "Task"})
	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreate_InternalError() {
	teamID := uuid.New().String()
	body, _ := json.Marshal(map[string]string{"team_id": teamID, "title": "Task"})
	s.taskService.On("Create", mock.Anything, uuid.MustParse(testUserID), uuid.MustParse(teamID), mock.Anything).
		Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Post("/tasks", s.api.Create)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.taskService.AssertExpectations(s.T())
}
