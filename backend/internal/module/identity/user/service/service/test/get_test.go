package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestGetByID_Success() {
	id := uuid.New().String()
	want := &model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.repo.On("GetByID", mock.Anything, mock.Anything, id).
		Return(want, nil).Once()

	got, err := s.svc.GetByID(s.ctx, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_Error() {
	id := uuid.New().String()

	s.repo.On("GetByID", mock.Anything, mock.Anything, id).
		Return((*model.User)(nil), model.ErrUserNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, id)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	assert.Nil(s.T(), got)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByEmail_Success() {
	email := "u@ex.com"
	want := &model.User{ID: uuid.New(), Email: email, Name: "User"}

	s.repo.On("GetByEmail", mock.Anything, mock.Anything, email).
		Return(want, nil).Once()

	got, err := s.svc.GetByEmail(s.ctx, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByEmail_Error() {
	email := "missing@ex.com"

	s.repo.On("GetByEmail", mock.Anything, mock.Anything, email).
		Return((*model.User)(nil), model.ErrUserNotFound).Once()

	got, err := s.svc.GetByEmail(s.ctx, email)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	assert.Nil(s.T(), got)
	s.repo.AssertExpectations(s.T())
}
