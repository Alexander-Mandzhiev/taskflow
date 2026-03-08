package team_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestCreate_Success() {
	ownerID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	input := &model.TeamInput{Name: "New Team"}
	team := &model.Team{
		ID:        uuid.New(),
		Name:      input.Name,
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.teamWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), input, ownerID).
		Return(team, nil).Once()

	got, err := s.repo.Create(s.ctx, nil, input, ownerID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WithTx() {
	tx := &sqlx.Tx{}
	ownerID := uuid.New()
	input := &model.TeamInput{Name: "Team"}
	team := &model.Team{ID: uuid.New(), Name: input.Name, CreatedBy: ownerID}

	s.teamWriter.On("Create", mock.Anything, tx, input, ownerID).
		Return(team, nil).Once()

	got, err := s.repo.Create(s.ctx, tx, input, ownerID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WriterError() {
	ownerID := uuid.New()
	input := &model.TeamInput{Name: "Team"}

	s.teamWriter.On("Create", mock.Anything, mock.Anything, input, ownerID).
		Return((*model.Team)(nil), assert.AnError).Once()

	got, err := s.repo.Create(s.ctx, nil, input, ownerID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamWriter.AssertExpectations(s.T())
}
