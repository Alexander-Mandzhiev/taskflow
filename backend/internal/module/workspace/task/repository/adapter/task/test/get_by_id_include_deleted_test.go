package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestGetByIDIncludeDeleted_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	deletedAt := time.Now()
	task := model.Task{
		ID:        taskID,
		Title:     "Deleted Task",
		TeamID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		DeletedAt: &deletedAt,
	}

	s.taskReader.On("GetByIDIncludeDeleted", mock.Anything, (*sqlx.Tx)(nil), taskID).
		Return(task, nil).Once()

	got, err := s.repo.GetByIDIncludeDeleted(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	assert.NotNil(s.T(), got.DeletedAt)
	s.taskReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByIDIncludeDeleted_TaskNotFound() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	s.taskReader.On("GetByIDIncludeDeleted", mock.Anything, mock.Anything, taskID).
		Return(model.Task{}, model.ErrTaskNotFound).Once()

	got, err := s.repo.GetByIDIncludeDeleted(s.ctx, nil, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Equal(s.T(), model.Task{}, got)
	s.taskReader.AssertExpectations(s.T())
}
