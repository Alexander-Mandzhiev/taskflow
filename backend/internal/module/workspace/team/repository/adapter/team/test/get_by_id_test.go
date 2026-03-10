package team_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestGetByID_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	team := model.Team{
		ID:        teamID,
		Name:      "My Team",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamReader.On("GetByID", mock.Anything, (*sqlx.Tx)(nil), teamID).
		Return(team, nil).Once()

	got, err := s.repo.GetByID(s.ctx, nil, teamID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.New()
	team := model.Team{ID: teamID, Name: "Team"}

	s.teamReader.On("GetByID", mock.Anything, tx, teamID).
		Return(team, nil).Once()

	got, err := s.repo.GetByID(s.ctx, tx, teamID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetByID_TeamNotFound() {
	teamID := uuid.New()
	s.teamReader.On("GetByID", mock.Anything, mock.Anything, teamID).
		Return(model.Team{}, model.ErrTeamNotFound).Once()

	got, err := s.repo.GetByID(s.ctx, nil, teamID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrTeamNotFound)
	assert.Equal(s.T(), model.Team{}, got)
	s.teamReader.AssertExpectations(s.T())
}
