package adapter_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

func (s *AdapterSuite) TestUpdate_Success_WithoutRegistry() {
	id := uuid.New().String()
	input := model.UserInput{Email: "new@ex.com", Name: "New"}
	user := model.User{ID: uuid.MustParse(id), Email: input.Email, Name: input.Name}

	s.writer.On("Update", mock.Anything, mock.Anything, id, input).
		Return(user, nil).Once()

	got, err := s.repo.Update(s.ctx, nil, id, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Set")
}

func (s *AdapterSuite) TestUpdate_Success_WithRegistry_RegistersHook() {
	registry := txmanager.NewHookRegistry()
	ctx := txmanager.WithHookRegistry(s.ctx, registry)
	id := uuid.New().String()
	input := model.UserInput{Email: "new@ex.com", Name: "New"}
	user := model.User{ID: uuid.MustParse(id), Email: input.Email, Name: input.Name}

	s.writer.On("Update", mock.Anything, mock.Anything, id, input).
		Return(user, nil).Once()
	s.cache.On("Set", mock.Anything, id, user).
		Return(nil).Once()

	got, err := s.repo.Update(ctx, nil, id, input)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)

	for _, hook := range registry.GetHooks() {
		assert.NoError(s.T(), hook(ctx))
	}
	s.writer.AssertExpectations(s.T())
	s.cache.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestUpdate_WriterError() {
	id := uuid.New().String()
	input := model.UserInput{Email: "u@ex.com", Name: "User"}
	s.writer.On("Update", mock.Anything, mock.Anything, id, input).
		Return(model.User{}, assert.AnError).Once()

	got, err := s.repo.Update(s.ctx, nil, id, input)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.User{}, got)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Set")
}
