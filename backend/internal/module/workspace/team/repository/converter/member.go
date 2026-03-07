package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// ToDomainTeamMember преобразует строку БД (TeamMemberRow) в доменную модель TeamMember.
func ToDomainTeamMember(r resources.TeamMemberRow) (model.TeamMember, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.TeamMember{}, err
	}
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.TeamMember{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamMember{}, err
	}
	return model.TeamMember{
		ID:        id,
		UserID:    userID,
		TeamID:    teamID,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}, nil
}
