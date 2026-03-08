package member_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestAddMember_Success() {
	teamID := uuid.New()
	userID := uuid.New()
	role := model.RoleMember
	member := &model.TeamMember{TeamID: teamID, UserID: userID, Role: role}

	s.memberWriter.On("AddMember", mock.Anything, (*sqlx.Tx)(nil), teamID, userID, role).
		Return(member, nil).Once()

	got, err := s.repo.AddMember(s.ctx, nil, teamID, userID, role)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestAddMember_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.New()
	userID := uuid.New()
	role := model.RoleAdmin
	member := &model.TeamMember{TeamID: teamID, UserID: userID, Role: role}

	s.memberWriter.On("AddMember", mock.Anything, tx, teamID, userID, role).
		Return(member, nil).Once()

	got, err := s.repo.AddMember(s.ctx, tx, teamID, userID, role)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), member, got)
	s.memberWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestAddMember_WriterError() {
	teamID := uuid.New()
	userID := uuid.New()
	role := model.RoleMember

	s.memberWriter.On("AddMember", mock.Anything, mock.Anything, teamID, userID, role).
		Return((*model.TeamMember)(nil), assert.AnError).Once()

	got, err := s.repo.AddMember(s.ctx, nil, teamID, userID, role)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.memberWriter.AssertExpectations(s.T())
}
