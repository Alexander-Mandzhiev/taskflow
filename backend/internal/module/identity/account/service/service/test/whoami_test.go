package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
)

func (s *ServiceSuite) TestWhoami_Success() {
	sessionID := uuid.New()
	userID := uuid.New()
	session := &accountmodel.Session{UserID: userID}

	s.sessionRepo.On("Get", mock.Anything, sessionID).Return(session, nil).Once()

	got, err := s.svc.Whoami(s.ctx, sessionID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userID, got)
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoami_SessionNotFound() {
	sessionID := uuid.New()
	s.sessionRepo.On("Get", mock.Anything, sessionID).Return((*accountmodel.Session)(nil), accountmodel.ErrSessionNotFound).Once()

	got, err := s.svc.Whoami(s.ctx, sessionID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, accountmodel.ErrSessionNotFound)
	assert.Equal(s.T(), uuid.Nil, got)
	s.sessionRepo.AssertExpectations(s.T())
}
