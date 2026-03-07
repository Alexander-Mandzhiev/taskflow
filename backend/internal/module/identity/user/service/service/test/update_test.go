package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestUpdate_NilInput() {
	got, err := s.svc.Update(s.ctx, uuid.New().String(), nil)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNilInput)
	assert.Nil(s.T(), got)
	s.repo.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdate_Success() {
	id := uuid.New().String()
	input := &model.UserInput{Email: "new@ex.com", Name: "New"}
	want := &model.User{ID: uuid.MustParse(id), Email: input.Email, Name: input.Name}

	s.repo.On("Update", mock.Anything, mock.Anything, id, input).
		Return(want, nil).Once()

	got, err := s.svc.Update(s.ctx, id, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestUpdate_RepoError() {
	id := uuid.New().String()
	input := &model.UserInput{Email: "u@ex.com", Name: "User"}

	s.repo.On("Update", mock.Anything, mock.Anything, id, input).
		Return((*model.User)(nil), assert.AnError).Once()

	got, err := s.svc.Update(s.ctx, id, input)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.repo.AssertExpectations(s.T())
}
