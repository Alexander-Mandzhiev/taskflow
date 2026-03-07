package team_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *ServiceSuite) TestGetByID_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	team := &model.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	members := []*model.TeamMember{
		{ID: uuid.New(), UserID: userID, TeamID: teamID, Role: model.RoleOwner, CreatedAt: time.Now()},
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

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return((*model.Team)(nil), model.ErrTeamNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTeamNotFound)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_Forbidden_NotMember() {
	teamID := uuid.New()
	userID := uuid.New()
	team := &model.Team{ID: teamID, Name: "Team", CreatedBy: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	members := []*model.TeamMember{}

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(team, nil).Once()
	s.teamRepo.On("GetMembersByTeamID", mock.Anything, mock.Anything, teamID).Return(members, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, userID).Return((*model.TeamMember)(nil), model.ErrMemberNotFound).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrForbidden)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetByID_MembersError() {
	teamID := uuid.New()
	userID := uuid.New()
	team := &model.Team{ID: teamID, Name: "Team", CreatedBy: userID, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(team, nil).Once()
	s.teamRepo.On("GetMembersByTeamID", mock.Anything, mock.Anything, teamID).Return(([]*model.TeamMember)(nil), assert.AnError).Once()

	got, err := s.svc.GetByID(s.ctx, teamID, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamRepo.AssertExpectations(s.T())
}
