package account_v1

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/identity/account/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Login обрабатывает вход по email и паролю.
// При успехе создаёт сессию (JWT), выставляет access_token и refresh_token в cookie, возвращает { success, message }.
func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	input := converter.LoginRequestToDomain(req, r.UserAgent(), pkghttp.ClientIP(r))
	accessToken, refreshToken, err := api.accountService.Login(r.Context(), input)
	if err != nil {
		mapError(w, r, err)
		return
	}

	// Access — в cookie без httpOnly (фронт может читать).
	pkghttp.SetCookie(w, api.accessTokenCookieName, accessToken, int(api.accessTTL.Seconds()), api.isSecure, api.cookieDomain, false)
	// Refresh — httpOnly, с защитой (только сервер при refresh).
	pkghttp.SetCookie(w, api.refreshTokenCookieName, refreshToken, int(api.refreshTTL.Seconds()), api.isSecure, api.cookieDomain, true)
	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.LoginResponse{Success: true, Message: "Успешный вход"})
}
