package adapter_test

import (
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *AdapterSuite) TestCreateInvitation_Success() {
	inv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Email:     "invited@example.com",
		Role:      model2.RoleMember,
		InvitedBy: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Status:    model2.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.invitationWriter.On("Create", mock.Anything, (*sqlx.Tx)(nil), inv).
		Return(nil).Once()

	err := s.repo.CreateInvitation(s.ctx, nil, inv)

	assert.NoError(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateInvitation_WithTx() {
	tx := &sqlx.Tx{}
	inv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Email:     "admin@example.com",
		Role:      model2.RoleAdmin,
		InvitedBy: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Status:    model2.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.invitationWriter.On("Create", mock.Anything, tx, inv).
		Return(nil).Once()

	err := s.repo.CreateInvitation(s.ctx, tx, inv)

	assert.NoError(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}

func (s *AdapterSuite) TestCreateInvitation_WriterError() {
	inv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Email:     "invited@example.com",
		Role:      model2.RoleMember,
		InvitedBy: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Status:    model2.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.invitationWriter.On("Create", mock.Anything, mock.Anything, inv).
		Return(assert.AnError).Once()

	err := s.repo.CreateInvitation(s.ctx, nil, inv)

	assert.Error(s.T(), err)
	s.invitationWriter.AssertExpectations(s.T())
}
