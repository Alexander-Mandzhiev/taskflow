package account_v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

const testRefreshTokenValue = "test_refresh_token_value"

func (s *APISuite) TestLogout_SuccessWithSession() {
	s.accountService.On("Logout", mock.Anything, testRefreshTokenValue).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshTokenValue})
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
	s.accountService.On("Logout", mock.Anything, "").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var resp map[string]interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp["success"].(bool))
}

func (s *APISuite) TestLogout_SessionNotFound() {
	s.accountService.On("Logout", mock.Anything, testRefreshTokenValue).Return(accountmodel.ErrSessionNotFound).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshTokenValue})
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestLogout_InternalError() {
	s.accountService.On("Logout", mock.Anything, testRefreshTokenValue).Return(assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: testRefreshTokenValue})
	rec := httptest.NewRecorder()

	s.api.Logout(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.accountService.AssertExpectations(s.T())
}
