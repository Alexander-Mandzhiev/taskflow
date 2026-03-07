package task_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestRestore_Success() {
	tx := &sqlx.Tx{}
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")

	s.taskWriter.On("Restore", mock.Anything, tx, taskID).
		Return(nil).Once()

	err := s.repo.Restore(s.ctx, tx, taskID)

	assert.NoError(s.T(), err)
	s.taskWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestRestore_TaskNotFound() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	s.taskWriter.On("Restore", mock.Anything, mock.Anything, taskID).
		Return(model.ErrTaskNotFound).Once()

	err := s.repo.Restore(s.ctx, nil, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	s.taskWriter.AssertExpectations(s.T())
}
