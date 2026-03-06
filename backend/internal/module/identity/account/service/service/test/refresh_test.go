package service_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
)

func (s *ServiceSuite) TestRefresh_EmptyToken() {
	accessToken, userID, err := s.svc.Refresh(s.ctx, "", "Mozilla/5.0", "192.168.1.1")

	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidRefreshToken)
	assert.Empty(s.T(), accessToken)
	assert.Equal(s.T(), uuid.Nil, userID)
	s.sessionRepo.AssertNotCalled(s.T(), "Get")
}

func (s *ServiceSuite) TestRefresh_InvalidToken() {
	accessToken, userID, err := s.svc.Refresh(s.ctx, "invalid_token", "Mozilla/5.0", "192.168.1.1")

	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidRefreshToken)
	assert.Empty(s.T(), accessToken)
	assert.Equal(s.T(), uuid.Nil, userID)
	s.sessionRepo.AssertNotCalled(s.T(), "Get")
}

func (s *ServiceSuite) TestRefresh_SessionNotFound() {
	userID := uuid.New()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID.String(), "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	s.sessionRepo.On("Get", mock.Anything, jti).Return((*accountmodel.Session)(nil), accountmodel.ErrSessionNotFound).Once()

	accessToken, gotUserID, err := s.svc.Refresh(s.ctx, refreshToken, "Mozilla/5.0", "192.168.1.1")

	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidRefreshToken)
	assert.Empty(s.T(), accessToken)
	assert.Equal(s.T(), uuid.Nil, gotUserID)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRefresh_GetSessionError() {
	userID := uuid.New()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID.String(), "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	s.sessionRepo.On("Get", mock.Anything, jti).Return((*accountmodel.Session)(nil), assert.AnError).Once()

	accessToken, gotUserID, err := s.svc.Refresh(s.ctx, refreshToken, "Mozilla/5.0", "192.168.1.1")

	assert.Error(s.T(), err)
	assert.Empty(s.T(), accessToken)
	assert.Equal(s.T(), uuid.Nil, gotUserID)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRefresh_SessionMetadataMismatch() {
	userID := uuid.New()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID.String(), "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	// Сессия с другим типом устройства (mobile vs desktop в запросе)
	session := &accountmodel.Session{
		UserID:     userID,
		DeviceType: "mobile",
		UserAgent:  "Mobile Safari",
		IP:         "192.168.1.1",
	}
	s.sessionRepo.On("Get", mock.Anything, jti).Return(session, nil).Once()

	accessToken, gotUserID, err := s.svc.Refresh(s.ctx, refreshToken, "Mozilla/5.0 (Desktop)", "192.168.1.1")

	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidRefreshToken)
	assert.Empty(s.T(), accessToken)
	assert.Equal(s.T(), uuid.Nil, gotUserID)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRefresh_Success() {
	userID := uuid.New()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID.String(), "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	session := &accountmodel.Session{
		UserID:     userID,
		DeviceType: "desktop",
		UserAgent:  "Mozilla/5.0",
		IP:         "192.168.1.1",
	}
	s.sessionRepo.On("Get", mock.Anything, jti).Return(session, nil).Once()

	accessToken, gotUserID, err := s.svc.Refresh(s.ctx, refreshToken, "Mozilla/5.0", "192.168.1.1")

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken)
	assert.Equal(s.T(), userID, gotUserID)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRefresh_Success_EmptySessionMetadata() {
	// Сессия без UserAgent/IP — сверка не применяется, refresh проходит
	userID := uuid.New()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID.String(), "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	session := &accountmodel.Session{
		UserID:     userID,
		DeviceType: "desktop",
		UserAgent:  "",
		IP:         "",
	}
	s.sessionRepo.On("Get", mock.Anything, jti).Return(session, nil).Once()

	accessToken, gotUserID, err := s.svc.Refresh(s.ctx, refreshToken, "Mozilla/5.0", "192.168.1.1")

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken)
	assert.Equal(s.T(), userID, gotUserID)
	s.sessionRepo.AssertExpectations(s.T())
}
