package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestRestore_Success() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}
	restored := &model.Task{ID: taskID, Title: "Task", TeamID: teamID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	s.taskRepo.On("GetByIDIncludeDeleted", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).Return(member, nil).Once()
	s.taskRepo.On("Restore", mock.Anything, mock.Anything, taskID).Return(nil).Once()
	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(restored, nil).Once()

	got, err := s.svc.Restore(s.ctx, userID, taskID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), restored, got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRestore_TaskNotFound() {
	userID := uuid.New()
	taskID := uuid.New()

	s.taskRepo.On("GetByIDIncludeDeleted", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.svc.Restore(s.ctx, userID, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertNotCalled(s.T(), "GetMember")
}

func (s *ServiceSuite) TestRestore_NotMember() {
	userID := uuid.New()
	taskID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}

	s.taskRepo.On("GetByIDIncludeDeleted", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamSvc.On("GetMember", mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.Restore(s.ctx, userID, taskID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamSvc.AssertExpectations(s.T())
	s.taskRepo.AssertNotCalled(s.T(), "Restore")
}
