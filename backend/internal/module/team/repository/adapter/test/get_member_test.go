package adapter_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

func (s *AdapterSuite) TestGetMember_Success() {
	teamID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440000"
	member := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		TeamID:    uuid.MustParse(teamID),
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.memberReader.On("GetMember", mock.Anything, (*sqlx.Tx)(nil), teamID, userID).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	assert.Equal(s.T(), model.RoleOwner, got.Role)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_WithTx() {
	tx := &sqlx.Tx{}
	teamID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440000"
	member := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(userID),
		TeamID:    uuid.MustParse(teamID),
		Role:      model.RoleAdmin,
		CreatedAt: time.Now(),
	}

	s.memberReader.On("GetMember", mock.Anything, tx, teamID, userID).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, tx, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_NotFound() {
	teamID := "550e8400-e29b-41d4-a716-446655440001"
	userID := "550e8400-e29b-41d4-a716-446655440000"

	s.memberReader.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return((*model.TeamMember)(nil), model.ErrMemberNotFound).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrMemberNotFound)
	assert.Nil(s.T(), got)
	s.memberReader.AssertExpectations(s.T())
}
