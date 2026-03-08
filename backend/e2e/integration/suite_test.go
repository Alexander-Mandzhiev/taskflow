//go:build integration

// Пакет integration — e2e-тесты с testcontainers (MySQL, Redis, backend) и Ginkgo.
// Запуск из backend/: go test -tags=integration -v -timeout 10m ./e2e/integration/
package integration

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Alexander-Mandzhiev/taskflow/backend/e2e/testcontainers"
)

var env *testcontainers.TestEnvironment

var (
	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Taskflow Backend Integration Test Suite")
}

var _ = BeforeSuite(func() {
	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	logMsg("Запуск тестового окружения...")
	var err error
	env, err = testcontainers.Setup(suiteCtx)
	if err != nil {
		panic(fmt.Sprintf("integration setup: %v", err))
	}
	logMsg("setup OK, BackendURL=%s", env.BackendURL)
})

var _ = AfterSuite(func() {
	logMsg("Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	if suiteCancel != nil {
		suiteCancel()
	}
})

func logMsg(format string, args ...interface{}) {
	GinkgoWriter.Println(fmt.Sprintf(format, args...))
}
