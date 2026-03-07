package task_v1_test

import (
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

func (s *APISuite) TestGetByID_Success() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	teamID := uuid.New()
	task := &model.Task{
		ID:          taskID,
		Title:       "My Task",
		Description: "Description",
		Status:      model.TaskStatusTodo,
		TeamID:      teamID,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.taskService.On("GetByID", mock.Anything, taskID, userID).Return(task, nil).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), testTaskID, resp.ID)
	assert.Equal(s.T(), "My Task", resp.Title)
	assert.Equal(s.T(), model.TaskStatusTodo, resp.Status)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_NoAuth() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestGetByID_InvalidUUID() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/not-a-uuid", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *APISuite) TestGetByID_NotFound() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.taskService.On("GetByID", mock.Anything, taskID, userID).Return(nil, model.ErrTaskNotFound).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_Forbidden() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.taskService.On("GetByID", mock.Anything, taskID, userID).Return(nil, model.ErrForbidden).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_InternalError() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.taskService.On("GetByID", mock.Anything, taskID, userID).Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.taskService.AssertExpectations(s.T())
}
