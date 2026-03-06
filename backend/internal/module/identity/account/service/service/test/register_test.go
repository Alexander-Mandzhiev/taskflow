package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestRegister_Success() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "newuser@example.com").
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.MatchedBy(func(in *usermodel.UserInput) bool {
		return in != nil && in.Email == "newuser@example.com" && in.Name == "New User"
	}), mock.AnythingOfType("string")).
		Return(&usermodel.User{}, nil).Once()

	err := s.svc.Register(s.ctx, "newuser@example.com", "password123", "New User")

	assert.NoError(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_EmailDuplicate() {
	existing := &usermodel.User{Email: "existing@example.com"}
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "existing@example.com").
		Return(existing, nil).Once()

	err := s.svc.Register(s.ctx, "existing@example.com", "password123", "User")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, usermodel.ErrEmailDuplicate)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_CreateReturnsDuplicate() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "race@example.com").
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).
		Return((*usermodel.User)(nil), usermodel.ErrEmailDuplicate).Once()

	err := s.svc.Register(s.ctx, "race@example.com", "password123", "User")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, usermodel.ErrEmailDuplicate)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_CreateError() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).
		Return((*usermodel.User)(nil), assert.AnError).Once()

	err := s.svc.Register(s.ctx, "user@example.com", "password123", "User")

	assert.Error(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}
