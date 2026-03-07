package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestUpdate_NilInput() {
	userID := uuid.New()
	taskID := uuid.New()

	got, err := s.svc.Update(s.ctx, userID, taskID, nil)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrNilInput)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertNotCalled(s.T(), "GetByID")
	s.taskRepo.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdate_InvalidStatus() {
	userID := uuid.New()
	taskID := uuid.New()
	input := &model.TaskInput{Title: "Updated", Status: "invalid"}

	got, err := s.svc.Update(s.ctx, userID, taskID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvalidStatus)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertNotCalled(s.T(), "GetByID")
	s.taskRepo.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdate_TaskNotFound() {
	userID := uuid.New()
	taskID := uuid.New()
	input := &model.TaskInput{Title: "Updated", Status: model.TaskStatusTodo}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.svc.Update(s.ctx, userID, taskID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertNotCalled(s.T(), "GetMember")
	s.taskRepo.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdate_NotMember() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	current := &model.Task{ID: taskID, Title: "Old", TeamID: teamID}
	input := &model.TaskInput{Title: "New", Status: model.TaskStatusTodo}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(current, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Update(s.ctx, userID, taskID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "Update")
	s.historyRepo.AssertNotCalled(s.T(), "CreateHistoryEntry")
}

func (s *ServiceSuite) TestUpdate_AssigneeNotInTeam() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	assigneeID := uuid.New()
	current := &model.Task{ID: taskID, Title: "Old", TeamID: teamID}
	input := &model.TaskInput{Title: "New", Status: model.TaskStatusTodo, AssigneeID: &assigneeID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(current, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).Return(member, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, assigneeID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Update(s.ctx, userID, taskID, input)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrAssigneeNotInTeam)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestUpdate_Success() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	current := &model.Task{
		ID: taskID, Title: "Old", Description: "Desc", Status: model.TaskStatusTodo,
		TeamID: teamID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	input := &model.TaskInput{Title: "New", Description: "Desc", Status: model.TaskStatusTodo}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}
	updated := &model.Task{
		ID: taskID, Title: "New", Description: "Desc", Status: model.TaskStatusTodo,
		TeamID: teamID, CreatedAt: current.CreatedAt, UpdatedAt: time.Now(),
	}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(current, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("Update", mock.Anything, mock.Anything, taskID, input).Return(nil).Once()
	s.historyRepo.On("CreateHistoryEntry", mock.Anything, mock.Anything, mock.MatchedBy(func(e *model.TaskHistory) bool {
		return e.FieldName == "title" && e.OldValue == "Old" && e.NewValue == "New"
	})).Return(nil).Once()
	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(updated, nil).Once()

	got, err := s.svc.Update(s.ctx, userID, taskID, input)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updated, got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertExpectations(s.T())
	s.historyRepo.AssertExpectations(s.T())
}
