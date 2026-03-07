package grpc

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// Notification — клиент сервиса уведомлений (email/delivery). Принимает доменные типы; реализация в notification/v1 конвертирует в proto.

// Ссылку «принять приглашение» собирает сам сервис уведомлений из inv.Token и своего конфига (base URL).

type Notification interface {
	// NotifyInvitation отправляет уведомление о приглашении. inv (в т.ч. Token), teamName, inviterName — данные для письма.

	NotifyInvitation(ctx context.Context, inv *model.TeamInvitation, teamName, inviterName string) error
}
