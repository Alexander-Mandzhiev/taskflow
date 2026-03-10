package report_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
)

func (s *AdapterSuite) TestTeamTaskStats_Success() {
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	stats := []model.TeamTaskStats{
		{
			TeamID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			TeamName:       "Team A",
			MemberCount:    5,
			DoneTasksCount: 12,
		},
	}

	s.reportReader.On("TeamTaskStats", mock.Anything, (*sqlx.Tx)(nil), since).
		Return(stats, nil).Once()

	got, err := s.repo.TeamTaskStats(s.ctx, nil, since)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "Team A", got[0].TeamName)
	assert.Equal(s.T(), 5, got[0].MemberCount)
	assert.Equal(s.T(), 12, got[0].DoneTasksCount)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTeamTaskStats_Empty() {
	since := time.Now()

	s.reportReader.On("TeamTaskStats", mock.Anything, mock.Anything, since).
		Return([]model.TeamTaskStats{}, nil).Once()

	got, err := s.repo.TeamTaskStats(s.ctx, nil, since)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestTeamTaskStats_WithTx_ReaderError() {
	tx := &sqlx.Tx{}
	since := time.Now()

	s.reportReader.On("TeamTaskStats", mock.Anything, tx, since).
		Return(nil, assert.AnError).Once()

	got, err := s.repo.TeamTaskStats(s.ctx, tx, since)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.reportReader.AssertExpectations(s.T())
}
