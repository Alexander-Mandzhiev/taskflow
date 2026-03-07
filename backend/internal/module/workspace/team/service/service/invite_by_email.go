package service

import (
	"context"
	"errors"
	"time"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// TODO: вынести в конфиг (например notification.invitation_ttl или app.invitation_expires_in).
const invitationExpiresIn = 7 * 24 * time.Hour

// InviteByEmail создаёт приглашение (запись в team_invitations). Проверки и данные для уведомления — в одной транзакции; отправка notifier — после коммита.
func (s *teamService) InviteByEmail(ctx context.Context, teamID, inviterUserID uuid.UUID, inviteeEmail, role string) (*model2.TeamInvitation, error) {
	var inv *model2.TeamInvitation
	var teamName, inviterName string

	err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := s.checkInviterPermissions(ctx, tx, teamID, inviterUserID); err != nil {
			return err
		}
		if err := s.validateInviteRole(role); err != nil {
			return err
		}
		if err := s.checkUserNotMember(ctx, tx, teamID, inviteeEmail); err != nil {
			return err
		}
		if err := s.checkNoActiveInvitation(ctx, tx, teamID, inviteeEmail); err != nil {
			return err
		}

		inv = &model2.TeamInvitation{
			ID:        uuid.New(),
			TeamID:    teamID,
			Email:     inviteeEmail,
			Role:      role,
			InvitedBy: inviterUserID,
			Status:    model2.InvitationStatusPending,
			Token:     uuid.New().String(),
			ExpiresAt: time.Now().UTC().Add(invitationExpiresIn),
		}
		if err := s.repo.CreateInvitation(ctx, tx, inv); err != nil {
			return err
		}

		s.prepareNotificationData(ctx, tx, teamID, inviterUserID, &teamName, &inviterName)
		return nil
	})
	if err != nil {
		if errors.Is(err, model2.ErrForbidden) || errors.Is(err, model2.ErrAlreadyMember) ||
			errors.Is(err, model2.ErrAlreadyInvited) || errors.Is(err, model2.ErrInvalidRole) {
			return nil, err
		}
		logger.Error(ctx, "InviteByEmail failed", zap.Error(err))
		return nil, err
	}

	if s.notifier != nil {
		if err := s.notifier.NotifyInvitation(ctx, inv, teamName, inviterName); err != nil {
			logger.Error(ctx, "notify invitation failed", zap.Error(err))
		}
	}
	return inv, nil
}

func (s *teamService) checkInviterPermissions(ctx context.Context, tx *sqlx.Tx, teamID, inviterUserID uuid.UUID) error {
	member, err := s.repo.GetMember(ctx, tx, teamID, inviterUserID)
	if err != nil {
		if errors.Is(err, model2.ErrMemberNotFound) {
			return model2.ErrForbidden
		}
		return err
	}
	if member.Role != model2.RoleOwner {
		return model2.ErrForbidden
	}
	return nil
}

func (s *teamService) validateInviteRole(role string) error {
	if role != model2.RoleMember && role != model2.RoleAdmin {
		return model2.ErrInvalidRole
	}
	return nil
}

func (s *teamService) checkUserNotMember(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, inviteeEmail string) error {
	user, err := s.userRepo.GetByEmail(ctx, tx, inviteeEmail)
	if err != nil {
		if !errors.Is(err, usermodel.ErrUserNotFound) {
			return err
		}
		return nil
	}
	existing, err := s.repo.GetMember(ctx, tx, teamID, user.ID)
	if err != nil && !errors.Is(err, model2.ErrMemberNotFound) {
		return err
	}
	if existing != nil {
		return model2.ErrAlreadyMember
	}
	return nil
}

func (s *teamService) checkNoActiveInvitation(ctx context.Context, tx *sqlx.Tx, teamID uuid.UUID, inviteeEmail string) error {
	pending, err := s.repo.GetPendingInvitationByTeamAndEmail(ctx, tx, teamID, inviteeEmail)
	if err != nil && !errors.Is(err, model2.ErrInvitationNotFound) {
		return err
	}
	if pending != nil && pending.ExpiresAt.After(time.Now().UTC()) {
		return model2.ErrAlreadyInvited
	}
	return nil
}

func (s *teamService) prepareNotificationData(ctx context.Context, tx *sqlx.Tx, teamID, inviterUserID uuid.UUID, teamName, inviterName *string) {
	if team, err := s.repo.GetByID(ctx, tx, teamID); err == nil {
		*teamName = team.Name
	} else {
		logger.Warn(ctx, "InviteByEmail: get team for notification failed", zap.Error(err))
	}
	if u, err := s.userRepo.GetByID(ctx, tx, inviterUserID.String()); err == nil && u != nil {
		*inviterName = u.Name
		if *inviterName == "" {
			*inviterName = u.Email
		}
	} else if err != nil {
		logger.Warn(ctx, "InviteByEmail: get inviter for notification failed", zap.Error(err))
	}
}
