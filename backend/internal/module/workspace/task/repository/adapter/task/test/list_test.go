package task_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestList_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	status := model.TaskStatusTodo
	filter := &model.TaskListFilter{TeamID: &teamID, Status: &status, Limit: 10, Offset: 0}
	tasks := []*model.Task{
		{ID: uuid.New(), Title: "Task 1", TeamID: teamID, Status: model.TaskStatusTodo},
		{ID: uuid.New(), Title: "Task 2", TeamID: teamID, Status: model.TaskStatusTodo},
	}

	s.taskReader.On("List", mock.Anything, (*sqlx.Tx)(nil), filter).
		Return(tasks, 2, nil).Once()

	got, total, err := s.repo.List(s.ctx, nil, filter)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 2)
	assert.Equal(s.T(), 2, total)
	s.taskReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestList_Empty() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 10, Offset: 0}

	s.taskReader.On("List", mock.Anything, mock.Anything, filter).
		Return([]*model.Task{}, 0, nil).Once()

	got, total, err := s.repo.List(s.ctx, nil, filter)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.taskReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestList_ReaderError() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	filter := &model.TaskListFilter{TeamID: &teamID, Limit: 10, Offset: 0}

	s.taskReader.On("List", mock.Anything, mock.Anything, filter).
		Return(([]*model.Task)(nil), 0, assert.AnError).Once()

	got, total, err := s.repo.List(s.ctx, nil, filter)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	assert.Equal(s.T(), 0, total)
	s.taskReader.AssertExpectations(s.T())
}
