package account_v1

import (
	"net/http"

	"github.com/google/uuid"

	"mkk/internal/api/account/v1/dto"
	pkghttp "mkk/pkg/http"
	"mkk/pkg/metadata"
)

// Logout обрабатывает выход: удаляет сессию и cookie session_id.
// Если сессии в контексте нет, cookie всё равно сбрасывается и возвращается 200.
func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := metadata.SessionID(r.Context())
	if err != nil || sessionID == uuid.Nil {
		api.clearSessionAndRespond(w)
		return
	}

	if err := api.accountService.Logout(r.Context(), sessionID); err != nil {
		mapError(w, r, err)
		return
	}

	api.clearSessionAndRespond(w)
}

// clearSessionAndRespond удаляет cookie и пишет успешный ответ о завершении сессии.
func (api *API) clearSessionAndRespond(w http.ResponseWriter) {
	pkghttp.DeleteCookie(w, "session_id", api.isSecure, api.cookieDomain)
	pkghttp.WriteJSON(w, http.StatusOK, dto.LogoutResponse{
		Success: true,
		Message: "Сессия завершена",
	})
}
