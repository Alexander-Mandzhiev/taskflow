package http_server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewServer создаёт и настраивает *http.Server, используя переданные параметры.
// Сервер получает адрес, таймауты чтения/записи/заголовков и лимит размера заголовков.
// Обработчиком служит переданный роутер r.
func NewServer(
	r *chi.Mux,
	address string,
	readHeaderTimeout, readTimeout, writeTimeout, idleTimeout time.Duration,
	maxHeaderBytes int,
) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}
