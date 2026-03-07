package service

import (
	"context"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) Create(ctx context.Context, input *model2.TeamInput, ownerUserID uuid.UUID) (*model2.Team, error) {
	if input == nil {
		logger.Warn(ctx, "Create team: nil input")
		return nil, model2.ErrNilInput
	}

	var team *model2.Team
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		team, errTx = s.repo.Create(ctx, tx, input, ownerUserID)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.repo.AddMember(ctx, tx, team.ID, ownerUserID, model2.RoleOwner)
		return errTx
	}); err != nil {
		logger.Error(ctx, "Create team failed", zap.Error(err))
		return nil, err
	}
	return team, nil
}
