package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/resources"
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
