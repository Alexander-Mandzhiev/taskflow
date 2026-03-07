package adapter_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestAddMember_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("770e8400-e29b-41d4-a716-446655440002")
	role := model2.RoleMember
	member := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    userID,
		TeamID:    teamID,
		Role:      role,
		CreatedAt: time.Now(),
	}

	s.memberWriter.On("AddMember", mock.Anything, (*sqlx.Tx)(nil), teamID.String(), userID.String(), role).
		Return(member, nil).Once()

	got, err := s.repo.AddMember(s.ctx, nil, teamID, userID, role)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	assert.Equal(s.T(), role, got.Role)
	s.memberWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestAddMember_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("770e8400-e29b-41d4-a716-446655440002")
	role := model2.RoleAdmin
	member := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    userID,
		TeamID:    teamID,
		Role:      role,
		CreatedAt: time.Now(),
	}

	s.memberWriter.On("AddMember", mock.Anything, tx, teamID.String(), userID.String(), role).
		Return(member, nil).Once()

	got, err := s.repo.AddMember(s.ctx, tx, teamID, userID, role)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestAddMember_WriterError() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	userID := uuid.MustParse("770e8400-e29b-41d4-a716-446655440002")
	role := model2.RoleMember

	s.memberWriter.On("AddMember", mock.Anything, mock.Anything, teamID.String(), userID.String(), role).
		Return((*model2.TeamMember)(nil), model2.ErrAlreadyMember).Once()

	got, err := s.repo.AddMember(s.ctx, nil, teamID, userID, role)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrAlreadyMember)
	assert.Nil(s.T(), got)
	s.memberWriter.AssertExpectations(s.T())
}
