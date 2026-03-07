package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// ToDomainTeamInvitation преобразует строку БД (TeamInvitationRow) в доменную модель TeamInvitation.
func ToDomainTeamInvitation(r resources.TeamInvitationRow) (model.TeamInvitation, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.TeamInvitation{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamInvitation{}, err
	}
	invitedBy, err := uuid.Parse(r.InvitedBy)
	if err != nil {
		return model.TeamInvitation{}, err
	}
	return model.TeamInvitation{
		ID:        id,
		TeamID:    teamID,
		Email:     r.Email,
		Role:      r.Role,
		InvitedBy: invitedBy,
		Status:    r.Status,
		Token:     r.Token,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}
