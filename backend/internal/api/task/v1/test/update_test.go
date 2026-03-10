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

func (s *APISuite) TestUpdate_Success() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	body, _ := json.Marshal(map[string]interface{}{
		"title":       "Updated Title",
		"description": "Updated desc",
		"status":      model.TaskStatusInProgress,
	})
	updated := model.Task{
		ID:          taskID,
		Title:       "Updated Title",
		Description: "Updated desc",
		Status:      model.TaskStatusInProgress,
		TeamID:      uuid.New(),
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.taskService.On("Update", mock.Anything, userID, taskID, mock.MatchedBy(func(in model.TaskInput) bool {
		return in.Title == "Updated Title" && in.Status == model.TaskStatusInProgress
	})).Return(updated, nil).Once()

	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+testTaskID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Title string `json:"title"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), "Updated Title", resp.Title)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdate_NoAuth() {
	body, _ := json.Marshal(map[string]string{"title": "T", "description": "", "status": model.TaskStatusTodo})
	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+testTaskID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestUpdate_InvalidTaskID() {
	body, _ := json.Marshal(map[string]string{"title": "T", "description": "", "status": model.TaskStatusTodo})
	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/not-a-uuid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "Update")
}

func (s *APISuite) TestUpdate_ValidationError_InvalidStatus() {
	body, _ := json.Marshal(map[string]string{"title": "T", "description": "", "status": "invalid"})
	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+testTaskID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "Update")
}

func (s *APISuite) TestUpdate_NotFound() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	body, _ := json.Marshal(map[string]string{"title": "T", "description": "", "status": model.TaskStatusTodo})
	s.taskService.On("Update", mock.Anything, userID, taskID, mock.Anything).Return(model.Task{}, model.ErrTaskNotFound).Once()

	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+testTaskID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdate_InternalError() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	body, _ := json.Marshal(map[string]string{"title": "T", "description": "", "status": model.TaskStatusTodo})
	s.taskService.On("Update", mock.Anything, userID, taskID, mock.Anything).Return(model.Task{}, assert.AnError).Once()

	r := chi.NewRouter()
	r.Put("/tasks/{id}", s.api.Update)
	req := httptest.NewRequest(http.MethodPut, "/tasks/"+testTaskID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.taskService.AssertExpectations(s.T())
}
