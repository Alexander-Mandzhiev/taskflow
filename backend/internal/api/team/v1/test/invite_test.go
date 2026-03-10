package team_v1_test

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

	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const invitedEmail = "invited@example.com"

func (s *APISuite) TestInvite_Success() {
	body, _ := json.Marshal(map[string]string{
		"email": invitedEmail,
		"role":  "member",
	})

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	invitation := teamModel.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    uuid.MustParse(testTeamID),
		Email:     invitedEmail,
		Role:      "member",
		ExpiresAt: expiresAt,
	}

	s.teamService.On("InviteByEmail", mock.Anything, uuid.MustParse(testTeamID), uuid.MustParse(testOwnerUserID), invitedEmail, "member").Return(invitation, nil).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp struct {
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		Invitation struct {
			ID        string `json:"id"`
			TeamID    string `json:"team_id"`
			Email     string `json:"email"`
			Role      string `json:"role"`
			ExpiresAt string `json:"expires_at"`
		} `json:"invitation"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp.Success)
	assert.Equal(s.T(), "Приглашение создано", resp.Message)
	assert.Equal(s.T(), invitation.ID.String(), resp.Invitation.ID)
	assert.Equal(s.T(), testTeamID, resp.Invitation.TeamID)
	assert.Equal(s.T(), invitedEmail, resp.Invitation.Email)
	assert.Equal(s.T(), "member", resp.Invitation.Role)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestInvite_Forbidden_NotOwnerOrAdmin() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	s.teamService.On("InviteByEmail", mock.Anything, uuid.MustParse(testTeamID), uuid.MustParse(testOwnerUserID), invitedEmail, "member").Return(teamModel.TeamInvitation{}, teamModel.ErrForbidden).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusForbidden, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestInvite_MemberNotFound() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	s.teamService.On("InviteByEmail", mock.Anything, uuid.MustParse(testTeamID), uuid.MustParse(testOwnerUserID), invitedEmail, "member").Return(teamModel.TeamInvitation{}, teamModel.ErrMemberNotFound).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestInvite_AlreadyMember() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	s.teamService.On("InviteByEmail", mock.Anything, uuid.MustParse(testTeamID), uuid.MustParse(testOwnerUserID), invitedEmail, "member").Return(teamModel.TeamInvitation{}, teamModel.ErrAlreadyMember).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestInvite_NoAuth() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *APISuite) TestInvite_InvalidJSON() {
	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestInvite_ValidationError_InvalidRole() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "invalid_role"})

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestInvite_ValidationError_MissingEmail() {
	body, _ := json.Marshal(map[string]string{"role": "member"})

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestInvite_InternalError() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	s.teamService.On("InviteByEmail", mock.Anything, uuid.MustParse(testTeamID), uuid.MustParse(testOwnerUserID), invitedEmail, "member").Return(teamModel.TeamInvitation{}, assert.AnError).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.teamService.AssertExpectations(s.T())
}
