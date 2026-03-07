package adapter_test

import (
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_Success() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	email := "invited@example.com"
	inv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    teamID,
		Email:     email,
		Role:      model2.RoleMember,
		InvitedBy: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Status:    model2.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, (*sqlx.Tx)(nil), teamID.String(), email).
		Return(inv, nil).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, nil, teamID, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), inv, got)
	assert.Equal(s.T(), email, got.Email)
	assert.Equal(s.T(), model2.InvitationStatusPending, got.Status)
	s.invitationReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	email := "pending@example.com"
	inv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    teamID,
		Email:     email,
		Role:      model2.RoleAdmin,
		InvitedBy: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Status:    model2.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(48 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, tx, teamID.String(), email).
		Return(inv, nil).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, tx, teamID, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), inv, got)
	s.invitationReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_NotFound() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	email := "unknown@example.com"

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, mock.Anything, teamID.String(), email).
		Return((*model2.TeamInvitation)(nil), model2.ErrInvitationNotFound).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, nil, teamID, email)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrInvitationNotFound)
	assert.Nil(s.T(), got)
	s.invitationReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_ReaderError() {
	teamID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	email := "invited@example.com"

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, mock.Anything, teamID.String(), email).
		Return((*model2.TeamInvitation)(nil), assert.AnError).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, nil, teamID, email)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
	s.invitationReader.AssertExpectations(s.T())
}
