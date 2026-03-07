package adapter_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

const (
	testTeamID      = "550e8400-e29b-41d4-a716-446655440001"
	testOwnerUserID = "550e8400-e29b-41d4-a716-446655440000"
)

func (s *AdapterSuite) TestGetByID_Success() {
	teamID := uuid.MustParse(testTeamID)
	team := &model2.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.MustParse(testOwnerUserID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamReader.On("GetByID", mock.Anything, (*sqlx.Tx)(nil), testTeamID).
		Return(team, nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, teamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), testTeamID, got.ID.String())
	assert.Equal(s.T(), "My Team", got.Name)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse(testTeamID)
	team := &model2.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.MustParse(testOwnerUserID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamReader.On("GetByID", mock.Anything, tx, testTeamID).
		Return(team, nil).Once()

	got, err := s.repo.GetByID(s.ctx, tx, teamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), testTeamID, got.ID.String())
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_TeamNotFound() {
	teamID := uuid.MustParse(testTeamID)
	s.teamReader.On("GetByID", mock.Anything, mock.Anything, testTeamID).
		Return((*model2.Team)(nil), model2.ErrTeamNotFound).Once()

	got, err := s.repo.GetByID(s.ctx, nil, teamID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrTeamNotFound)
	assert.Nil(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
}
