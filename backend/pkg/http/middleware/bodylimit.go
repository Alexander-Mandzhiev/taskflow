package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// BodyLimitMiddleware ограничивает размер тела запроса до maxBytes.
// При превышении лимита или пустом теле возвращает JSON-ошибку, не доходя до хендлера.
func BodyLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IsBodyError проверяет, является ли ошибка проблемой чтения тела (EOF, превышение лимита)
// и пишет соответствующий HTTP-ответ. Возвращает true, если ошибка обработана.
func IsBodyError(ctx context.Context, w http.ResponseWriter, err error) bool {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		pkghttp.WriteJSON(ctx, w, http.StatusBadRequest, pkghttp.ErrorBody{
			Code:    http.StatusBadRequest,
			Message: "Тело запроса пустое или обрезано",
		})
		return true
	}

	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		pkghttp.WriteJSON(ctx, w, http.StatusRequestEntityTooLarge, pkghttp.ErrorBody{
			Code:    http.StatusRequestEntityTooLarge,
			Message: "Тело запроса слишком большое",
		})
		return true
	}

	return false
}
