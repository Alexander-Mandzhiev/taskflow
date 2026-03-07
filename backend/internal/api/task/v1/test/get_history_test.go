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

func (s *APISuite) TestGetHistory_Success() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	entries := []*model.TaskHistory{
		{
			ID:        uuid.New(),
			TaskID:    taskID,
			ChangedBy: userID,
			FieldName: "title",
			OldValue:  "Old",
			NewValue:  "New",
			ChangedAt: time.Now(),
		},
	}

	s.taskService.On("GetHistory", mock.Anything, taskID, userID).Return(entries, nil).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/history", s.api.GetHistory)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/history", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		TaskID  string `json:"task_id"`
		Entries []struct {
			FieldName string `json:"field_name"`
			OldValue  string `json:"old_value"`
			NewValue  string `json:"new_value"`
		} `json:"entries"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), testTaskID, resp.TaskID)
	assert.Len(s.T(), resp.Entries, 1)
	assert.Equal(s.T(), "title", resp.Entries[0].FieldName)
	assert.Equal(s.T(), "Old", resp.Entries[0].OldValue)
	assert.Equal(s.T(), "New", resp.Entries[0].NewValue)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetHistory_NoAuth() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}/history", s.api.GetHistory)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/history", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestGetHistory_InvalidTaskID() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}/history", s.api.GetHistory)
	req := httptest.NewRequest(http.MethodGet, "/tasks/not-a-uuid/history", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "GetHistory")
}

func (s *APISuite) TestGetHistory_NotFound() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.taskService.On("GetHistory", mock.Anything, taskID, userID).Return(nil, model.ErrTaskNotFound).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/history", s.api.GetHistory)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/history", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetHistory_InternalError() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.taskService.On("GetHistory", mock.Anything, taskID, userID).Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/history", s.api.GetHistory)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/history", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.taskService.AssertExpectations(s.T())
}
