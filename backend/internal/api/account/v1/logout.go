package account_v1

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// Logout обрабатывает выход: удаляет сессию и cookie session_id.
// Если сессии в контексте нет, cookie всё равно сбрасывается и возвращается 200.
// При ошибке «сессия не найдена/истекла» cookie также удаляется, чтобы клиент не оставался с мёртвой cookie; ответ — 401.
func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := metadata.SessionID(r.Context())
	if err != nil || sessionID == uuid.Nil {
		api.clearSessionAndRespond(w, r)
		return
	}

	if err := api.accountService.Logout(r.Context(), sessionID); err != nil {
		// При ошибках «сессия не найдена/истекла» удаляем cookie, чтобы клиент не оставался с мёртвой cookie.
		if isSessionInvalidOrExpiredError(err) {
			pkghttp.DeleteCookie(w, string(ctxkey.SessionID), api.isSecure, api.cookieDomain)
		}
		mapError(w, r, err)
		return
	}

	api.clearSessionAndRespond(w, r)
}

// clearSessionAndRespond удаляет cookie и пишет успешный ответ о завершении сессии.
func (api *API) clearSessionAndRespond(w http.ResponseWriter, r *http.Request) {
	pkghttp.DeleteCookie(w, string(ctxkey.SessionID), api.isSecure, api.cookieDomain)
	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.LogoutResponse{
		Success: true,
		Message: "Сессия завершена",
	})
}
