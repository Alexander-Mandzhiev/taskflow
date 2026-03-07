package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestGetByID_Success() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	task := &model.Task{
		ID:        taskID,
		Title:     "Task",
		Status:    model.TaskStatusTodo,
		TeamID:    teamID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.taskReader.On("GetByID", mock.Anything, (*sqlx.Tx)(nil), taskID).
		Return(task, nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, taskID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	s.taskReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_WithTx() {
	tx := &sqlx.Tx{}
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	task := &model.Task{ID: taskID, Title: "Task", TeamID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")}

	s.taskReader.On("GetByID", mock.Anything, tx, taskID).
		Return(task, nil).Once()

	got, err := s.repo.GetByID(s.ctx, tx, taskID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	s.taskReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_TaskNotFound() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	s.taskReader.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.repo.GetByID(s.ctx, nil, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskReader.AssertExpectations(s.T())
}
