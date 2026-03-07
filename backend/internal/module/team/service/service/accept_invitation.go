package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// AcceptInvitation — заглушка: принятие приглашения по токену пока не реализовано.
func (s *teamService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) (*model.TeamMember, error) {
	return nil, model.ErrNotImplemented
}
