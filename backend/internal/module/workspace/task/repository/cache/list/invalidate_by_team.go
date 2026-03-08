package list

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// InvalidateByTeam удаляет все закешированные страницы списка задач команды.
func (r *Repository) InvalidateByTeam(ctx context.Context, teamID uuid.UUID) error {
	prefix := PrefixForTeam(teamID)
	if err := r.redis.DelByPrefix(ctx, prefix); err != nil {
		return fmt.Errorf("task list cache invalidate: %w", err)
	}
	return nil
}
