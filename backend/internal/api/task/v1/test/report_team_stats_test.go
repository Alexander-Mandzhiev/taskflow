package task_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

func (s *APISuite) TestReportTeamStats_Success() {
	userID := uuid.MustParse(testUserID)
	teamID := uuid.New()
	items := []model.TeamTaskStats{
		{
			TeamID:         teamID,
			TeamName:       "My Team",
			MemberCount:    5,
			DoneTasksCount: 12,
		},
	}

	s.reportService.On("TeamTaskStats", mock.Anything, userID, mock.AnythingOfType("time.Time")).Return(items, nil).Once()

	r := chi.NewRouter()
	r.Get("/reports/team-stats", s.api.ReportTeamStats)
	req := httptest.NewRequest(http.MethodGet, "/reports/team-stats?since_days=7", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []struct {
			TeamID         string `json:"team_id"`
			TeamName       string `json:"team_name"`
			MemberCount    int    `json:"member_count"`
			DoneTasksCount int    `json:"done_tasks_count"`
		} `json:"items"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp.Items, 1)
	assert.Equal(s.T(), "My Team", resp.Items[0].TeamName)
	assert.Equal(s.T(), 5, resp.Items[0].MemberCount)
	assert.Equal(s.T(), 12, resp.Items[0].DoneTasksCount)
	s.reportService.AssertExpectations(s.T())
}

func (s *APISuite) TestReportTeamStats_NoAuth() {
	r := chi.NewRouter()
	r.Get("/reports/team-stats", s.api.ReportTeamStats)
	req := httptest.NewRequest(http.MethodGet, "/reports/team-stats", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestReportTeamStats_InvalidSinceDays() {
	r := chi.NewRouter()
	r.Get("/reports/team-stats", s.api.ReportTeamStats)
	req := httptest.NewRequest(http.MethodGet, "/reports/team-stats?since_days=invalid", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.reportService.AssertNotCalled(s.T(), "TeamTaskStats")
}

func (s *APISuite) TestReportTeamStats_InternalError() {
	s.reportService.On("TeamTaskStats", mock.Anything, uuid.MustParse(testUserID), mock.AnythingOfType("time.Time")).
		Return([]model.TeamTaskStats(nil), assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/reports/team-stats", s.api.ReportTeamStats)
	req := httptest.NewRequest(http.MethodGet, "/reports/team-stats", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.reportService.AssertExpectations(s.T())
}
