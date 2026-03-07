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

func (s *APISuite) TestCreateComment_Success() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	content := "New comment"
	created := &model.TaskComment{
		ID:        uuid.New(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.commentService.On("Create", mock.Anything, taskID, userID, content).Return(created, nil).Once()

	body, _ := json.Marshal(map[string]string{"content": content})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp struct {
		ID      string `json:"id"`
		TaskID  string `json:"task_id"`
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), content, resp.Content)
	assert.Equal(s.T(), testTaskID, resp.TaskID)
	assert.Equal(s.T(), testUserID, resp.UserID)
	s.commentService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateComment_NoAuth() {
	body, _ := json.Marshal(map[string]string{"content": "Comment"})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestCreateComment_InvalidTaskID() {
	body, _ := json.Marshal(map[string]string{"content": "Comment"})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/not-a-uuid/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.commentService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreateComment_InvalidJSON() {
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.commentService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreateComment_ValidationError_EmptyContent() {
	body, _ := json.Marshal(map[string]string{"content": ""})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.commentService.AssertNotCalled(s.T(), "Create")
}

func (s *APISuite) TestCreateComment_NotFound() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	content := "Comment"
	s.commentService.On("Create", mock.Anything, taskID, userID, content).Return(nil, model.ErrTaskNotFound).Once()

	body, _ := json.Marshal(map[string]string{"content": content})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.commentService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateComment_InternalError() {
	taskID := uuid.MustParse(testTaskID)
	userID := uuid.MustParse(testUserID)
	content := "Comment"
	s.commentService.On("Create", mock.Anything, taskID, userID, content).Return(nil, assert.AnError).Once()

	body, _ := json.Marshal(map[string]string{"content": content})
	r := chi.NewRouter()
	r.Post("/tasks/{id}/comments", s.api.CreateComment)
	req := httptest.NewRequest(http.MethodPost, "/tasks/"+testTaskID+"/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.commentService.AssertExpectations(s.T())
}
