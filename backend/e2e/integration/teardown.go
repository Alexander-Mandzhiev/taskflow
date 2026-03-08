//go:build integration

package integration

import (
	"context"
	"log"

	"github.com/Alexander-Mandzhiev/taskflow/backend/e2e/testcontainers"
)

// teardownTestEnvironment останавливает контейнеры и логирует результат.
func teardownTestEnvironment(ctx context.Context, env *testcontainers.TestEnvironment) {
	log.Println("🧹 Очистка тестового окружения...")
	env.Cleanup(ctx)
	log.Println("✅ Тестовое окружение успешно очищено")
}
