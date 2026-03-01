package account_v1

import (
	"encoding/json"
	"net/http"

	"mkk/internal/api/account/v1/dto"
	pkghttp "mkk/pkg/http"
)

// Login обрабатывает вход по email и паролю.
// При успехе создаёт сессию и устанавливает session_id в httpOnly cookie.
func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	sessionID, err := api.accountService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.SetCookie(w, "session_id", sessionID.String(), int(api.sessionTTL.Seconds()), api.isSecure, api.cookieDomain)

	pkghttp.WriteJSON(w, http.StatusOK, dto.LoginResponse{
		Success: true,
		Message: "Сессия создана, session_id установлен в httpOnly cookie",
	})
}
