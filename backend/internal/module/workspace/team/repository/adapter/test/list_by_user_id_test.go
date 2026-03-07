package adapter_test

import (
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *AdapterSuite) TestListByUserID_Success() {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	teams := []*model2.TeamWithRole{
		{
			Team: model2.Team{
				ID:        uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
				Name:      "My Team",
				CreatedBy: userID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Role: model2.RoleOwner,
		},
	}

	s.teamReader.On("ListByUserID", mock.Anything, (*sqlx.Tx)(nil), userID.String()).
		Return(teams, nil).Once()

	got, err := s.repo.ListByUserID(s.ctx, nil, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), got, 1)
	assert.Equal(s.T(), "My Team", got[0].Name)
	assert.Equal(s.T(), model2.RoleOwner, got[0].Role)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListByUserID_Empty() {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	s.teamReader.On("ListByUserID", mock.Anything, mock.Anything, userID.String()).
		Return([]*model2.TeamWithRole{}, nil).Once()

	got, err := s.repo.ListByUserID(s.ctx, nil, userID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListByUserID_ReaderError() {
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	s.teamReader.On("ListByUserID", mock.Anything, mock.Anything, userID.String()).
		Return(([]*model2.TeamWithRole)(nil), assert.AnError).Once()

	got, err := s.repo.ListByUserID(s.ctx, nil, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
}
