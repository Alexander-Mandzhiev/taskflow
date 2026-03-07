package converter

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

const timeFormat = time.RFC3339

// TeamToResponse конвертирует доменную команду в DTO ответа.
func TeamToResponse(t *model.Team) dto.TeamResponse {
	if t == nil {
		return dto.TeamResponse{}
	}
	return dto.TeamResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		CreatedBy: t.CreatedBy.String(),
		CreatedAt: t.CreatedAt.Format(timeFormat),
		UpdatedAt: t.UpdatedAt.Format(timeFormat),
	}
}

// MemberToResponse конвертирует доменного участника в DTO ответа.
func MemberToResponse(m *model.TeamMember) dto.MemberResponse {
	if m == nil {
		return dto.MemberResponse{}
	}
	return dto.MemberResponse{
		ID:        m.ID.String(),
		UserID:    m.UserID.String(),
		TeamID:    m.TeamID.String(),
		Role:      m.Role,
		CreatedAt: m.CreatedAt.Format(timeFormat),
	}
}

// TeamWithMembersToResponse конвертирует доменную команду с участниками в DTO ответа.
func TeamWithMembersToResponse(t *model.TeamWithMembers) dto.TeamWithMembersResponse {
	if t == nil {
		return dto.TeamWithMembersResponse{}
	}
	members := make([]dto.MemberResponse, 0, len(t.Members))
	for _, m := range t.Members {
		members = append(members, MemberToResponse(m))
	}
	return dto.TeamWithMembersResponse{
		Team:    TeamToResponse(&t.Team),
		Members: members,
	}
}

// TeamWithRoleToResponse конвертирует доменную команду с ролью в DTO ответа.
func TeamWithRoleToResponse(t *model.TeamWithRole) dto.TeamWithRoleResponse {
	if t == nil {
		return dto.TeamWithRoleResponse{}
	}
	return dto.TeamWithRoleResponse{
		TeamResponse: TeamToResponse(&t.Team),
		Role:         t.Role,
	}
}

// TeamsWithRolesToResponse конвертирует список команд с ролями в DTO ответа.
func TeamsWithRolesToResponse(teams []*model.TeamWithRole) []dto.TeamWithRoleResponse {
	if len(teams) == 0 {
		return []dto.TeamWithRoleResponse{}
	}
	out := make([]dto.TeamWithRoleResponse, 0, len(teams))
	for _, t := range teams {
		out = append(out, TeamWithRoleToResponse(t))
	}
	return out
}

// InvitationToResponse конвертирует приглашение в DTO ответа (без token в JSON).
func InvitationToResponse(inv *model.TeamInvitation) dto.InvitationResponse {
	if inv == nil {
		return dto.InvitationResponse{}
	}
	return dto.InvitationResponse{
		ID:        inv.ID.String(),
		TeamID:    inv.TeamID.String(),
		Email:     inv.Email,
		Role:      inv.Role,
		ExpiresAt: inv.ExpiresAt.Format(timeFormat),
	}
}
