package converter

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/client"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

// ToSendInvitationRequest конвертирует доменную модель приглашения и данные для письма в клиентскую модель запроса.
// Ссылку принять приглашение собирает сервис уведомлений из Token и своего конфига.
func ToSendInvitationRequest(inv *model.TeamInvitation, teamName, inviterName string) *client.SendInvitationRequest {
	if inv == nil {
		return nil
	}
	return &client.SendInvitationRequest{
		Email:       inv.Email,
		TeamName:    teamName,
		InviterName: inviterName,
		Token:       inv.Token,
		ExpiresAt:   inv.ExpiresAt,
	}
}
