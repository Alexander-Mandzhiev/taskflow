package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/user/repository/resources"
)

// Get возвращает пользователя из кеша по id. При промахе или повреждённых данных — (model.User{}, false, nil).
// При ошибке Redis — (model.User{}, false, err). При успехе — (user, true, nil).
func (r *repository) Get(ctx context.Context, id string) (model.User, bool, error) {
	key := Key(id)
	data, err := r.redis.Get(ctx, key)
	if err != nil {
		return model.User{}, false, fmt.Errorf("cache get: %w", err)
	}
	if data == nil {
		return model.User{}, false, nil
	}
	var c resources.UserCache
	if err := json.Unmarshal(data, &c); err != nil {
		_ = r.redis.Del(ctx, key)
		return model.User{}, false, nil
	}
	user, err := converter.FromCache(c)
	if err != nil {
		_ = r.redis.Del(ctx, key)
		return model.User{}, false, nil
	}
	return user, true, nil
}
