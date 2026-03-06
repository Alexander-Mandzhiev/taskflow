package adapter_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

func (s *AdapterSuite) TestCreate_Success_WithoutRegistry() {
	input := &model.UserInput{Email: "u@ex.com", Name: "User"}
	hash := "hashed"
	user := &model.User{ID: uuid.New(), Email: input.Email, Name: input.Name}

	s.writer.On("Create", mock.Anything, mock.Anything, input, hash).
		Return(user, nil).Once()

	got, err := s.repo.Create(s.ctx, nil, input, hash)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Set")
}

func (s *AdapterSuite) TestCreate_Success_WithRegistry_RegistersHook() {
	registry := txmanager.NewHookRegistry()
	ctx := txmanager.WithHookRegistry(s.ctx, registry)
	input := &model.UserInput{Email: "u@ex.com", Name: "User"}
	hash := "hashed"
	user := &model.User{ID: uuid.New(), Email: input.Email, Name: input.Name}

	s.writer.On("Create", mock.Anything, mock.Anything, input, hash).
		Return(user, nil).Once()
	s.cache.On("Set", mock.Anything, user.ID.String(), user).
		Return(nil).Once()

	got, err := s.repo.Create(ctx, nil, input, hash)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)

	for _, hook := range registry.GetHooks() {
		assert.NoError(s.T(), hook(ctx))
	}
	s.writer.AssertExpectations(s.T())
	s.cache.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WriterError() {
	input := &model.UserInput{Email: "u@ex.com", Name: "User"}
	s.writer.On("Create", mock.Anything, mock.Anything, input, mock.AnythingOfType("string")).
		Return((*model.User)(nil), assert.AnError).Once()

	got, err := s.repo.Create(s.ctx, nil, input, "hash")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.writer.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Set")
}
