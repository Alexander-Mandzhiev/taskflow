package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/resources"
)

// ToDomainUser преобразует строку БД (UserRow) в доменную модель User.
func ToDomainUser(r resources.UserRow) (model.User, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:           id,
		Email:        r.Email,
		Name:         r.Name,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		DeletedAt:    r.DeletedAt,
	}, nil
}
