package adapter_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *AdapterSuite) TestGetByEmail_Success() {
	email := "u@ex.com"
	user := &model.User{Email: email, Name: "User"}

	s.reader.On("GetByEmail", mock.Anything, mock.Anything, email).
		Return(user, nil).Once()

	got, err := s.repo.GetByEmail(s.ctx, nil, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.reader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByEmail_NotFound() {
	email := "missing@ex.com"
	s.reader.On("GetByEmail", mock.Anything, mock.Anything, email).
		Return((*model.User)(nil), model.ErrUserNotFound).Once()

	got, err := s.repo.GetByEmail(s.ctx, nil, email)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	assert.Nil(s.T(), got)
	s.reader.AssertExpectations(s.T())
}
