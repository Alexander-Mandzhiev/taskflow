package team

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func (s *teamService) GetByID(ctx context.Context, teamID, userID uuid.UUID) (*model.TeamWithMembers, error) {
	var result *model.TeamWithMembers
	if err := s.txManager.WithTx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		team, errTx := s.teamRepo.GetByID(ctx, tx, teamID)
		if errTx != nil {
			return errTx
		}
		members, errTx := s.memberRepo.GetMembersByTeamID(ctx, tx, teamID)
		if errTx != nil {
			return errTx
		}
		_, errTx = s.memberRepo.GetMember(ctx, tx, teamID, userID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrMemberNotFound) {
				return model.ErrForbidden
			}
			return errTx
		}
		result = &model.TeamWithMembers{Team: *team, Members: members}
		return nil
	}); err != nil {
		if !errors.Is(err, model.ErrForbidden) {
			logger.Error(ctx, "GetByID failed", zap.Error(err))
		}
		return nil, err
	}

	return result, nil
}
