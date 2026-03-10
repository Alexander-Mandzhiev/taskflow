package member_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestGetMember_Success() {
	teamID := uuid.New()
	userID := uuid.New()
	member := model.TeamMember{TeamID: teamID, UserID: userID, Role: model.RoleAdmin}

	s.memberReader.On("GetMember", mock.Anything, (*sqlx.Tx)(nil), teamID, userID).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.New()
	userID := uuid.New()
	member := model.TeamMember{TeamID: teamID, UserID: userID, Role: model.RoleMember}

	s.memberReader.On("GetMember", mock.Anything, tx, teamID, userID).
		Return(member, nil).Once()

	got, err := s.repo.GetMember(s.ctx, tx, teamID, userID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMember_NotFound() {
	teamID := uuid.New()
	userID := uuid.New()
	s.memberReader.On("GetMember", mock.Anything, mock.Anything, teamID, userID).
		Return(model.TeamMember{}, model.ErrMemberNotFound).Once()

	got, err := s.repo.GetMember(s.ctx, nil, teamID, userID)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrMemberNotFound)
	assert.Equal(s.T(), model.TeamMember{}, got)
	s.memberReader.AssertExpectations(s.T())
}
