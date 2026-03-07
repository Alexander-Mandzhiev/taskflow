package task_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestUpdate_Success() {
	tx := &sqlx.Tx{}
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	input := &model.TaskInput{Title: "Updated", Description: "New desc", Status: model.TaskStatusInProgress}

	s.taskWriter.On("Update", mock.Anything, tx, taskID, input).
		Return(nil).Once()

	err := s.repo.Update(s.ctx, tx, taskID, input)

	assert.NoError(s.T(), err)
	s.taskWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestUpdate_WriterError() {
	taskID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440002")
	input := &model.TaskInput{Title: "Updated", Status: model.TaskStatusDone}

	s.taskWriter.On("Update", mock.Anything, mock.Anything, taskID, input).
		Return(model.ErrTaskNotFound).Once()

	err := s.repo.Update(s.ctx, nil, taskID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	s.taskWriter.AssertExpectations(s.T())
}
