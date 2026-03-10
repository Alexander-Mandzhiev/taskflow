package user_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestCreate_NilInput() {
	got, err := s.svc.Create(s.ctx, model.UserInput{}, "hash")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNilInput)
	assert.Equal(s.T(), model.User{}, got)
	s.repo.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreate_Success() {
	input := model.UserInput{Email: "u@ex.com", Name: "User"}
	hash := "hash"
	want := model.User{ID: uuid.New(), Email: input.Email, Name: input.Name}

	s.repo.On("Create", mock.Anything, mock.Anything, input, hash).
		Return(want, nil).Once()

	got, err := s.svc.Create(s.ctx, input, hash)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_RepoError() {
	input := model.UserInput{Email: "u@ex.com", Name: "User"}
	hash := "hash"

	s.repo.On("Create", mock.Anything, mock.Anything, input, hash).
		Return(model.User{}, assert.AnError).Once()

	got, err := s.svc.Create(s.ctx, input, hash)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.User{}, got)
	s.repo.AssertExpectations(s.T())
}
