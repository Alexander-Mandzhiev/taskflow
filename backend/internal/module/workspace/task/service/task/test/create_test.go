package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// Валидация входящих данных (title, status и т.д.) выполняется в API слое на DTO (CreateTaskRequest, UpdateTaskRequest).

func (s *ServiceSuite) TestCreate_UserNotMember() {
	userID := uuid.New()
	teamID := uuid.New()
	input := model.TaskInput{Title: "Task", Status: model.TaskStatusTodo}

	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return(teamModel.TeamMember{}, teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Create(s.ctx, userID, teamID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Equal(s.T(), model.Task{}, got)
	s.memberRepo.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreate_AssigneeNotInTeam() {
	userID := uuid.New()
	teamID := uuid.New()
	assigneeID := uuid.New()
	input := model.TaskInput{Title: "Task", Status: model.TaskStatusTodo, AssigneeID: &assigneeID}
	member := teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, assigneeID).
		Return(teamModel.TeamMember{}, teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Create(s.ctx, userID, teamID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrAssigneeNotInTeam)
	assert.Equal(s.T(), model.Task{}, got)
	s.memberRepo.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreate_Success() {
	userID := uuid.New()
	teamID := uuid.New()
	input := model.TaskInput{Title: "New Task", Description: "Desc", Status: model.TaskStatusTodo}
	member := teamModel.TeamMember{UserID: userID, TeamID: teamID}
	created := model.Task{
		ID: uuid.New(), Title: input.Title, Description: input.Description, Status: model.TaskStatusTodo,
		TeamID: teamID, CreatedBy: userID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("Create", mock.Anything, mock.Anything, teamID, mock.Anything, userID).
		Return(created, nil).Once()

	got, err := s.svc.Create(s.ctx, userID, teamID, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), created, got)
	s.memberRepo.AssertExpectations(s.T())
	s.taskRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreate_RepoError() {
	userID := uuid.New()
	teamID := uuid.New()
	input := model.TaskInput{Title: "Task", Status: model.TaskStatusTodo}
	member := teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.memberRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("Create", mock.Anything, mock.Anything, teamID, mock.Anything, userID).
		Return(model.Task{}, assert.AnError).Once()

	got, err := s.svc.Create(s.ctx, userID, teamID, input)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Task{}, got)
	s.memberRepo.AssertExpectations(s.T())
	s.taskRepo.AssertExpectations(s.T())
}
