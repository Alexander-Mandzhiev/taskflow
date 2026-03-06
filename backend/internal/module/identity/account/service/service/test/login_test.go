package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestLogin_Success() {
	hash, _ := s.hasher.Hash("password123")
	user := &usermodel.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: hash}

	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(user, nil).Once()
	s.sessionRepo.On("Set", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.MatchedBy(func(sess *accountmodel.Session) bool {
		return sess != nil && sess.UserID == user.ID
	}), mock.AnythingOfType("time.Duration")).
		Return(nil).Once()

	sessionID, err := s.svc.Login(s.ctx, "user@example.com", "password123", "Mozilla/5.0", "192.168.1.1")

	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, sessionID)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_UserNotFound() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "nonexistent@example.com").
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()

	sessionID, err := s.svc.Login(s.ctx, "nonexistent@example.com", "password123", "", "")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidCredentials)
	assert.Equal(s.T(), uuid.Nil, sessionID)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_InvalidPassword() {
	hash, _ := s.hasher.Hash("correct")
	user := &usermodel.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: hash}

	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(user, nil).Once()

	sessionID, err := s.svc.Login(s.ctx, "user@example.com", "wrongpassword", "", "")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidCredentials)
	assert.Equal(s.T(), uuid.Nil, sessionID)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_SetSessionError() {
	hash, _ := s.hasher.Hash("password123")
	user := &usermodel.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: hash}

	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(user, nil).Once()
	s.sessionRepo.On("Set", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.AnythingOfType("time.Duration")).
		Return(assert.AnError).Once()

	sessionID, err := s.svc.Login(s.ctx, "user@example.com", "password123", "", "")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), uuid.Nil, sessionID)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}
