package task

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository"
	taskRepo "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/repository/task"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/database/txmanager"
)

var _ repository.TaskRepository = (*Adapter)(nil)

// Adapter — адаптер репозитория задач (таблица tasks).
// listCache — опционально; при nil кеш списка и post-commit инвалидация не используются.
type Adapter struct {
	taskReader taskRepo.TaskReaderRepository
	taskWriter taskRepo.TaskWriterRepository
	listCache  repository.TaskListCacheRepository
}

// NewAdapter создаёт адаптер задач. listCache может быть nil — тогда кеширование отключено.
func NewAdapter(
	taskReader taskRepo.TaskReaderRepository,
	taskWriter taskRepo.TaskWriterRepository,
	listCache repository.TaskListCacheRepository,
) *Adapter {
	return &Adapter{
		taskReader: taskReader,
		taskWriter: taskWriter,
		listCache:  listCache,
	}
}

// registerInvalidateHook регистрирует post-commit хук инвалидации кеша списка задач по teamID.
// Вызывается из Create/Update/SoftDelete/Restore; при отсутствии listCache или HookRegistry в ctx — no-op.
func (r *Adapter) registerInvalidateHook(ctx context.Context, teamID uuid.UUID) {
	if r.listCache == nil {
		return
	}
	reg := txmanager.GetHookRegistry(ctx)
	if reg == nil {
		return
	}
	reg.Register(fmt.Sprintf("task:list:invalidate:team:%s", teamID.String()), func(ctx context.Context) error {
		return r.listCache.InvalidateByTeam(ctx, teamID)
	})
}
