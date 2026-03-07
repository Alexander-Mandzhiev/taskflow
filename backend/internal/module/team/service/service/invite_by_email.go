package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

const invitationExpiresIn = 7 * 24 * time.Hour

// InviteByEmail создаёт приглашение (запись в team_invitations). Проверки: права owner/admin, пользователь не в команде, нет pending по (team_id, email). Отправка письма — позже через отдельный сервис.
func (s *teamService) InviteByEmail(ctx context.Context, teamID, inviterUserID uuid.UUID, inviteeEmail, role string) (*model.TeamInvitation, error) {
	member, err := s.repo.GetMember(ctx, nil, teamID.String(), inviterUserID.String())
	if err != nil {
		logger.Error(ctx, "InviteByEmail GetMember failed", zap.Error(err))
		return nil, err
	}
	if member.Role != model.RoleOwner && member.Role != model.RoleAdmin {
		return nil, model.ErrForbidden
	}

	// Если по email есть пользователь — проверяем, что он ещё не в команде.
	user, err := s.userRepo.GetByEmail(ctx, nil, inviteeEmail)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			// Пользователь не зарегистрирован — приглашение всё равно создаём (примет после регистрации).
		} else {
			return nil, err
		}
	} else {
		existing, err := s.repo.GetMember(ctx, nil, teamID.String(), user.ID.String())
		if err != nil && !errors.Is(err, model.ErrMemberNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, model.ErrAlreadyMember
		}
	}

	pending, err := s.repo.GetPendingInvitationByTeamAndEmail(ctx, nil, teamID.String(), inviteeEmail)
	if err != nil && !errors.Is(err, model.ErrInvitationNotFound) {
		return nil, err
	}
	if pending != nil {
		return nil, model.ErrAlreadyInvited
	}

	inv := &model.TeamInvitation{
		ID:        uuid.New(),
		TeamID:    teamID,
		Email:     inviteeEmail,
		Role:      role,
		InvitedBy: inviterUserID,
		Status:    model.InvitationStatusPending,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(invitationExpiresIn),
	}

	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		return s.repo.CreateInvitation(ctx, tx, inv)
	}); err != nil {
		logger.Error(ctx, "InviteByEmail CreateInvitation failed", zap.Error(err))
		return nil, err
	}

	return inv, nil
}
