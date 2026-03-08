package member_test

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestGetMembersByTeamID_Success() {
	teamID := uuid.New()
	want := []*model.TeamMember{
		{TeamID: teamID, UserID: uuid.New(), Role: model.RoleOwner},
		{TeamID: teamID, UserID: uuid.New(), Role: model.RoleMember},
	}

	s.memberReader.On("GetByTeamID", mock.Anything, (*sqlx.Tx)(nil), teamID).
		Return(want, nil).Once()

	got, err := s.repo.GetMembersByTeamID(s.ctx, nil, teamID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), want, got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMembersByTeamID_Empty() {
	teamID := uuid.New()
	s.memberReader.On("GetByTeamID", mock.Anything, mock.Anything, teamID).
		Return([]*model.TeamMember{}, nil).Once()

	got, err := s.repo.GetMembersByTeamID(s.ctx, nil, teamID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Empty(s.T(), got)
	s.memberReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetMembersByTeamID_ReaderError() {
	teamID := uuid.New()
	s.memberReader.On("GetByTeamID", mock.Anything, mock.Anything, teamID).
		Return(([]*model.TeamMember)(nil), assert.AnError).Once()

	got, err := s.repo.GetMembersByTeamID(s.ctx, &sqlx.Tx{}, teamID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.memberReader.AssertExpectations(s.T())
}
