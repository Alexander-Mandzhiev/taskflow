package txmanager

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

const maxAttrLen = 200

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func formatPanic(r any) string {
	if r == nil {
		return ""
	}
	return truncateStr(fmt.Sprintf("%T: %v", r, r), maxAttrLen)
}

// executeHooks выполняет все post-commit hooks с защитой от паник и логированием.
func (m *Manager) executeHooks(ctx context.Context, hooks []PostCommitHook) {
	for i, hook := range hooks {
		if hook == nil {
			continue
		}

		select {
		case <-ctx.Done():
			logger.Warn(ctx, "Context cancelled, stopping hooks",
				zap.Int("remaining_hooks", len(hooks)-i),
			)
			return
		default:
		}

		if deadline, ok := ctx.Deadline(); ok {
			if remaining := time.Until(deadline); remaining <= 100*time.Millisecond {
				logger.Warn(ctx, "Skipping hook, insufficient time before deadline",
					zap.Int("hook_index", i),
					zap.Duration("remaining", remaining),
				)
				continue
			}
		}

		hookCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

		hookStart := time.Now()
		func() {
			defer cancel()
			defer func() {
				if r := recover(); r != nil {
					logger.Error(ctx, "Panic in hook",
						zap.Int("hook_index", i),
						zap.String("panic", formatPanic(r)),
					)
				}
			}()

			if err := hook(hookCtx); err != nil {
				hookDuration := time.Since(hookStart)
				logger.Warn(ctx, "Post-commit hook failed",
					zap.Int("hook_index", i),
					zap.Error(err),
					zap.Duration("duration", hookDuration),
				)
				if span := trace.SpanFromContext(ctx); span.IsRecording() {
					span.RecordError(err)
					span.SetAttributes(
						attribute.String(fmt.Sprintf("db.hook.%d.status", i), "error"),
						attribute.String(fmt.Sprintf("db.hook.%d.error", i), truncateStr(err.Error(), maxAttrLen)),
					)
				}
			} else {
				hookDuration := time.Since(hookStart)
				if span := trace.SpanFromContext(ctx); span.IsRecording() {
					span.SetAttributes(
						attribute.String(fmt.Sprintf("db.hook.%d.status", i), "success"),
						attribute.Float64(fmt.Sprintf("db.hook.%d.duration_us", i), float64(hookDuration.Microseconds())),
					)
				}
			}
		}()
	}
}
