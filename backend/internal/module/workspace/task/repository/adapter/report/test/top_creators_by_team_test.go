package report_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestTopCreatorsByTeam_Success() {
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	limit := 5
	creators := []model.TeamTopCreator{
		{
			TeamID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			UserID:       uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
			Rank:         1,
			CreatedCount: 42,
		},
	}

	s.reportReader.On("TopCreatorsByTeam", mock.Anything, (*sqlx.Tx)(nil), since, limit).
		Return(creators, nil).Once()

	got, err := s.repo.TopCreatorsByTeam(s.ctx, nil, since, limit)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), 1, got[0].Rank)
	assert.Equal(s.T(), int64(42), got[0].CreatedCount)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTopCreatorsByTeam_Empty() {
	since := time.Now()
	limit := 3

	s.reportReader.On("TopCreatorsByTeam", mock.Anything, mock.Anything, since, limit).
		Return([]model.TeamTopCreator{}, nil).Once()

	got, err := s.repo.TopCreatorsByTeam(s.ctx, nil, since, limit)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTopCreatorsByTeam_ReaderError() {
	since := time.Now()
	limit := 10

	s.reportReader.On("TopCreatorsByTeam", mock.Anything, mock.Anything, since, limit).
		Return(nil, assert.AnError).Once()

	got, err := s.repo.TopCreatorsByTeam(s.ctx, nil, since, limit)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}
