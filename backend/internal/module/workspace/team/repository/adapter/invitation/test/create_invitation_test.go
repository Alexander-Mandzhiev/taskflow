package invitation_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

func (s *AdapterSuite) TestCreateInvitation_Success() {
	inv := &model.TeamInvitation{
		TeamID:    uuid.New(),
		Email:     "user@example.com",
		Token:     "token123",
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.invitationWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), inv).
		Return(nil).Once()

	err := s.repo.CreateInvitation(s.ctx, nil, inv)

	assert.NoError(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateInvitation_WithTx() {
	tx := &sqlx.Tx{}
	inv := &model.TeamInvitation{
		TeamID: uuid.New(),
		Email:  "a@b.com",
		Token:  "t",
		Status: model.InvitationStatusPending,
	}

	s.invitationWriter.On("Create", mock.Anything, tx, inv).
		Return(nil).Once()

	err := s.repo.CreateInvitation(s.ctx, tx, inv)

	assert.NoError(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateInvitation_WriterError() {
	inv := &model.TeamInvitation{TeamID: uuid.New(), Email: "x@y.com", Status: model.InvitationStatusPending}

	s.invitationWriter.On("Create", mock.Anything, mock.Anything, inv).
		Return(assert.AnError).Once()

	err := s.repo.CreateInvitation(s.ctx, nil, inv)

	assert.Error(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}
