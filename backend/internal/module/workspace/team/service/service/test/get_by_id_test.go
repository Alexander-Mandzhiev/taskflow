package service_test

import (
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestGetByID_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	team := &model2.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	members := []*model2.TeamMember{
		{ID: uuid.New(), UserID: userID, TeamID: teamID, Role: model2.RoleOwner, CreatedAt: time.Now()},
	}

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(team, nil).Once()
	s.teamRepo.On("GetMembersByTeamID", mock.Anything, mock.Anything, teamID).Return(members, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return(members[0], nil).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), teamID, got.ID)
	assert.Equal(s.T(), "My Team", got.Name)
	assert.Len(s.T(), got.Members, 1)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_TeamNotFound() {
	teamID := uuid.New()
	userID := uuid.New()

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return((*model2.Team)(nil), model2.ErrTeamNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrTeamNotFound)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_Forbidden_NotMember() {
	teamID := uuid.New()
	userID := uuid.New()
	team := &model2.Team{ID: teamID, Name: "Team", CreatedBy: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	members := []*model2.TeamMember{}

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(team, nil).Once()
	s.teamRepo.On("GetMembersByTeamID", mock.Anything, mock.Anything, teamID).Return(members, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return((*model2.TeamMember)(nil), model2.ErrMemberNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrForbidden)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_MembersError() {
	teamID := uuid.New()
	userID := uuid.New()
	team := &model2.Team{ID: teamID, Name: "Team", CreatedBy: userID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(team, nil).Once()
	s.teamRepo.On("GetMembersByTeamID", mock.Anything, mock.Anything, teamID).Return(([]*model2.TeamMember)(nil), assert.AnError).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}
