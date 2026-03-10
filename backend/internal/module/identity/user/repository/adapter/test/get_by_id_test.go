package adapter_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

func (s *AdapterSuite) TestGetByID_WithTx_ReadsFromDB() {
	tx := &sqlx.Tx{}
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.reader.On("GetByID", mock.Anything, tx, id).
		Return(user, nil).Once()

	got, err := s.repo.GetByID(s.ctx, tx, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.reader.AssertExpectations(s.T())
	s.cache.AssertNotCalled(s.T(), "Get")
	s.cache.AssertNotCalled(s.T(), "Set")
}

func (s *AdapterSuite) TestGetByID_NoTx_CacheHit() {
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.cache.On("Get", mock.Anything, id).
		Return(user, true, nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertNotCalled(s.T(), "GetByID")
}

func (s *AdapterSuite) TestGetByID_NoTx_CacheHit_IgnoresCacheError() {
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	// Даже если Redis вернул ошибку, но значение есть — считаем кеш-хитом.
	s.cache.On("Get", mock.Anything, id).
		Return(user, true, assert.AnError).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertNotCalled(s.T(), "GetByID")
	s.cache.AssertNotCalled(s.T(), "Set")
}

func (s *AdapterSuite) TestGetByID_NoTx_CacheMiss_FallbackToDB() {
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.cache.On("Get", mock.Anything, id).
		Return(model.User{}, false, nil).Once()
	s.reader.On("GetByID", mock.Anything, mock.Anything, id).
		Return(user, nil).Once()
	s.cache.On("Set", mock.Anything, id, user).
		Return(nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_NoTx_CacheMiss_CacheSetError_Ignored() {
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.cache.On("Get", mock.Anything, id).
		Return(model.User{}, false, nil).Once()
	s.reader.On("GetByID", mock.Anything, mock.Anything, id).
		Return(user, nil).Once()
	s.cache.On("Set", mock.Anything, id, user).
		Return(assert.AnError).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_NoTx_CacheError_FallbackToDB() {
	id := uuid.New().String()
	user := model.User{ID: uuid.MustParse(id), Email: "u@ex.com", Name: "User"}

	s.cache.On("Get", mock.Anything, id).
		Return(model.User{}, false, assert.AnError).Once()
	s.reader.On("GetByID", mock.Anything, mock.Anything, id).
		Return(user, nil).Once()
	s.cache.On("Set", mock.Anything, id, user).
		Return(nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_NoTx_DBError() {
	id := uuid.New().String()
	s.cache.On("Get", mock.Anything, id).
		Return(model.User{}, false, nil).Once()
	s.reader.On("GetByID", mock.Anything, mock.Anything, id).
		Return(model.User{}, model.ErrUserNotFound).Once()

	got, err := s.repo.GetByID(s.ctx, nil, id)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrUserNotFound)
	assert.Equal(s.T(), model.User{}, got)
	s.cache.AssertExpectations(s.T())
	s.reader.AssertExpectations(s.T())
}
