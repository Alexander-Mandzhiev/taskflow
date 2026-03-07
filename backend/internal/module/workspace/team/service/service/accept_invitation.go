package service

import (
	"context"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/google/uuid"
)

// AcceptInvitation — заглушка: принятие приглашения по токену пока не реализовано.
func (s *teamService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) (*model2.TeamMember, error) {
	return nil, model2.ErrNotImplemented
}
