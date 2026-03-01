package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"mkk/internal/module/identity/user/model"
	"mkk/internal/module/identity/user/repository/converter"
	"mkk/internal/module/identity/user/repository/resources"
)

// Get возвращает пользователя из кеша по id. При промахе — (nil, nil).
// При повреждённых данных удаляет запись и возвращает (nil, nil) — self-healing cache miss.
func (r *repository) Get(ctx context.Context, id string) (*model.User, error) {
	key := Key(id)
	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get: %w", err)
	}
	if data == nil {
		return nil, nil
	}
	var c resources.UserCache
	if err := json.Unmarshal(data, &c); err != nil {
		_ = r.redis.Del(ctx, key)
		return nil, nil
	}
	user, err := converter.FromCache(c)
	if err != nil {
		_ = r.redis.Del(ctx, key)
		return nil, nil
	}
	return &user, nil
}
