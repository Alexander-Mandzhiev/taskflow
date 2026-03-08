package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

var _ grpc.Notification = (*countingMockNotification)(nil)

// countingMockNotification считает вызовы NotifyInvitation и возвращает заданную ошибку (или nil после N вызовов).
type countingMockNotification struct {
	mu              sync.Mutex
	callCount       int
	err             error
	returnNilAfter  int // после стольких вызовов возвращать nil (0 = всегда err)
}

func (m *countingMockNotification) NotifyInvitation(ctx context.Context, inv *model.TeamInvitation, teamName, inviterName string) error {
	m.mu.Lock()
	m.callCount++
	n := m.callCount
	err := m.err
	if m.returnNilAfter > 0 && n > m.returnNilAfter {
		err = nil
	}
	m.mu.Unlock()
	return err
}

func (m *countingMockNotification) getCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestNotificationWithCircuitBreaker_opens_after_consecutive_failures(t *testing.T) {
	inner := &countingMockNotification{err: errors.New("service unavailable")}
	st := DefaultNotificationCBSettings()
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 3 // для теста достаточно 3
	}
	w := NewNotificationWithCircuitBreaker(inner, st)
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "a@b.c"}

	// 3 вызова с ошибкой — inner вызывается 3 раза
	for i := 0; i < 3; i++ {
		err := w.NotifyInvitation(ctx, inv, "Team", "User")
		require.Error(t, err)
	}
	assert.Equal(t, 3, inner.getCallCount(), "inner должен быть вызван 3 раза")

	// 4-й вызов — цепь открыта, inner не вызывается, возвращается ошибка
	err := w.NotifyInvitation(ctx, inv, "Team", "User")
	require.Error(t, err)
	assert.Equal(t, 3, inner.getCallCount(), "при Open inner не должен вызываться")
}

func TestNotificationWithCircuitBreaker_halfopen_allows_probe_after_timeout(t *testing.T) {
	inner := &countingMockNotification{
		err:            errors.New("service unavailable"),
		returnNilAfter: 3, // с 4-го вызова inner возвращает nil (успех пробного запроса в HalfOpen)
	}
	st := DefaultNotificationCBSettings()
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		return counts.ConsecutiveFailures >= 3
	}
	st.Timeout = 30 * time.Millisecond // короткий таймаут для теста
	w := NewNotificationWithCircuitBreaker(inner, st)
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "a@b.c"}

	// 3 ошибки -> переход в Open
	for i := 0; i < 3; i++ {
		_ = w.NotifyInvitation(ctx, inv, "Team", "User")
	}
	assert.Equal(t, 3, inner.getCallCount())

	// 4-й вызов — Open, inner не вызывается
	_ = w.NotifyInvitation(ctx, inv, "Team", "User")
	assert.Equal(t, 3, inner.getCallCount())

	// Ждём перехода в HalfOpen
	time.Sleep(50 * time.Millisecond)

	// Пробный запрос: inner вызывается, возвращаем nil -> цепь закрывается
	err := w.NotifyInvitation(ctx, inv, "Team", "User")
	require.NoError(t, err)
	assert.Equal(t, 4, inner.getCallCount(), "один пробный вызов в HalfOpen")

	// Следующий вызов снова идёт в inner (цепь закрыта)
	err = w.NotifyInvitation(ctx, inv, "Team", "User")
	require.NoError(t, err)
	assert.Equal(t, 5, inner.getCallCount())
}

// TestNotifyInvitation_success — при успешном ответе inner ошибки нет, цепь остаётся закрытой.
func TestNotifyInvitation_success(t *testing.T) {
	inner := &countingMockNotification{err: nil}
	w := NewNotificationWithCircuitBreaker(inner, gobreaker.Settings{})
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "a@b.c"}

	err := w.NotifyInvitation(ctx, inv, "Team", "User")
	require.NoError(t, err)
	assert.Equal(t, 1, inner.getCallCount())

	err = w.NotifyInvitation(ctx, inv, "Team", "User")
	require.NoError(t, err)
	assert.Equal(t, 2, inner.getCallCount())
}

// TestNewNotificationWithCircuitBreaker_empty_settings — пустой st: используются только дефолты из DefaultNotificationCBSettings.
func TestNewNotificationWithCircuitBreaker_empty_settings(t *testing.T) {
	inner := &countingMockNotification{err: errors.New("fail")}
	w := NewNotificationWithCircuitBreaker(inner, gobreaker.Settings{})
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "x@y.z"}

	// Дефолтный порог — 5 подряд ошибок
	for i := 0; i < 5; i++ {
		_ = w.NotifyInvitation(ctx, inv, "T", "U")
	}
	assert.Equal(t, 5, inner.getCallCount())
	// 6-й вызов — Open, inner не вызывается
	err := w.NotifyInvitation(ctx, inv, "T", "U")
	require.Error(t, err)
	assert.Equal(t, 5, inner.getCallCount())
}

// TestNewNotificationWithCircuitBreaker_merge_custom_settings — кастомные поля перезаписывают дефолты.
func TestNewNotificationWithCircuitBreaker_merge_custom_settings(t *testing.T) {
	inner := &countingMockNotification{err: errors.New("fail")}
	stateChanges := 0
	st := gobreaker.Settings{
		Name:    "custom-name",
		Timeout: 10 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			stateChanges++
		},
	}
	w := NewNotificationWithCircuitBreaker(inner, st)
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "a@b.c"}

	// 2 ошибки — кастомный порог
	for i := 0; i < 2; i++ {
		_ = w.NotifyInvitation(ctx, inv, "T", "U")
	}
	assert.Equal(t, 2, inner.getCallCount())
	// 3-й — Open
	_ = w.NotifyInvitation(ctx, inv, "T", "U")
	assert.Equal(t, 2, inner.getCallCount())
	assert.GreaterOrEqual(t, stateChanges, 1, "OnStateChange должен вызываться при переходе в Open")
}

// TestNewNotificationWithCircuitBreaker_merge_optional_fields — покрытие веток MaxRequests, Interval, BucketPeriod, IsSuccessful, IsExcluded.
func TestNewNotificationWithCircuitBreaker_merge_optional_fields(t *testing.T) {
	inner := &countingMockNotification{err: errors.New("fail")}
	st := gobreaker.Settings{
		MaxRequests:  2,
		Interval:     time.Second,
		BucketPeriod: 100 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 1
		},
		IsSuccessful: func(err error) bool { return err == nil },
		IsExcluded:   func(err error) bool { return false },
	}
	w := NewNotificationWithCircuitBreaker(inner, st)
	ctx := context.Background()
	inv := &model.TeamInvitation{ID: uuid.New(), TeamID: uuid.New(), Email: "a@b.c"}

	// 1 ошибка -> Open (кастомный ReadyToTrip)
	_ = w.NotifyInvitation(ctx, inv, "T", "U")
	assert.Equal(t, 1, inner.getCallCount())
	_ = w.NotifyInvitation(ctx, inv, "T", "U")
	assert.Equal(t, 1, inner.getCallCount())
}
