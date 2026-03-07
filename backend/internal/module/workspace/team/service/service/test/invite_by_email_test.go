package service_test

import (
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
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
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model2.TeamInvitation)(nil), model2.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.MatchedBy(func(inv *model2.TeamInvitation) bool {
		return inv != nil && inv.TeamID == teamID && inv.Email == inviteeEmail && inv.Role == model2.RoleMember && inv.InvitedBy == inviterID
	})).Return(nil).Once()
	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).
		Return(&model2.Team{ID: teamID, Name: "Test Team"}, nil).Maybe()
	s.userRepo.On("GetByID", mock.Anything, mock.Anything, testInviterID).
		Return(&usermodel.User{ID: inviterID, Email: "owner@example.com", Name: "Owner"}, nil).Maybe()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), inv)
	assert.Equal(s.T(), teamID, inv.TeamID)
	assert.Equal(s.T(), inviteeEmail, inv.Email)
	assert.Equal(s.T(), model2.RoleMember, inv.Role)
	assert.Equal(s.T(), model2.InvitationStatusPending, inv.Status)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Forbidden_NotOwner() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	adminMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleAdmin,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(adminMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrForbidden)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Forbidden_MemberNotFound() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return((*model2.TeamMember)(nil), model2.ErrMemberNotFound).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrForbidden)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_InvalidRole() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, "owner")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrInvalidRole)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_AlreadyMember() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	inviteeID := uuid.MustParse(inviteeUserIDStr)
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}
	existingUser := &usermodel.User{ID: inviteeID, Email: inviteeEmail, Name: "Invited"}
	existingMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviteeID,
		TeamID:    teamID,
		Role:      model2.RoleMember,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return(existingUser, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviteeID).
		Return(existingMember, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrAlreadyMember)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_AlreadyInvited() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}
	pendingInv := &model2.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    teamID,
		Email:     inviteeEmail,
		Role:      model2.RoleMember,
		Status:    model2.InvitationStatusPending,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return(pendingInv, nil).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, model2.ErrAlreadyInvited)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_CreateInvitationError() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return((*usermodel.User)(nil), usermodel.ErrUserNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model2.TeamInvitation)(nil), model2.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.Anything).
		Return(assert.AnError).Once()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleMember)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), inv)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestInviteByEmail_Success_UserRegisteredNotInTeam() {
	teamID := uuid.MustParse(testTeamID)
	inviterID := uuid.MustParse(testInviterID)
	inviteeID := uuid.MustParse(inviteeUserIDStr)
	ownerMember := &model2.TeamMember{
		ID:        uuid.New(),
		UserID:    inviterID,
		TeamID:    teamID,
		Role:      model2.RoleOwner,
		CreatedAt: time.Now(),
	}
	existingUser := &usermodel.User{ID: inviteeID, Email: inviteeEmail, Name: "Invited"}

	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviterID).
		Return(ownerMember, nil).Once()
	s.userRepo.On("GetByEmail", mock.Anything, mock.Anything, inviteeEmail).
		Return(existingUser, nil).Once()
	s.teamRepo.On("GetMember", mock.Anything, mock.Anything, teamID, inviteeID).
		Return((*model2.TeamMember)(nil), model2.ErrMemberNotFound).Once()
	s.teamRepo.On("GetPendingInvitationByTeamAndEmail", mock.Anything, mock.Anything, teamID, inviteeEmail).
		Return((*model2.TeamInvitation)(nil), model2.ErrInvitationNotFound).Once()
	s.teamRepo.On("CreateInvitation", mock.Anything, mock.Anything, mock.MatchedBy(func(inv *model2.TeamInvitation) bool {
		return inv != nil && inv.TeamID == teamID && inv.Email == inviteeEmail && inv.InvitedBy == inviterID
	})).Return(nil).Once()
	s.teamRepo.On("GetByID", mock.Anything, mock.Anything, teamID).Return(&model2.Team{ID: teamID, Name: "Test Team"}, nil).Maybe()
	s.userRepo.On("GetByID", mock.Anything, mock.Anything, testInviterID).Return(&usermodel.User{ID: inviterID, Email: "owner@example.com", Name: "Owner"}, nil).Maybe()

	inv, err := s.svc.InviteByEmail(s.ctx, teamID, inviterID, inviteeEmail, model2.RoleAdmin)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), inv)
	assert.Equal(s.T(), model2.RoleAdmin, inv.Role)
	s.teamRepo.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}
