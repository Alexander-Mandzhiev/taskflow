package team_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestListByUserID_Success() {
	userID := uuid.New()
	want := []model.TeamWithRole{
		{Team: model.Team{ID: uuid.New(), Name: "Team A"}, Role: model.RoleOwner},
		{Team: model.Team{ID: uuid.New(), Name: "Team B"}, Role: model.RoleMember},
	}

	s.teamReader.On("ListByUserID", mock.Anything, (*sqlx.Tx)(nil), userID).
		Return(want, nil).Once()

	got, err := s.repo.ListByUserID(s.ctx, nil, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListByUserID_Empty() {
	userID := uuid.New()
	s.teamReader.On("ListByUserID", mock.Anything, mock.Anything, userID).
		Return([]model.TeamWithRole{}, nil).Once()

	got, err := s.repo.ListByUserID(s.ctx, nil, userID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Empty(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestListByUserID_ReaderError() {
	userID := uuid.New()
	s.teamReader.On("ListByUserID", mock.Anything, mock.Anything, userID).
		Return(nil, assert.AnError).Once()

	got, err := s.repo.ListByUserID(s.ctx, &sqlx.Tx{}, userID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.teamReader.AssertExpectations(s.T())
}
