package task_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestGetHistory_Success() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}
	entries := []*model.TaskHistory{
		{
			ID: uuid.New(), TaskID: taskID, ChangedBy: userID,
			FieldName: "title", OldValue: "Old", NewValue: "New", ChangedAt: time.Now(),
		},
	}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.historyRepo.On("ListHistoryByTaskID", mock.Anything, mock.Anything, taskID).Return(entries, nil).Once()

	got, err := s.svc.GetHistory(s.ctx, taskID, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "title", got[0].FieldName)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.historyRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetHistory_TaskNotFound() {
	taskID := uuid.New()
	userID := uuid.New()

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).
		Return((*model.Task)(nil), model.ErrTaskNotFound).Once()

	got, err := s.svc.GetHistory(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertNotCalled(s.T(), "GetMember")
	s.historyRepo.AssertNotCalled(s.T(), "ListHistoryByTaskID")
}

func (s *ServiceSuite) TestGetHistory_NotMember() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return((*teamModel.TeamMember)(nil), teamModel.ErrMemberNotFound).Once()

	got, err := s.svc.GetHistory(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTaskNotFound)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.historyRepo.AssertNotCalled(s.T(), "ListHistoryByTaskID")
}

func (s *ServiceSuite) TestGetHistory_HistoryRepoError() {
	taskID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	task := &model.Task{ID: taskID, Title: "Task", TeamID: teamID}
	member := &teamModel.TeamMember{UserID: userID, TeamID: teamID}

	s.taskRepo.On("GetByID", mock.Anything, mock.Anything, taskID).Return(task, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(member, nil).Once()
	s.historyRepo.On("ListHistoryByTaskID", mock.Anything, mock.Anything, taskID).
		Return(([]*model.TaskHistory)(nil), assert.AnError).Once()

	got, err := s.svc.GetHistory(s.ctx, taskID, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.taskRepo.AssertExpectations(s.T())
	s.teamRepo.AssertExpectations(s.T())
	s.historyRepo.AssertExpectations(s.T())
}
