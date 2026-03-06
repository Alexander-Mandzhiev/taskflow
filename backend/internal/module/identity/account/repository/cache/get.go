package cache

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/repository/resources"
)

// Get возвращает сессию из кеша по sessionID.
// При отсутствии или истечении ключа — model.ErrSessionNotFound.
// При повреждённых данных удаляет запись и возвращает model.ErrSessionNotFound (self-healing).
func (r *repository) Get(ctx context.Context, jti uuid.UUID) (*model.Session, error) {
	key := Key(jti)
	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, model.ErrSessionNotFound
	}

	var c resources.SessionCache
	if err := json.Unmarshal(data, &c); err != nil {
		_ = r.redis.Del(ctx, key)
		return nil, model.ErrSessionNotFound
	}

	session, err := converter.FromCache(c)
	if err != nil {
		_ = r.redis.Del(ctx, key)
		return nil, model.ErrSessionNotFound
	}

	return session, nil
}
