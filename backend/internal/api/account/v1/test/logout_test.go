package account_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

func (s *APISuite) TestLogout_SuccessWithSession() {
	sessionID := uuid.New()

	s.accountService.On("Logout", mock.Anything, sessionID).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req = req.WithContext(metadata.SetSessionIDUUID(req.Context(), sessionID))
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp["success"].(bool))
	assert.Equal(s.T(), "Сессия завершена", resp["message"])
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestLogout_SuccessNoSession_ClearsCookie() {
	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp["success"].(bool))
}

func (s *APISuite) TestLogout_SessionNotFound() {
	sessionID := uuid.New()

	s.accountService.On("Logout", mock.Anything, sessionID).Return(accountmodel.ErrSessionNotFound).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req = req.WithContext(metadata.SetSessionIDUUID(req.Context(), sessionID))
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestLogout_InternalError() {
	sessionID := uuid.New()

	s.accountService.On("Logout", mock.Anything, sessionID).Return(assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req = req.WithContext(metadata.SetSessionIDUUID(req.Context(), sessionID))
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.accountService.AssertExpectations(s.T())
}
