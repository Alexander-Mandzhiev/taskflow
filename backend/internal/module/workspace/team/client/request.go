package client

import "time"

// SendInvitationRequest — модель запроса, уходящего из клиента в сервис уведомлений (email/delivery).
// AcceptLink в проде не передаётся: сервис уведомлений собирает ссылку из Token и своего base URL.
type SendInvitationRequest struct {
	Email       string
	TeamName    string
	InviterName string
	Token       string
	ExpiresAt   time.Time
}
