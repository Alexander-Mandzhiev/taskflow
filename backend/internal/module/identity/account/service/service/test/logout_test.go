package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestLogout_Success() {
	sessionID := uuid.New()
	s.sessionRepo.On("Delete", mock.Anything, sessionID).Return(nil).Once()

	err := s.svc.Logout(s.ctx, sessionID)

	assert.NoError(s.T(), err)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogout_DeleteError() {
	sessionID := uuid.New()
	s.sessionRepo.On("Delete", mock.Anything, sessionID).Return(assert.AnError).Once()

	err := s.svc.Logout(s.ctx, sessionID)

	assert.Error(s.T(), err)
	s.sessionRepo.AssertExpectations(s.T())
}
