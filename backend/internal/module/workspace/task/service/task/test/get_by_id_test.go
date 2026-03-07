package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestGetByID_Success() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{
		ID: taskID, Title: "Task", TeamID: teamID, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()

	got, err := s.svc.GetByID(s.ctx, taskID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), task, got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_TaskNotFound() {
	taskID := uuid.New()
	userID := uuid.New()

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertNotCalled(s.T(), "GetMember")
}

func (s *ServiceSuite) TestGetByID_NotMember() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
}
