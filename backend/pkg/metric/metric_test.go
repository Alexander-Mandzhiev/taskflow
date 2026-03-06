package metric

import (
	"context"
	"testing"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

func TestNewWithLogger(t *testing.T) {
	t.Run("nil logger uses noop", func(t *testing.T) {
		m := NewWithLogger(nil)
		if m == nil {
			t.Fatal("NewWithLogger(nil) = nil")
		}
	})
	t.Run("with logger", func(t *testing.T) {
		m := NewWithLogger(&logger.NoopLogger{})
		if m == nil {
			t.Fatal("NewWithLogger(noop) = nil")
		}
	})
}

func TestMetrics_SetLogger(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	m.SetLogger(nil) // не паникует
	m.SetLogger(&logger.NoopLogger{})
}

func TestMetrics_Init_Disabled(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()
	err := m.Init(ctx, WithEnable(false))
	if err != nil {
		t.Errorf("Init(WithEnable(false)) = %v", err)
	}
	if m.GetMeterProvider() != nil {
		t.Error("GetMeterProvider() after Init(disabled) should be nil")
	}
}

func TestMetrics_GetMeterProvider_NotInitialized(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	if m.GetMeterProvider() != nil {
		t.Error("GetMeterProvider() without Init should be nil")
	}
}

func TestMetrics_Shutdown_Nil(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()
	err := m.Shutdown(ctx, time.Second)
	if err != nil {
		t.Errorf("Shutdown when not initialized = %v", err)
	}
}

func TestMetrics_IsInitialized(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	if m.IsInitialized() {
		t.Error("IsInitialized() before Init = true, want false")
	}
	_ = m.Init(context.Background(), WithEnable(false))
	if m.IsInitialized() {
		t.Error("IsInitialized() after Init(disabled) = true, want false")
	}
}

func TestMetrics_IsEnabled(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	if m.IsEnabled() {
		t.Error("IsEnabled() before Init = true, want false")
	}
	_ = m.Init(context.Background(), WithEnable(false))
	if m.IsEnabled() {
		t.Error("IsEnabled() after Init(WithEnable(false)) = true, want false")
	}
}

func TestMetrics_GetConfig(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	if m.GetConfig() != nil {
		t.Error("GetConfig() before Init should be nil")
	}
	_ = m.Init(context.Background(), WithEnable(false))
	if m.GetConfig() == nil {
		t.Error("GetConfig() after Init should be non-nil")
	}
}

func TestMetrics_HealthCheck_NotInitialized(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	ctx := context.Background()
	err := m.HealthCheck(ctx)
	if err == nil {
		t.Fatal("HealthCheck() when not initialized expected error")
	}
	if err.Error() != "metrics not initialized" {
		t.Errorf("HealthCheck() error = %q, want \"metrics not initialized\"", err.Error())
	}
}

func TestMetrics_getMetricName(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	_ = m.Init(context.Background(), WithEnable(false), WithNamespace("ns"), WithAppName("app"))
	got := m.getMetricName("http_requests_total")
	want := "ns_app_http_requests_total"
	if got != want {
		t.Errorf("getMetricName = %q, want %q", got, want)
	}
}

func TestMetrics_getMetricName_NilConfig(t *testing.T) {
	m := NewWithLogger(&logger.NoopLogger{})
	// без Init config == nil
	got := m.getMetricName("foo")
	if got != "foo" {
		t.Errorf("getMetricName with nil config = %q, want \"foo\"", got)
	}
}
