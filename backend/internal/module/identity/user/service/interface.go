package service

import (
	"context"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
)

// UserService — слой сервиса пользователей (регистрация, профиль, CRUD).
// Транзакции открываются внутри сервиса (txmanager.WithTx); вызывающий не передаёт tx.
// Create и Update: сервис проверяет input на nil до вызова репозитория и при nil возвращает model.ErrNilInput.
type UserService interface {
	Create(ctx context.Context, input *model.UserInput, passwordHash string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, input *model.UserInput) (*model.User, error)
	ChangePassword(ctx context.Context, id, passwordHash string) error
	Delete(ctx context.Context, id string) error
}
