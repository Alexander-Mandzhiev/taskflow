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

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const testTeamID = "550e8400-e29b-41d4-a716-446655440001"

func (s *APISuite) TestGetByID_Success() {
	teamID := uuid.MustParse(testTeamID)
	userID := uuid.MustParse(testOwnerUserID)
	teamWithMembers := model2.TeamWithMembers{
		Team: model2.Team{
			ID:        teamID,
			Name:      "My Team",
			CreatedBy: userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Members: []model2.TeamMember{},
	}

	s.teamService.On("GetByID", mock.Anything, teamID, userID).Return(teamWithMembers, nil).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
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

func (s *APISuite) TestGetByID_NoAuth() {
	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestGetByID_Forbidden() {
	teamID := uuid.MustParse(testTeamID)
	userID := uuid.MustParse(testOwnerUserID)
	s.teamService.On("GetByID", mock.Anything, teamID, userID).Return(model2.TeamWithMembers{}, model2.ErrForbidden).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_NotFound() {
	teamID := uuid.MustParse(testTeamID)
	userID := uuid.MustParse(testOwnerUserID)
	s.teamService.On("GetByID", mock.Anything, teamID, userID).Return(model2.TeamWithMembers{}, model2.ErrTeamNotFound).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetByID_InternalError() {
	teamID := uuid.MustParse(testTeamID)
	userID := uuid.MustParse(testOwnerUserID)
	s.teamService.On("GetByID", mock.Anything, teamID, userID).Return(model2.TeamWithMembers{}, assert.AnError).Once()

	r := chi.NewRouter()
	r.Get("/teams/{id}", s.api.GetByID)
	req := httptest.NewRequest(http.MethodGet, "/teams/"+testTeamID, nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.teamService.AssertExpectations(s.T())
}
