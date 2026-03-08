package circuitbreaker

import (
	"context"
	"time"

	"github.com/sony/gobreaker/v2"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
)

// Настройки по умолчанию для circuit breaker сервиса уведомлений (email при invite).
const (
	DefaultNotificationCBName        = "notification"
	DefaultNotificationCBTimeout    = 60 * time.Second // время в состоянии Open перед переходом в HalfOpen
	DefaultNotificationCBMaxFailure = 5                // после стольких подряд ошибок — переход в Open
)

// NotificationWithCircuitBreaker оборачивает grpc.Notification в circuit breaker:
// при серии ошибок вызовы временно не выполняются (Open), затем один пробный (HalfOpen).
var _ grpc.Notification = (*NotificationWithCircuitBreaker)(nil)

type NotificationWithCircuitBreaker struct {
	inner grpc.Notification
	cb    *gobreaker.CircuitBreaker[struct{}]
}

// NewNotificationWithCircuitBreaker возвращает реализацию grpc.Notification, которая проксирует
// вызовы во inner через gobreaker. За основу берутся DefaultNotificationCBSettings(); переданные
// в st непустые поля перезаписывают значения по умолчанию.
func NewNotificationWithCircuitBreaker(inner grpc.Notification, st gobreaker.Settings) *NotificationWithCircuitBreaker {
	base := DefaultNotificationCBSettings()
	if st.Name != "" {
		base.Name = st.Name
	}
	if st.Timeout != 0 {
		base.Timeout = st.Timeout
	}
	if st.MaxRequests != 0 {
		base.MaxRequests = st.MaxRequests
	}
	if st.Interval != 0 {
		base.Interval = st.Interval
	}
	if st.BucketPeriod > 0 {
		base.BucketPeriod = st.BucketPeriod
	}
	if st.ReadyToTrip != nil {
		base.ReadyToTrip = st.ReadyToTrip
	}
	if st.OnStateChange != nil {
		base.OnStateChange = st.OnStateChange
	}
	if st.IsSuccessful != nil {
		base.IsSuccessful = st.IsSuccessful
	}
	if st.IsExcluded != nil {
		base.IsExcluded = st.IsExcluded
	}
	cb := gobreaker.NewCircuitBreaker[struct{}](base)
	return &NotificationWithCircuitBreaker{inner: inner, cb: cb}
}

// DefaultNotificationCBSettings возвращает настройки по умолчанию для notification circuit breaker.
func DefaultNotificationCBSettings() gobreaker.Settings {
	return gobreaker.Settings{
		Name:    DefaultNotificationCBName,
		Timeout: DefaultNotificationCBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= DefaultNotificationCBMaxFailure
		},
	}
}

// NotifyInvitation выполняет вызов inner.NotifyInvitation через circuit breaker.
// Если breaker в состоянии Open, вызов не выполняется и возвращается ошибка.
func (w *NotificationWithCircuitBreaker) NotifyInvitation(ctx context.Context, inv *model.TeamInvitation, teamName, inviterName string) error {
	_, err := w.cb.Execute(func() (struct{}, error) {
		return struct{}{}, w.inner.NotifyInvitation(ctx, inv, teamName, inviterName)
	})
	return err
}
