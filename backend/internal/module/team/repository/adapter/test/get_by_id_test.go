package adapter_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

const testTeamID = "550e8400-e29b-41d4-a716-446655440001"
const testOwnerUserID = "550e8400-e29b-41d4-a716-446655440000"

func (s *AdapterSuite) TestGetByID_Success() {
	teamID := uuid.MustParse(testTeamID)
	team := &model.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.MustParse(testOwnerUserID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	members := []*model.TeamMember{
		{
			ID:        uuid.New(),
			UserID:    uuid.MustParse(testOwnerUserID),
			TeamID:    teamID,
			Role:      model.RoleOwner,
			CreatedAt: time.Now(),
		},
	}

	s.teamReader.On("GetByID", mock.Anything, (*sqlx.Tx)(nil), testTeamID).
		Return(team, nil).Once()
	s.memberReader.On("GetByTeamID", mock.Anything, (*sqlx.Tx)(nil), testTeamID).
		Return(members, nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, testTeamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), testTeamID, got.Team.ID.String())
	assert.Equal(s.T(), "My Team", got.Team.Name)
	assert.Len(s.T(), got.Members, 1)
	assert.Equal(s.T(), model.RoleOwner, got.Members[0].Role)
	s.teamReader.AssertExpectations(s.T())
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse(testTeamID)
	team := &model.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.MustParse(testOwnerUserID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamReader.On("GetByID", mock.Anything, tx, testTeamID).
		Return(team, nil).Once()
	s.memberReader.On("GetByTeamID", mock.Anything, tx, testTeamID).
		Return([]*model.TeamMember{}, nil).Once()

	got, err := s.repo.GetByID(s.ctx, tx, testTeamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), testTeamID, got.Team.ID.String())
	assert.Empty(s.T(), got.Members)
	s.teamReader.AssertExpectations(s.T())
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_TeamNotFound() {
	s.teamReader.On("GetByID", mock.Anything, mock.Anything, testTeamID).
		Return((*model.Team)(nil), model.ErrTeamNotFound).Once()

	got, err := s.repo.GetByID(s.ctx, nil, testTeamID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTeamNotFound)
	assert.Nil(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
	s.memberReader.AssertNotCalled(s.T(), "GetByTeamID")
}

func (s *AdapterSuite) TestGetByID_MembersError() {
	teamID := uuid.MustParse(testTeamID)
	team := &model.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.MustParse(testOwnerUserID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamReader.On("GetByID", mock.Anything, mock.Anything, testTeamID).
		Return(team, nil).Once()
	s.memberReader.On("GetByTeamID", mock.Anything, mock.Anything, testTeamID).
		Return(([]*model.TeamMember)(nil), assert.AnError).Once()

	got, err := s.repo.GetByID(s.ctx, nil, testTeamID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
	s.memberReader.AssertExpectations(s.T())
}
