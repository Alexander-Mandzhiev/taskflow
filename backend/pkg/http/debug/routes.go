// Package debug регистрирует маршруты для профилирования (pprof) под /debug/pprof/.
package debug

import (
	"net/http"
	"net/http/pprof"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes вешает стандартные pprof-обработчики на /debug/pprof/...
// CPU profile только один одновременно — повторный запрос получает 503.
func RegisterRoutes(r chi.Router) {
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", profileOnce(pprof.Profile))
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
}

var profileInUse int32

// profileOnce разрешает только один одновременный вызов CPU profile; второй получает 503.
func profileOnce(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !atomic.CompareAndSwapInt32(&profileInUse, 0, 1) {
			w.Header().Set("Retry-After", "30")
			http.Error(w, "CPU profile already in progress (only one at a time). Retry in ~30s.", http.StatusServiceUnavailable)
			return
		}
		defer atomic.StoreInt32(&profileInUse, 0)
		handler(w, r)
	}
}
