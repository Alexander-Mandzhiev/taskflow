package adapter_test

import (
	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *AdapterSuite) TestCreate_Success() {
	input := &model2.TeamInput{Name: "My Team"}
	ownerUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	team := &model2.Team{
		ID:        uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		Name:      input.Name,
		CreatedBy: ownerUserID,
	}

	s.teamWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), input, ownerUserID.String()).
		Return(team, nil).Once()

	got, err := s.repo.Create(s.ctx, nil, input, ownerUserID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WithTx() {
	tx := &sqlx.Tx{}
	input := &model2.TeamInput{Name: "My Team"}
	ownerUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	team := &model2.Team{
		ID:        uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
		Name:      input.Name,
		CreatedBy: ownerUserID,
	}

	s.teamWriter.On("Create", mock.Anything, tx, input, ownerUserID.String()).
		Return(team, nil).Once()

	got, err := s.repo.Create(s.ctx, tx, input, ownerUserID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), team, got)
	s.teamWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreate_WriterError() {
	input := &model2.TeamInput{Name: "My Team"}
	ownerUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	s.teamWriter.On("Create", mock.Anything, mock.Anything, input, ownerUserID.String()).
		Return((*model2.Team)(nil), assert.AnError).Once()

	got, err := s.repo.Create(s.ctx, nil, input, ownerUserID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamWriter.AssertExpectations(s.T())
}
