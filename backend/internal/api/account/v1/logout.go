package account_v1

import (
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// Logout обрабатывает выход: читает refresh из cookie, передаёт в сервис; сервис валидирует и удаляет сессию по jti.
// Cookie очищаются в любом случае; при ошибке сервиса (сессия не найдена, невалидный токен) — 401, иначе 200.
func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := ""
	if c, err := r.Cookie(api.refreshTokenCookieName); err == nil && c != nil && c.Value != "" {
		refreshToken = c.Value
	}

	if err := api.accountService.Logout(r.Context(), refreshToken); err != nil {
		if isSessionInvalidOrExpiredError(err) {
			pkghttp.DeleteCookie(w, api.accessTokenCookieName, api.isSecure, api.cookieDomain)
			pkghttp.DeleteCookie(w, api.refreshTokenCookieName, api.isSecure, api.cookieDomain)
		}
		mapError(w, r, err)
		return
	}

	api.clearSessionAndRespond(w, r)
}

// clearSessionAndRespond удаляет cookie и пишет успешный ответ о завершении сессии.
func (api *API) clearSessionAndRespond(w http.ResponseWriter, r *http.Request) {
	pkghttp.DeleteCookie(w, api.accessTokenCookieName, api.isSecure, api.cookieDomain)
	pkghttp.DeleteCookie(w, api.refreshTokenCookieName, api.isSecure, api.cookieDomain)
	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.LogoutResponse{
		Success: true,
		Message: "Сессия завершена",
	})
}
