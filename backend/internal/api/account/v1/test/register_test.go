package account_v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *APISuite) TestRegister_Success() {
	body, _ := json.Marshal(map[string]string{
		"email": "newuser@example.com", "password": "password123", "name": "New User",
	})

	s.accountService.On("Register", mock.Anything, "newuser@example.com", "password123", "New User").
		Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var resp map[string]interface{}
	assert.NoError(s.T(), json.NewDecoder(rec.Body).Decode(&resp))
	assert.True(s.T(), resp["success"].(bool))
	assert.Equal(s.T(), "Пользователь успешно зарегистрирован", resp["message"])
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestRegister_EmailDuplicate() {
	body, _ := json.Marshal(map[string]string{
		"email": "existing@example.com", "password": "password123", "name": "User",
	})

	s.accountService.On("Register", mock.Anything, "existing@example.com", "password123", "User").
		Return(model.ErrEmailDuplicate).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusConflict, rec.Code)
	s.accountService.AssertExpectations(s.T())
}

func (s *APISuite) TestRegister_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestRegister_ValidationError_InvalidEmail() {
	body, _ := json.Marshal(map[string]string{
		"email": "not-an-email", "password": "password123", "name": "User",
	})

	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestRegister_ValidationError_ShortPassword() {
	body, _ := json.Marshal(map[string]string{
		"email": "user@example.com", "password": "short", "name": "User",
	})

	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *APISuite) TestRegister_InternalError() {
	body, _ := json.Marshal(map[string]string{
		"email": "user@example.com", "password": "password123", "name": "User",
	})

	s.accountService.On("Register", mock.Anything, "user@example.com", "password123", "User").
		Return(assert.AnError).Once()

	req := httptest.NewRequest(http.MethodPost, "/account/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.api.Register(rec, req)

	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
	s.accountService.AssertExpectations(s.T())
}
