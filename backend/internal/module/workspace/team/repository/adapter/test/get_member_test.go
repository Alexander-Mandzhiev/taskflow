package adapter_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestGetMember_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	member := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    userID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.memberReader.On("GetMember", mock.Anything, (*sqlx.Tx)(nil), teamID.String(), userID.String()).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	assert.Equal(s.T(), model2.RoleOwner, got.Role)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	member := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    userID,
		TeamID:    teamID,
		Role:      model2.RoleAdmin,
		CreatedAt: time.Now(),
	}

	s.memberReader.On("GetMember", mock.Anything, tx, teamID.String(), userID.String()).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, tx, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_NotFound() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	s.memberReader.On("GetMember", mock.Anything, mock.Anything, teamID.String(), userID.String()).
		Return((*model2.TeamMember)(nil), model2.ErrMemberNotFound).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrMemberNotFound)
	assert.Nil(s.T(), got)
	s.memberReader.AssertExpectations(s.T())
}
