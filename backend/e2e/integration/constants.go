//go:build integration

package integration

import "time"

const (
	// projectName — имя проекта для логов и идентификации контейнеров.
	projectName = "taskflow-backend"

	// Таймауты
	testsTimeout   = 5 * time.Minute
	requestTimeout = 30 * time.Second
	startupTimeout = 5 * time.Minute

	// Количество повторных созданий для полноценной проверки CRUD и отчётов.
	createTeamsCount    = 3
	createTasksCount    = 3
	createCommentsCount = 3
	updateHistoryCount  = 3
	reportsCallCount    = 2

	// Пауза перед каждым запросом к API, чтобы не упираться в rate limit (7 req/s на register/login и т.д.).
	requestDelay = 150 * time.Millisecond
)
