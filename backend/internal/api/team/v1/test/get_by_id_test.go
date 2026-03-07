package team_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

const testTeamID = "550e8400-e29b-41d4-a716-446655440001"

func (s *APISuite) TestGetByID_Success() {
	teamID := uuid.MustParse(testTeamID)
	teamWithMembers := &model.TeamWithMembers{
		Team: model.Team{
			ID:        teamID,
			Name:      "My Team",
			CreatedBy: uuid.MustParse(testOwnerUserID),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Members: []*model.TeamMember{},
	}

	s.teamService.On("GetByID", mock.Anything, testTeamID).Return(teamWithMembers, nil).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp struct {
		Team struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"team"`
		Members []interface{} `json:"members"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), testTeamID, resp.Team.ID)
	assert.Equal(s.T(), "My Team", resp.Team.Name)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_NotFound() {
	s.teamService.On("GetByID", mock.Anything, testTeamID).Return(nil, model.ErrTeamNotFound).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_InternalError() {
	s.teamService.On("GetByID", mock.Anything, testTeamID).Return(nil, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.teamService.AssertExpectations(s.T())
}
