package report_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestTasksWithInvalidAssignee_Success() {
	userID := uuid.New()
	teamID := uuid.New()
	userTeams := []*teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}
	allTasks := []*model.Task{
		{ID: uuid.New(), Title: "Bad assignee task", TeamID: teamID},
		{ID: uuid.New(), Title: "Other team task", TeamID: uuid.New()},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TasksWithInvalidAssignee", mock.Anything, mock.Anything).Return(allTasks, nil).Once()

	got, err := s.svc.TasksWithInvalidAssignee(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "Bad assignee task", got[0].Title)
	assert.Equal(s.T(), teamID, got[0].TeamID)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestTasksWithInvalidAssignee_Empty() {
	userID := uuid.New()
	userTeams := []*teamModel.TeamWithRole{}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TasksWithInvalidAssignee", mock.Anything, mock.Anything).Return([]*model.Task{}, nil).Once()

	got, err := s.svc.TasksWithInvalidAssignee(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestTasksWithInvalidAssignee_ListByUserIDError() {
	userID := uuid.New()

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(([]*teamModel.TeamWithRole)(nil), assert.AnError).Once()

	got, err := s.svc.TasksWithInvalidAssignee(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertNotCalled(s.T(), "TasksWithInvalidAssignee")
}

func (s *ServiceSuite) TestTasksWithInvalidAssignee_RepoError() {
	userID := uuid.New()
	teamID := uuid.New()
	userTeams := []*teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TasksWithInvalidAssignee", mock.Anything, mock.Anything).
		Return(([]*model.Task)(nil), assert.AnError).Once()

	got, err := s.svc.TasksWithInvalidAssignee(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}
