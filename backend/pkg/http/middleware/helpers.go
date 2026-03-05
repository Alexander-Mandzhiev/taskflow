package middleware

import (
	"net/http"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// WriteError записывает ошибку в формате JSON через pkghttp.WriteJSON.
func WriteError(w http.ResponseWriter, code int, message string) {
	pkghttp.WriteJSON(w, code, pkghttp.ErrorBody{Code: code, Message: message})
}
