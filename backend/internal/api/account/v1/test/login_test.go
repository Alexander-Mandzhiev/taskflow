package account_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

func containsCookie(setCookies []string, prefix string) bool {
	for _, c := range setCookies {
		if strings.HasPrefix(strings.TrimSpace(c), prefix) {
			return true
		}
	}
	return false
}

func (s *APISuite) TestLogin_Success() {
	accessToken := "access-jwt-token"
	refreshToken := "refresh-jwt-token"
	body, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "password123"})

	s.accountService.On("Login", mock.Anything, mock.MatchedBy(func(in accountmodel.LoginInput) bool {
		return in.Email == "user@example.com" && in.Password == "password123"
	})).Return(accessToken, refreshToken, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	setCookies := rec.Header().Values("Set-Cookie")
	assert.Len(s.T(), setCookies, 2, "expected access_token and refresh_token cookies")
	assert.True(s.T(), containsCookie(setCookies, "access_token="), "missing access_token cookie")
	assert.True(s.T(), containsCookie(setCookies, "refresh_token="), "missing refresh_token cookie")
	var resp map[string]interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp["success"].(bool))
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestLogin_InvalidCredentials() {
	body, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "wrongpassword"})

	s.accountService.On("Login", mock.Anything, mock.MatchedBy(func(in accountmodel.LoginInput) bool {
		return in.Email == "user@example.com" && in.Password == "wrongpassword"
	})).Return("", "", accountmodel.ErrInvalidCredentials).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusUnauthorized, rec.Code)
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestLogin_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestLogin_ValidationError_EmptyEmail() {
	body, _ := json.Marshal(map[string]string{"email": "", "password": "password123"})

	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestLogin_ValidationError_ShortPassword() {
	body, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "short"})

	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestLogin_InternalError() {
	body, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "password123"})

	s.accountService.On("Login", mock.Anything, mock.MatchedBy(func(in accountmodel.LoginInput) bool {
		return in.Email == "user@example.com" && in.Password == "password123"
	})).Return("", "", assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Login(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.accountService.AssertExpectations(s.T())
}
