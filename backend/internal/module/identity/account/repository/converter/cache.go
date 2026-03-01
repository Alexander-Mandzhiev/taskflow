package converter

import (
	"github.com/google/uuid"

	"mkk/internal/module/identity/account/model"
	"mkk/internal/module/identity/account/repository/resources"
)

// ToCache преобразует доменную модель Session в модель для записи в кеш.
func ToCache(session *model.Session) resources.SessionCache {
	if session == nil {
		return resources.SessionCache{}
	}
	return resources.SessionCache{
		UserID:     session.UserID.String(),
		CreatedAt:  session.CreatedAt,
		DeviceType: session.DeviceType,
		UserAgent:  session.UserAgent,
		IP:         session.IP,
	}
}

// FromCache преобразует модель из кеша в доменную Session.
func FromCache(c resources.SessionCache) (*model.Session, error) {
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		UserID:     userID,
		CreatedAt:  c.CreatedAt,
		DeviceType: c.DeviceType,
		UserAgent:  c.UserAgent,
		IP:         c.IP,
	}, nil
}
