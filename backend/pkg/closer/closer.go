package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/logger"
)

// shutdownTimeout — максимальное время ожидания graceful shutdown.
// Должен быть >= максимального shutdown_timeout у внешних зависимостей (OTel и т.п.).
const shutdownTimeout = 60 * time.Second

// ErrNotSet возвращается при попытке использовать Closer, который не был установлен (nil).
var ErrNotSet = errors.New("closer not set")

// Closer управляет процессом graceful shutdown приложения.
type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(context.Context) error
	logger Logger
}

// New создаёт новый Closer с NoopLogger. Если переданы сигналы, запускает их обработку.
// После logger.Init() вызвать SetLogger(logger.Logger()).
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(&logger.NoopLogger{}, signals...)
}

// NewWithLogger создаёт новый Closer с указанным логгером. Если переданы сигналы, при получении вызывается CloseAll.
// Если l == nil (например logger ещё не инициализирован), используется NoopLogger.
func NewWithLogger(l Logger, signals ...os.Signal) *Closer {
	if l == nil {
		l = &logger.NoopLogger{}
	}
	c := &Closer{
		done:   make(chan struct{}),
		logger: l,
	}
	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}
	return c
}

// SetLogger устанавливает логгер для Closer. Если l == nil, не меняет текущий логгер.
func (c *Closer) SetLogger(l Logger) {
	if l != nil {
		c.logger = l
	}
}

// handleSignals обрабатывает системные сигналы и вызывает CloseAll с контекстом shutdownTimeout.
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "Received system signal, starting graceful shutdown")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()
		_ = c.CloseAll(shutdownCtx)
	case <-c.done:
		// CloseAll уже вызван вручную
	}
}

// AddNamed добавляет функцию закрытия с именем для логирования.
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("Closing %s", name))
		err := f(ctx)
		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("Failed to close %s: %v (%s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("Closed %s in %s", name, duration))
		}
		return err
	})
}

// Add добавляет одну или несколько функций закрытия.
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll вызывает все зарегистрированные функции закрытия последовательно в обратном порядке добавления.
// Возвращает первую возникшую ошибку. Вызов идемпотентен — выполняется один раз.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "No close functions registered")
			return
		}

		c.logger.Info(ctx, "Starting graceful shutdown")

		for i := len(funcs) - 1; i >= 0; i-- {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "Context cancelled during shutdown", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			default:
			}

			fn := funcs[i]
			func() {
				defer func() {
					if r := recover(); r != nil {
						err := errors.New("panic recovered in closer")
						c.logger.Error(ctx, "Panic in close function", zap.Any("error", r))
						if result == nil {
							result = err
						}
					}
				}()
				if err := fn(ctx); err != nil {
					c.logger.Error(ctx, "Close function failed", zap.Error(err))
					if result == nil {
						result = err
					}
				}
			}()
		}

		c.logger.Info(ctx, "All resources closed")
	})
	return result
}
