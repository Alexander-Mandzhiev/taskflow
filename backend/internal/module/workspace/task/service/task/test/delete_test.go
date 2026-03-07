package task_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestDelete_Success() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("SoftDelete", mock.Anything, mock.Anything, taskID).Return(nil).Once()

	err := s.svc.Delete(s.ctx, userID, taskID)

	assert.NoError(s.T(), err)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestDelete_TaskNotFound() {
	userID := uuid.New()
	taskID := uuid.New()

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	err := s.svc.Delete(s.ctx, userID, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertNotCalled(s.T(), "GetMember")
}

func (s *ServiceSuite) TestDelete_NotMember() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	err := s.svc.Delete(s.ctx, userID, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "SoftDelete")
}

func (s *ServiceSuite) TestDelete_SoftDeleteError() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("SoftDelete", mock.Anything, mock.Anything, taskID).Return(assert.AnError).Once()

	err := s.svc.Delete(s.ctx, userID, taskID)

	assert.Error(s.T(), err)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
}
