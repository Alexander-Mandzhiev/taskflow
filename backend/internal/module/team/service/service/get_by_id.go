package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) GetByID(ctx context.Context, teamID uuid.UUID) (*model.TeamWithMembers, error) {
	var team *model.TeamWithMembers
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		var errTx error
		team, errTx = s.repo.GetByID(ctx, tx, teamID.String())
		return errTx
	}); err != nil {
		logger.Error(ctx, "GetByID failed", zap.Error(err))
		return nil, err
	}

	return team, nil
}
