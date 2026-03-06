package middleware

import (
	"context"
	"net/http"

	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// AuthResponse — ответ при невалидной аутентификации (сессия или JWT).
type AuthResponse struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message,omitempty"`
}

// writeAuthResponse пишет 401 и JSON { authenticated: false, message }.
func writeAuthResponse(ctx context.Context, w http.ResponseWriter, message string) {
	pkghttp.WriteJSON(ctx, w, http.StatusUnauthorized, AuthResponse{
		Authenticated: false,
		Message:       message,
	})
}
