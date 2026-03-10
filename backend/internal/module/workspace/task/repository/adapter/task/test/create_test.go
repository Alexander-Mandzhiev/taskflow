package task_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestCreate_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	createdBy := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	input := model.TaskInput{Title: "Task", Description: "Desc", Status: model.TaskStatusTodo}
	task := model.Task{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		TeamID:      teamID,
		CreatedBy:   createdBy,
	}

	s.taskWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), teamID, input, createdBy).
		Return(task, nil).Once()

	got, err := s.repo.Create(s.ctx, nil, teamID, input, createdBy)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	s.taskWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	createdBy := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	input := model.TaskInput{Title: "Task", Status: model.TaskStatusTodo}
	task := model.Task{
		ID: uuid.New(), Title: input.Title, Status: input.Status, TeamID: teamID, CreatedBy: createdBy,
	}

	s.taskWriter.On("Create", mock.Anything, tx, teamID, input, createdBy).
		Return(task, nil).Once()

	got, err := s.repo.Create(s.ctx, tx, teamID, input, createdBy)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	s.taskWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WriterError() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	createdBy := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	input := model.TaskInput{Title: "Task", Status: model.TaskStatusTodo}

	s.taskWriter.On("Create", mock.Anything, mock.Anything, teamID, input, createdBy).
		Return(model.Task{}, assert.AnError).Once()

	got, err := s.repo.Create(s.ctx, nil, teamID, input, createdBy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Task{}, got)
	s.taskWriter.AssertExpectations(s.T())
}
