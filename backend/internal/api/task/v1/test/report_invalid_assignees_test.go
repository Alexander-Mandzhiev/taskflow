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

func (s *APISuite) TestReportInvalidAssignees_Success() {
	userID := uuid.MustParse(testUserID)
	tasks := []*model.Task{
		{
			ID:          uuid.New(),
			Title:       "Task with bad assignee",
			TeamID:      uuid.New(),
			CreatedBy:   userID,
			Status:      model.TaskStatusTodo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	s.reportService.On("TasksWithInvalidAssignee", mock.Anything, userID).Return(tasks, nil).Once()

	r := chi.NewRouter()
	r.Get("/reports/invalid-assignees", s.api.ReportInvalidAssignees)
	req := httptest.NewRequest(http.MethodGet, "/reports/invalid-assignees", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp.Items, 1)
	assert.Equal(s.T(), "Task with bad assignee", resp.Items[0].Title)
	s.reportService.AssertExpectations(s.T())
}

func (s *APISuite) TestReportInvalidAssignees_NoAuth() {
	r := chi.NewRouter()
	r.Get("/reports/invalid-assignees", s.api.ReportInvalidAssignees)
	req := httptest.NewRequest(http.MethodGet, "/reports/invalid-assignees", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestReportInvalidAssignees_EmptyList() {
	s.reportService.On("TasksWithInvalidAssignee", mock.Anything, uuid.MustParse(testUserID)).
		Return([]*model.Task{}, nil).Once()

	r := chi.NewRouter()
	r.Get("/reports/invalid-assignees", s.api.ReportInvalidAssignees)
	req := httptest.NewRequest(http.MethodGet, "/reports/invalid-assignees", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []interface{} `json:"items"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Empty(s.T(), resp.Items)
	s.reportService.AssertExpectations(s.T())
}

func (s *APISuite) TestReportInvalidAssignees_InternalError() {
	s.reportService.On("TasksWithInvalidAssignee", mock.Anything, uuid.MustParse(testUserID)).
		Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/reports/invalid-assignees", s.api.ReportInvalidAssignees)
	req := httptest.NewRequest(http.MethodGet, "/reports/invalid-assignees", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.reportService.AssertExpectations(s.T())
}
