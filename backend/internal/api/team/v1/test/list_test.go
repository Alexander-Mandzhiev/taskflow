package team_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

func (s *APISuite) TestList_Success() {
	teamID := uuid.MustParse(testTeamID)
	teams := []*model.TeamWithRole{
		{
			Team: model.Team{
				ID:        teamID,
				Name:      "My Team",
				CreatedBy: uuid.MustParse(testOwnerUserID),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Role: model.RoleOwner,
		},
	}

	s.teamService.On("ListByUserID", mock.Anything, uuid.MustParse(testOwnerUserID)).Return(teams, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.List(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Len(s.T(), resp, 1)
	assert.Equal(s.T(), "My Team", resp[0].Name)
	assert.Equal(s.T(), model.RoleOwner, resp[0].Role)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestList_Empty() {
	s.teamService.On("ListByUserID", mock.Anything, uuid.MustParse(testOwnerUserID)).Return([]*model.TeamWithRole{}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.List(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp []interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Empty(s.T(), resp)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestList_NoAuth() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams", nil)
	rec := httptest.NewRecorder()

	s.api.List(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestList_InternalError() {
	s.teamService.On("ListByUserID", mock.Anything, uuid.MustParse(testOwnerUserID)).Return(nil, assert.AnError).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/teams", nil)
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.List(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.teamService.AssertExpectations(s.T())
}
