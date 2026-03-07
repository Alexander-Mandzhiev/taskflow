package converter

import (
	"github.com/google/uuid"

	model2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	resources2 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/resources"
)

// ToDomainTeam преобразует строку БД (TeamRow) в доменную модель Team.
func ToDomainTeam(r resources2.TeamRow) (model2.Team, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model2.Team{}, err
	}
	createdBy, err := uuid.Parse(r.CreatedBy)
	if err != nil {
		return model2.Team{}, err
	}
	return model2.Team{
		ID:        id,
		Name:      r.Name,
		CreatedBy: createdBy,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}, nil
}

// ToDomainTeamMember преобразует строку БД (TeamMemberRow) в доменную модель TeamMember.
func ToDomainTeamMember(r resources2.TeamMemberRow) (model2.TeamMember, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model2.TeamMember{}, err
	}
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model2.TeamMember{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model2.TeamMember{}, err
	}
	return model2.TeamMember{
		ID:        id,
		UserID:    userID,
		TeamID:    teamID,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}, nil
}

// ToDomainTeamWithRole преобразует строку БД (TeamWithRoleRow) в доменную модель TeamWithRole.
func ToDomainTeamWithRole(r resources2.TeamWithRoleRow) (model2.TeamWithRole, error) {
	team, err := ToDomainTeam(r.TeamRow)
	if err != nil {
		return model2.TeamWithRole{}, err
	}
	return model2.TeamWithRole{
		Team: team,
		Role: r.Role,
	}, nil
}

// ToDomainTeamInvitation преобразует строку БД (TeamInvitationRow) в доменную модель TeamInvitation.
func ToDomainTeamInvitation(r resources2.TeamInvitationRow) (model2.TeamInvitation, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model2.TeamInvitation{}, err
	}
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model2.TeamInvitation{}, err
	}
	invitedBy, err := uuid.Parse(r.InvitedBy)
	if err != nil {
		return model2.TeamInvitation{}, err
	}
	return model2.TeamInvitation{
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
