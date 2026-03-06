package adapter_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

func (s *AdapterSuite) TestDelete_Success_WithoutRegistry() {
	id := uuid.New().String()
	s.writer.On("Delete", mock.Anything, mock.Anything, id).
		Return(nil).Once()

	err := s.repo.Delete(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Delete")
}

func (s *AdapterSuite) TestDelete_Success_WithRegistry_RegistersHook() {
	registry := txmanager.NewHookRegistry()
	ctx := txmanager.WithHookRegistry(s.ctx, registry)
	id := uuid.New().String()

	s.writer.On("Delete", mock.Anything, mock.Anything, id).
		Return(nil).Once()
	s.cache.On("Delete", mock.Anything, id).
		Return(nil).Once()

	err := s.repo.Delete(ctx, nil, id)
	assert.NoError(s.T(), err)

	for _, hook := range registry.GetHooks() {
		assert.NoError(s.T(), hook(ctx))
	}
	s.writer.AssertExpectations(s.T())
	s.cache.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestDelete_WriterError() {
	id := uuid.New().String()
	s.writer.On("Delete", mock.Anything, mock.Anything, id).
		Return(model.ErrUserNotFound).Once()

	err := s.repo.Delete(s.ctx, nil, id)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Delete")
}
