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

func (s *APISuite) TestReportTopCreators_Success() {
	userID := uuid.MustParse(testUserID)
	teamID := uuid.New()
	items := []*model.TeamTopCreator{
		{
			TeamID:       teamID,
			UserID:      userID,
			Rank:        1,
			CreatedCount: 10,
		},
	}

	s.reportService.On("TopCreatorsByTeam", mock.Anything, userID, mock.AnythingOfType("time.Time"), 3).Return(items, nil).Once()

	r := chi.NewRouter()
	r.Get("/reports/top-creators", s.api.ReportTopCreators)
	req := httptest.NewRequest(http.MethodGet, "/reports/top-creators?since_days=30&limit=3", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Items []struct {
			TeamID       string `json:"team_id"`
			UserID       string `json:"user_id"`
			Rank         int    `json:"rank"`
			CreatedCount int64  `json:"created_count"`
		} `json:"items"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp.Items, 1)
	assert.Equal(s.T(), 1, resp.Items[0].Rank)
	assert.Equal(s.T(), int64(10), resp.Items[0].CreatedCount)
	s.reportService.AssertExpectations(s.T())
}

func (s *APISuite) TestReportTopCreators_NoAuth() {
	r := chi.NewRouter()
	r.Get("/reports/top-creators", s.api.ReportTopCreators)
	req := httptest.NewRequest(http.MethodGet, "/reports/top-creators", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestReportTopCreators_InvalidLimit() {
	r := chi.NewRouter()
	r.Get("/reports/top-creators", s.api.ReportTopCreators)
	req := httptest.NewRequest(http.MethodGet, "/reports/top-creators?limit=0", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
	s.reportService.AssertNotCalled(s.T(), "TopCreatorsByTeam")
}

func (s *APISuite) TestReportTopCreators_InternalError() {
	s.reportService.On("TopCreatorsByTeam", mock.Anything, uuid.MustParse(testUserID), mock.AnythingOfType("time.Time"), 3).
		Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/reports/top-creators", s.api.ReportTopCreators)
	req := httptest.NewRequest(http.MethodGet, "/reports/top-creators?limit=3", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.reportService.AssertExpectations(s.T())
}
