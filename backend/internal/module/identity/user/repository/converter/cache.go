package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/resources"
)

// ToCache преобразует доменную модель User в модель для кеша.
// PasswordHash не сохраняется в кеш из соображений безопасности.
func ToCache(user model.User) resources.UserCache {
	return resources.UserCache{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

// FromCache преобразует модель из кеша в доменную User.
// PasswordHash будет пустым — для операций с паролем используется БД напрямую.
func FromCache(c resources.UserCache) (model.User, error) {
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:        id,
		Email:     c.Email,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}, nil
}
