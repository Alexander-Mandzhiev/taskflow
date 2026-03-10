package report_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	teamModel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestTopCreatorsByTeam_InvalidLimit() {
	userID := uuid.New()
	since := time.Now()

	got, err := s.svc.TopCreatorsByTeam(s.ctx, userID, since, 0)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvalidLimit)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertNotCalled(s.T(), "ListByUserID")
	s.reportRepo.AssertNotCalled(s.T(), "TopCreatorsByTeam")
}

func (s *ServiceSuite) TestTopCreatorsByTeam_NegativeLimit() {
	userID := uuid.New()
	since := time.Now()

	got, err := s.svc.TopCreatorsByTeam(s.ctx, userID, since, -1)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvalidLimit)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertNotCalled(s.T(), "ListByUserID")
	s.reportRepo.AssertNotCalled(s.T(), "TopCreatorsByTeam")
}

func (s *ServiceSuite) TestTopCreatorsByTeam_Success() {
	userID := uuid.New()
	teamID := uuid.New()
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	limit := 5
	userTeams := []teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}
	allCreators := []model.TeamTopCreator{
		{TeamID: teamID, UserID: userID, Rank: 1, CreatedCount: 42},
		{TeamID: uuid.New(), UserID: uuid.New(), Rank: 1, CreatedCount: 10},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TopCreatorsByTeam", mock.Anything, mock.Anything, since, limit).Return(allCreators, nil).Once()

	got, err := s.svc.TopCreatorsByTeam(s.ctx, userID, since, limit)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), teamID, got[0].TeamID)
	assert.Equal(s.T(), int64(42), got[0].CreatedCount)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestTopCreatorsByTeam_ListByUserIDError() {
	userID := uuid.New()
	since := time.Now()
	limit := 3

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(nil, assert.AnError).Once()

	got, err := s.svc.TopCreatorsByTeam(s.ctx, userID, since, limit)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertNotCalled(s.T(), "TopCreatorsByTeam")
}

func (s *ServiceSuite) TestTopCreatorsByTeam_RepoError() {
	userID := uuid.New()
	teamID := uuid.New()
	since := time.Now()
	limit := 5
	userTeams := []teamModel.TeamWithRole{
		{Team: teamModel.Team{ID: teamID, Name: "My Team"}, Role: teamModel.RoleOwner},
	}

	s.teamSvc.On("ListByUserID", mock.Anything, userID).Return(userTeams, nil).Once()
	s.reportRepo.On("TopCreatorsByTeam", mock.Anything, mock.Anything, since, limit).
		Return(nil, assert.AnError).Once()

	got, err := s.svc.TopCreatorsByTeam(s.ctx, userID, since, limit)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamSvc.AssertExpectations(s.T())
	s.reportRepo.AssertExpectations(s.T())
}
