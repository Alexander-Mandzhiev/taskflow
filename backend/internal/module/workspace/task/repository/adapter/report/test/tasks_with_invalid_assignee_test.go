package report_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestTasksWithInvalidAssignee_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	tasks := []model.Task{
		{
			ID:     uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
			Title:  "Task with bad assignee",
			TeamID: teamID,
		},
	}

	s.reportReader.On("TasksWithInvalidAssignee", mock.Anything, (*sqlx.Tx)(nil)).
		Return(tasks, nil).Once()

	got, err := s.repo.TasksWithInvalidAssignee(s.ctx, nil)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "Task with bad assignee", got[0].Title)
	assert.Equal(s.T(), teamID, got[0].TeamID)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTasksWithInvalidAssignee_Empty() {
	s.reportReader.On("TasksWithInvalidAssignee", mock.Anything, mock.Anything).
		Return([]model.Task{}, nil).Once()

	got, err := s.repo.TasksWithInvalidAssignee(s.ctx, nil)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTasksWithInvalidAssignee_WithTx_ReaderError() {
	tx := &sqlx.Tx{}
	s.reportReader.On("TasksWithInvalidAssignee", mock.Anything, tx).
		Return(nil, assert.AnError).Once()

	got, err := s.repo.TasksWithInvalidAssignee(s.ctx, tx)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}
