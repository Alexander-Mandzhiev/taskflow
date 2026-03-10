package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	accountmodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestRegister_Success() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "newuser@example.com").
		Return(usermodel.User{}, usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.MatchedBy(func(in usermodel.UserInput) bool {
		return in.Email == "newuser@example.com" && in.Name == "New User"
	}), mock.AnythingOfType("string")).
		Return(usermodel.User{}, nil).Once()

	err := s.svc.Register(s.ctx, accountmodel.RegisterInput{Email: "newuser@example.com", Password: "password123", Name: "New User"})

	assert.NoError(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_EmailDuplicate() {
	existing := usermodel.User{ID: uuid.New(), Email: "existing@example.com"}
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "existing@example.com").
		Return(existing, nil).Once()

	err := s.svc.Register(s.ctx, accountmodel.RegisterInput{Email: "existing@example.com", Password: "password123", Name: "User"})

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, usermodel.ErrEmailDuplicate)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_CreateReturnsDuplicate() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "race@example.com").
		Return(usermodel.User{}, usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).
		Return(usermodel.User{}, usermodel.ErrEmailDuplicate).Once()

	err := s.svc.Register(s.ctx, accountmodel.RegisterInput{Email: "race@example.com", Password: "password123", Name: "User"})

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, usermodel.ErrEmailDuplicate)
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegister_CreateError() {
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, "user@example.com").
		Return(usermodel.User{}, usermodel.ErrUserNotFound).Once()
	s.userRepo.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("string")).
		Return(usermodel.User{}, assert.AnError).Once()

	err := s.svc.Register(s.ctx, accountmodel.RegisterInput{Email: "user@example.com", Password: "password123", Name: "User"})

	assert.Error(s.T(), err)
	s.userRepo.AssertExpectations(s.T())
}
