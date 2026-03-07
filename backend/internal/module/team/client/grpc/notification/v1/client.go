package v1

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/client/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/client/grpc"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
)

var _ grpc.Notification = (*Client)(nil)

// Client — мок-реализация grpc.Notification. Не выполняет реальный вызов к сервису уведомлений;

// домен конвертируется в SendInvitationRequest, но запрос никуда не отправляется.

// После подключения generated gRPC client здесь будет вызов notificationv1.NotificationServiceClient.

type Client struct {
	// generatedClient notificationv1.NotificationServiceClient
}

// NewClient создаёт мок-клиент. После добавления proto передавать generatedClient и заменить noop на реальный вызов.

func NewClient( /* generatedClient notificationv1.NotificationServiceClient */ ) *Client {
	return &Client{}
}

// NotifyInvitation конвертирует домен в client.SendInvitationRequest. Мок: запрос не отправляется.

// AcceptLink в проде собирается на стороне сервиса уведомлений из inv.Token и своего base URL.

func (c *Client) NotifyInvitation(ctx context.Context, inv *model.TeamInvitation, teamName, inviterName string) error {
	_ = ctx

	req := converter.ToSendInvitationRequest(inv, teamName, inviterName)

	_ = req

	return nil
}
