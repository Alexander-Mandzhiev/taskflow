package user_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestChangePassword_Success() {
	id := uuid.New().String()
	hash := "new_hash"

	s.repo.On("UpdatePasswordHash", mock.Anything, mock.Anything, id, hash).
		Return(nil).Once()

	err := s.svc.ChangePassword(s.ctx, id, hash)

	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestChangePassword_RepoError() {
	id := uuid.New().String()
	hash := "new_hash"

	s.repo.On("UpdatePasswordHash", mock.Anything, mock.Anything, id, hash).
		Return(model.ErrUserNotFound).Once()

	err := s.svc.ChangePassword(s.ctx, id, hash)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	s.repo.AssertExpectations(s.T())
}
