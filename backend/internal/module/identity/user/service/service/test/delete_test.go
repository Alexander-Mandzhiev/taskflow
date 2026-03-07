package service_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *ServiceSuite) TestDelete_Success() {
	id := uuid.New().String()

	s.repo.On("Delete", mock.Anything, mock.Anything, id).
		Return(nil).Once()

	err := s.svc.Delete(s.ctx, id)

	assert.NoError(s.T(), err)
	s.repo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestDelete_RepoError() {
	id := uuid.New().String()

	s.repo.On("Delete", mock.Anything, mock.Anything, id).
		Return(model.ErrUserNotFound).Once()

	err := s.svc.Delete(s.ctx, id)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	s.repo.AssertExpectations(s.T())
}
