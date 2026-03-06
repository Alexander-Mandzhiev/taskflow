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

	input := accountmodel.LoginInput{Email: "user@example.com", Password: "password123", UserAgent: "Mozilla/5.0", IP: "192.168.1.1"}
	accessToken, refreshToken, err := s.svc.Login(s.ctx, input)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken)
	assert.NotEmpty(s.T(), refreshToken)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_UserNotFound() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "nonexistent@example.com").
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()

	input := accountmodel.LoginInput{Email: "nonexistent@example.com", Password: "password123"}
	_, _, err := s.svc.Login(s.ctx, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidCredentials)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_InvalidPassword() {
	hash, _ := s.hasher.Hash("correct")
	user := &usermodel.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: hash}

	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(user, nil).Once()

	input := accountmodel.LoginInput{Email: "user@example.com", Password: "wrongpassword"}
	_, _, err := s.svc.Login(s.ctx, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, accountmodel.ErrInvalidCredentials)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLogin_SetSessionError() {
	hash, _ := s.hasher.Hash("password123")
	user := &usermodel.User{ID: uuid.New(), Email: "user@example.com", PasswordHash: hash}

	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(user, nil).Once()
	s.sessionRepo.On("Set", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.Anything, mock.AnythingOfType("time.Duration")).
		Return(assert.AnError).Once()

	input := accountmodel.LoginInput{Email: "user@example.com", Password: "password123"}
	_, _, err := s.svc.Login(s.ctx, input)

	assert.Error(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
	s.sessionRepo.AssertExpectations(s.T())
}
