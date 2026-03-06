package adapter_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

func (s *AdapterSuite) TestUpdatePasswordHash_Success_WithoutRegistry() {
	id := uuid.New().String()
	hash := "new_hash"
	s.writer.On("UpdatePasswordHash", mock.Anything, mock.Anything, id, hash).
		Return(nil).Once()

	err := s.repo.UpdatePasswordHash(s.ctx, nil, id, hash)

	assert.NoError(s.T(), err)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Delete")
}

func (s *AdapterSuite) TestUpdatePasswordHash_Success_WithRegistry_RegistersHook() {
	registry := txmanager.NewHookRegistry()
	ctx := txmanager.WithHookRegistry(s.ctx, registry)
	id := uuid.New().String()
	hash := "new_hash"

	s.writer.On("UpdatePasswordHash", mock.Anything, mock.Anything, id, hash).
		Return(nil).Once()
	s.cache.On("Delete", mock.Anything, id).
		Return(nil).Once()

	err := s.repo.UpdatePasswordHash(ctx, nil, id, hash)
	assert.NoError(s.T(), err)

	for _, hook := range registry.GetHooks() {
		assert.NoError(s.T(), hook(ctx))
	}
	s.writer.AssertExpectations(s.T())
	s.cache.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestUpdatePasswordHash_WriterError() {
	id := uuid.New().String()
	hash := "new_hash"
	s.writer.On("UpdatePasswordHash", mock.Anything, mock.Anything, id, hash).
		Return(model.ErrUserNotFound).Once()

	err := s.repo.UpdatePasswordHash(s.ctx, nil, id, hash)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Delete")
}
