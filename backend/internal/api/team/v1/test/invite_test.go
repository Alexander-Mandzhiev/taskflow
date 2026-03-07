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

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const invitedEmail = "invited@example.com"
const invitedUserID = "770e8400-e29b-41d4-a716-446655440002"

func (s *APISuite) TestInvite_Success() {
	body, _ := json.Marshal(map[string]string{
		"email": invitedEmail,
		"role":  "member",
	})

	invitedMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(invitedUserID),
		TeamID:    uuid.MustParse(testTeamID),
		Role:      model.RoleMember,
		CreatedAt: time.Now(),
	}

	s.teamService.On("InviteByEmail", mock.Anything, testTeamID, testOwnerUserID, invitedEmail, "member").Return(invitedMember, nil).Once()

	r := chi.NewRouter()
	r.Post("/teams/{id}/invite", s.api.Invite)
	req := httptest.NewRequest(http.MethodPost, "/teams/"+testTeamID+"/invite", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(metadata.SetUserID(req.Context(), testOwnerUserID))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Member  struct {
			UserID string `json:"user_id"`
			Role   string `json:"role"`
		} `json:"member"`
	}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp.Success)
	assert.Equal(s.T(), "Пользователь приглашён в команду", resp.Message)
	assert.Equal(s.T(), invitedUserID, resp.Member.UserID)
	assert.Equal(s.T(), "member", resp.Member.Role)
	s.teamService.AssertExpectations(s.T())
}

func (s *APISuite) TestInvite_Forbidden_NotOwnerOrAdmin() {
	body, _ := json.Marshal(map[string]string{"email": invitedEmail, "role": "member"})

	s.teamService.On("InviteByEmail", mock.Anything, testTeamID, testOwnerUserID, invitedEmail, "member").Return(nil, model.ErrForbidden).Once()

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

	s.teamService.On("InviteByEmail", mock.Anything, testTeamID, testOwnerUserID, invitedEmail, "member").Return(nil, model.ErrMemberNotFound).Once()

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

	s.teamService.On("InviteByEmail", mock.Anything, testTeamID, testOwnerUserID, invitedEmail, "member").Return(nil, model.ErrAlreadyMember).Once()

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

	s.teamService.On("InviteByEmail", mock.Anything, testTeamID, testOwnerUserID, invitedEmail, "member").Return(nil, assert.AnError).Once()

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
