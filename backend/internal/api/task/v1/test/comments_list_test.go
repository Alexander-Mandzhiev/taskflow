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

func (s *APISuite) TestListComments_Success() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	comments := []model.TaskComment{
		{
			ID:        uuid.New(),
			TaskID:    taskID,
			UserID:    userID,
			Content:   "First comment",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.commentService.On("ListByTaskID", mock.Anything, taskID, userID).Return(comments, nil).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/comments", s.api.ListComments)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/comments", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []struct {
			ID      string `json:"id"`
			TaskID  string `json:"task_id"`
			UserID  string `json:"user_id"`
			Content string `json:"content"`
		} `json:"items"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp.Items, 1)
	assert.Equal(s.T(), "First comment", resp.Items[0].Content)
	assert.Equal(s.T(), testTaskID, resp.Items[0].TaskID)
	s.commentService.AssertExpectations(s.T())
}

func (s *APISuite) TestListComments_NoAuth() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}/comments", s.api.ListComments)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/comments", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestListComments_InvalidTaskID() {
	r := chi.NewRouter()
	r.Get("/tasks/{id}/comments", s.api.ListComments)
	req := httptest.NewRequest(http.MethodGet, "/tasks/not-a-uuid/comments", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.commentService.AssertNotCalled(s.T(), "ListByTaskID")
}

func (s *APISuite) TestListComments_NotFound() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.commentService.On("ListByTaskID", mock.Anything, taskID, userID).Return([]model.TaskComment(nil), model.ErrTaskNotFound).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/comments", s.api.ListComments)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/comments", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.commentService.AssertExpectations(s.T())
}

func (s *APISuite) TestListComments_InternalError() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	s.commentService.On("ListByTaskID", mock.Anything, taskID, userID).Return([]model.TaskComment(nil), assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/tasks/{id}/comments", s.api.ListComments)
	req := httptest.NewRequest(http.MethodGet, "/tasks/"+testTaskID+"/comments", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.commentService.AssertExpectations(s.T())
}
