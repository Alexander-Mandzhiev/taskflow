package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// ToDomainTeam преобразует строку БД (TeamRow) в доменную модель Team.
func ToDomainTeam(r resources.TeamRow) (model.Team, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.Team{}, err
	}
	createdBy, err := uuid.Parse(r.CreatedBy)
	if err != nil {
		return model.Team{}, err
	}
	return model.Team{
		ID:        id,
		Name:      r.Name,
		CreatedBy: createdBy,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}, nil
}

// ToDomainTeamWithRole преобразует строку БД (TeamWithRoleRow) в доменную модель TeamWithRole.
func ToDomainTeamWithRole(r resources.TeamWithRoleRow) (model.TeamWithRole, error) {
	team, err := ToDomainTeam(r.TeamRow)
	if err != nil {
		return model.TeamWithRole{}, err
	}
	return model.TeamWithRole{
		Team: team,
		Role: r.Role,
	}, nil
}

// ToRepoTeamInput преобразует доменный TeamInput в ресурс репозитория.
func ToRepoTeamInput(m model.TeamInput) resources.TeamInput {
	return resources.TeamInput{
		Name: m.Name,
	}
}
