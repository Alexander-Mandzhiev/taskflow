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

func (s *APISuite) TestList_Success() {
	teamID := uuid.New()
	userID := uuid.MustParse(testUserID)
	tasks := []*model.Task{
		{
			ID:        uuid.New(),
			Title:     "Task 1",
			TeamID:    teamID,
			CreatedBy: userID,
			Status:    model.TaskStatusTodo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	s.taskService.On("List", mock.Anything, userID, mock.MatchedBy(func(f *model.TaskListFilter) bool {
		return f != nil && f.TeamID != nil && *f.TeamID == teamID && f.Limit == 20 && f.Offset == 0
	})).Return(tasks, 1, nil).Once()

	r := chi.NewRouter()
	r.Get("/tasks", s.api.List)
	req := httptest.NewRequest(http.MethodGet, "/tasks?team_id="+teamID.String()+"&limit=20", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
		Total int `json:"total"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp.Items, 1)
	assert.Equal(s.T(), "Task 1", resp.Items[0].Title)
	assert.Equal(s.T(), 1, resp.Total)
	s.taskService.AssertExpectations(s.T())
}

func (s *APISuite) TestList_NoAuth() {
	r := chi.NewRouter()
	r.Get("/tasks", s.api.List)
	req := httptest.NewRequest(http.MethodGet, "/tasks?team_id="+uuid.New().String(), nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestList_TeamIDRequired() {
	r := chi.NewRouter()
	r.Get("/tasks", s.api.List)
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "List")
}

func (s *APISuite) TestList_InvalidLimit() {
	r := chi.NewRouter()
	r.Get("/tasks", s.api.List)
	req := httptest.NewRequest(http.MethodGet, "/tasks?team_id="+uuid.New().String()+"&limit=0", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.taskService.AssertNotCalled(s.T(), "List")
}

func (s *APISuite) TestList_InternalError() {
	teamID := uuid.New()
	s.taskService.On("List", mock.Anything, uuid.MustParse(testUserID), mock.Anything).
		Return(nil, 0, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/tasks", s.api.List)
	req := httptest.NewRequest(http.MethodGet, "/tasks?team_id="+teamID.String()+"&limit=10", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.taskService.AssertExpectations(s.T())
}
