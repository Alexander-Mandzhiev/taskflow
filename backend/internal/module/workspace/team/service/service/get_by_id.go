package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) GetByID(ctx context.Context, teamID, userID uuid.UUID) (*model2.TeamWithMembers, error) {
	var result *model2.TeamWithMembers
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		team, errTx := s.repo.GetByID(ctx, tx, teamID)
		if errTx != nil {
			return errTx
		}
		members, errTx := s.repo.GetMembersByTeamID(ctx, tx, teamID)
		if errTx != nil {
			return errTx
		}
		_, errTx = s.repo.GetMember(ctx, tx, teamID, userID)
		if errTx != nil {
			if errors.Is(errTx, model2.ErrMemberNotFound) {
				return model2.ErrForbidden
			}
			return errTx
		}
		result = &model2.TeamWithMembers{Team: *team, Members: members}
		return nil
	}); err != nil {
		if !errors.Is(err, model2.ErrForbidden) {
			logger.Error(ctx, "GetByID failed", zap.Error(err))
		}
		return nil, err
	}

	return result, nil
}
