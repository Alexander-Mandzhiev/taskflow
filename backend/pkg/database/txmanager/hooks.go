package txmanager

import (
	"context"
	"sync"
)

// hookRegistryKey ключ для хранения HookRegistry в context.
type hookRegistryKey struct{}

// HookRegistry управляет регистрацией post-commit hooks.
// Позволяет адаптерам регистрировать hooks для инвалидации кэша и т.п.
// Дедуплицирует по ключу, сохраняя порядок первой регистрации.
type HookRegistry struct {
	mu    sync.RWMutex
	hooks map[string]PostCommitHook
	order []string
}

// NewHookRegistry создаёт новый registry для hooks.
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[string]PostCommitHook),
	}
}

// Register регистрирует hook для выполнения после успешного коммита.
// Если hook с таким ключом уже зарегистрирован, он будет перезаписан (порядок сохраняется).
func (r *HookRegistry) Register(key string, hook PostCommitHook) {
	if hook == nil || key == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.hooks[key]; !exists {
		r.order = append(r.order, key)
	}
	r.hooks[key] = hook
}

// GetHooks возвращает все зарегистрированные hooks в порядке регистрации.
func (r *HookRegistry) GetHooks() []PostCommitHook {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hooks := make([]PostCommitHook, 0, len(r.order))
	for _, key := range r.order {
		hooks = append(hooks, r.hooks[key])
	}
	return hooks
}

// WithHookRegistry добавляет HookRegistry в context.
func WithHookRegistry(ctx context.Context, registry *HookRegistry) context.Context {
	return context.WithValue(ctx, hookRegistryKey{}, registry)
}

// GetHookRegistry извлекает HookRegistry из context. Возвращает nil, если не найден.
func GetHookRegistry(ctx context.Context) *HookRegistry {
	if registry, ok := ctx.Value(hookRegistryKey{}).(*HookRegistry); ok {
		return registry
	}
	return nil
}
