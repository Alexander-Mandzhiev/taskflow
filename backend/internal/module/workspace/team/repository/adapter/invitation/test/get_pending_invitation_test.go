package invitation_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_Success() {
	teamID := uuid.New()
	email := "user@example.com"
	inv := model.TeamInvitation{
		TeamID:    teamID,
		Email:     email,
		Token:     "token",
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, (*sqlx.Tx)(nil), teamID, email).
		Return(inv, nil).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, nil, teamID, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), inv, got)
	s.invitationReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_WithTx() {
	tx := &sqlx.Tx{}
	teamID := uuid.New()
	email := "a@b.com"
	inv := model.TeamInvitation{TeamID: teamID, Email: email, Status: model.InvitationStatusPending}

	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, tx, teamID, email).
		Return(inv, nil).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, tx, teamID, email)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), inv, got)
	s.invitationReader.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestGetPendingInvitationByTeamAndEmail_NotFound() {
	teamID := uuid.New()
	email := "nobody@example.com"
	s.invitationReader.On("GetPendingByTeamAndEmail", mock.Anything, mock.Anything, teamID, email).
		Return(model.TeamInvitation{}, model.ErrInvitationNotFound).Once()

	got, err := s.repo.GetPendingInvitationByTeamAndEmail(s.ctx, nil, teamID, email)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvitationNotFound)
	assert.Equal(s.T(), model.TeamInvitation{}, got)
	s.invitationReader.AssertExpectations(s.T())
}
