package report_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestTeamTaskStats_Success() {
	userID := uuid.New()
	teamID := uuid.New()
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	userTeams := []*teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}
	allStats := []*model.TeamTaskStats{
		{TeamID: teamID, TeamName: "My Team", MemberCount: 5, DoneTasksCount: 12},
		{TeamID: uuid.New(), TeamName: "Other", MemberCount: 3, DoneTasksCount: 0},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TeamTaskStats", mock.Anything, mock.Anything, since).Return(allStats, nil).Once()

	got, err := s.svc.TeamTaskStats(s.ctx, userID, since)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), teamID, got[0].TeamID)
	assert.Equal(s.T(), "My Team", got[0].TeamName)
	assert.Equal(s.T(), 5, got[0].MemberCount)
	assert.Equal(s.T(), 12, got[0].DoneTasksCount)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestTeamTaskStats_EmptyTeams() {
	userID := uuid.New()
	since := time.Now()

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return([]*teamModel.TeamWithRole{}, nil).Once()
	s.reportRepo.On("TeamTaskStats", mock.Anything, mock.Anything, since).Return([]*model.TeamTaskStats{}, nil).Once()

	got, err := s.svc.TeamTaskStats(s.ctx, userID, since)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestTeamTaskStats_ListByUserIDError() {
	userID := uuid.New()
	since := time.Now()

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(([]*teamModel.TeamWithRole)(nil), assert.AnError).Once()

	got, err := s.svc.TeamTaskStats(s.ctx, userID, since)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertNotCalled(s.T(), "TeamTaskStats")
}

func (s *ServiceSuite) TestTeamTaskStats_RepoError() {
	userID := uuid.New()
	teamID := uuid.New()
	since := time.Now()
	userTeams := []*teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TeamTaskStats", mock.Anything, mock.Anything, since).
		Return(([]*model.TeamTaskStats)(nil), assert.AnError).Once()

	got, err := s.svc.TeamTaskStats(s.ctx, userID, since)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}
