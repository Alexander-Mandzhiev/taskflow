package team_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const testOwnerUserID = "550e8400-e29b-41d4-a716-446655440000"

func (s *APISuite) TestCreate_Success() {
	body, _ := json.Marshal(map[string]string{"name": "My Team"})

	ownerID := uuid.MustParse(testOwnerUserID)
	createdTeam := &model2.Team{
		ID:        uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		Name:      "My Team",
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamService.On("Create", mock.Anything, mock.MatchedBy(func(in *model2.TeamInput) bool {
		return in != nil && in.Name == "My Team"
	}), ownerID).Return(createdTeam, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.Create(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedBy string `json:"created_by"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(s.T(), "My Team", resp.Name)
	assert.Equal(s.T(), testOwnerUserID, resp.CreatedBy)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreate_NoAuth() {
	body, _ := json.Marshal(map[string]string{"name": "My Team"})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Create(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestCreate_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.Create(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestCreate_ValidationError_EmptyName() {
	body, _ := json.Marshal(map[string]string{"name": ""})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.Create(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestCreate_InternalError() {
	body, _ := json.Marshal(map[string]string{"name": "My Team"})

	s.teamService.On("Create", mock.Anything, mock.MatchedBy(func(in *model2.TeamInput) bool {
		return in != nil && in.Name == "My Team"
	}), uuid.MustParse(testOwnerUserID)).Return(nil, assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	s.api.Create(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.teamService.AssertExpectations(s.T())
}
