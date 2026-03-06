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

func writeAuthResponse(ctx context.Context, w http.ResponseWriter, authenticated bool, message string) {
	code := http.StatusOK
	if !authenticated {
		code = http.StatusUnauthorized
	}
	pkghttp.WriteJSON(ctx, w, code, AuthResponse{
		Authenticated: authenticated,
		Message:       message,
	})
}
