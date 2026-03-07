package service_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

const (
	testTeamID       = "550e8400-e29b-41d4-a716-446655440001"
	testInviterID    = "550e8400-e29b-41d4-a716-446655440000"
	inviteeEmail     = "invited@example.com"
	inviteeUserIDStr = "770e8400-e29b-41d4-a716-446655440002"
)

func (s *ServiceSuite) TestInviteByEmail_Success() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model.TeamInvitation)(nil), model.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.MatchedBy(func(inv *model.TeamInvitation) bool {
		return inv != nil && inv.TeamID == teamID && inv.Email == inviteeEmail && inv.Role == model.RoleMember && inv.InvitedBy == inviterID
	})).Return(nil).Once()
	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).
		Return(&model.Team{ID: teamID, Name: "Test Team"}, nil).Maybe()
	s.userRepo.On("GetByID", mock.Anything, mock.Anything, testInviterID).
		Return(&usermodel.User{ID: inviterID, Email: "owner@example.com", Name: "Owner"}, nil).Maybe()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), inv)
	assert.Equal(s.T(), teamID, inv.TeamID)
	assert.Equal(s.T(), inviteeEmail, inv.Email)
	assert.Equal(s.T(), model.RoleMember, inv.Role)
	assert.Equal(s.T(), model.InvitationStatusPending, inv.Status)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Forbidden_NotOwner() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	adminMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleAdmin,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(adminMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrForbidden)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Forbidden_MemberNotFound() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return((*model.TeamMember)(nil), model.ErrMemberNotFound).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrForbidden)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_InvalidRole() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, "owner")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrInvalidRole)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_AlreadyMember() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	inviteeID := uuid.MustParse(inviteeUserIDStr)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}
	existingUser := &usermodel.User{ID: inviteeID, Email: inviteeEmail, Name: "Invited"}
	existingMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviteeID,
		TeamID:    teamID,
		Role:      model.RoleMember,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return(existingUser, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviteeID).
		Return(existingMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrAlreadyMember)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_AlreadyInvited() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}
	pendingInv := &model.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    teamID,
		Email:     inviteeEmail,
		Role:      model.RoleMember,
		Status:    model.InvitationStatusPending,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return(pendingInv, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model.ErrAlreadyInvited)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_CreateInvitationError() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model.TeamInvitation)(nil), model.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.Anything).
		Return(assert.AnError).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleMember)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Success_UserRegisteredNotInTeam() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	inviteeID := uuid.MustParse(inviteeUserIDStr)
	ownerMember := &model.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model.RoleOwner,
		CreatedAt: time.Now(),
	}
	existingUser := &usermodel.User{ID: inviteeID, Email: inviteeEmail, Name: "Invited"}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return(existingUser, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviteeID).
		Return((*model.TeamMember)(nil), model.ErrMemberNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model.TeamInvitation)(nil), model.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.MatchedBy(func(inv *model.TeamInvitation) bool {
		return inv != nil && inv.TeamID == teamID && inv.Email == inviteeEmail && inv.InvitedBy == inviterID
	})).Return(nil).Once()
	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(&model.Team{ID: teamID, Name: "Test Team"}, nil).Maybe()
	s.userRepo.On("GetByID", mock.Anything, mock.Anything, testInviterID).Return(&usermodel.User{ID: inviterID, Email: "owner@example.com", Name: "Owner"}, nil).Maybe()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model.RoleAdmin)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), inv)
	assert.Equal(s.T(), model.RoleAdmin, inv.Role)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}
