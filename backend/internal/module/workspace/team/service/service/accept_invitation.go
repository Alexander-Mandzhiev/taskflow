package service

import (
	"context"

	"github.com/google/uuid"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// AcceptInvitation — заглушка: принятие приглашения по токену пока не реализовано.
func (s *teamService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) (*model2.TeamMember, error) {
	return nil, model2.ErrNotImplemented
}
