package service_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/jwt"
)

func (s *ServiceSuite) TestLogout_EmptyToken_Success() {
	err := s.svc.Logout(s.ctx, "")

	assert.NoError(s.T(), err)
	s.sessionRepo.AssertNotCalled(s.T(), "Delete")
}

func (s *ServiceSuite) TestLogout_Success() {
	userID := uuid.New().String()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID, "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	s.sessionRepo.On("Delete", mock.Anything, jti).Return(nil).Once()

	err = s.svc.Logout(s.ctx, refreshToken)

	assert.NoError(s.T(), err)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogout_DeleteError() {
	userID := uuid.New().String()
	refreshToken, jti, err := jwt.GenerateRefreshToken(userID, "web", "test-refresh-secret", time.Hour)
	s.Require().NoError(err)

	s.sessionRepo.On("Delete", mock.Anything, jti).Return(assert.AnError).Once()

	err = s.svc.Logout(s.ctx, refreshToken)

	assert.Error(s.T(), err)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogout_InvalidToken() {
	err := s.svc.Logout(s.ctx, "invalid_token")

	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidRefreshToken)
	s.sessionRepo.AssertNotCalled(s.T(), "Delete")
}
