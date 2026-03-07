package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	usermodel "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// InviteByEmail приглашает по email: проверка прав по inviterUserID (owner/admin), резолв inviteeEmail, добавление в команду.
// Позже флоу будет: invite = запись в team_invitations + письмо; addMember только при принятии приглашения.
func (s *teamService) InviteByEmail(ctx context.Context, teamID, inviterUserID uuid.UUID, inviteeEmail, role string) (*model.TeamMember, error) {
	member, err := s.repo.GetMember(ctx, nil, teamID.String(), inviterUserID.String())
	if err != nil {
		return nil, err
	}
	if member.Role != model.RoleOwner && member.Role != model.RoleAdmin {
		return nil, model.ErrForbidden
	}
	user, err := s.userRepo.GetByEmail(ctx, nil, inviteeEmail)
	if err != nil {
		if errors.Is(err, usermodel.ErrUserNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	inviteeUserID := user.ID
	var added *model.TeamMember
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		added, errTx = s.repo.AddMember(ctx, tx, teamID.String(), inviteeUserID.String(), role)
		return errTx
	}); err != nil {
		logger.Error(ctx, "InviteByEmail AddMember failed", zap.Error(err))
		return nil, err
	}
	return added, nil
}
