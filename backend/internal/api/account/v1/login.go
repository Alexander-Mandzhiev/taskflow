package account_v1

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/ctxkey"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Login обрабатывает вход по email и паролю.
// При успехе создаёт сессию и устанавливает session_id в httpOnly cookie.
// userAgent и IP передаются в сервис для отображения в списке сессий.
// Лимит тела запроса задаётся BodyLimitMiddleware в routes.RegisterAPIs (публичная группа).
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

	sessionID, err := api.accountService.Login(r.Context(), req.Email, req.Password, r.UserAgent(), pkghttp.ClientIP(r))
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.SetCookie(w, string(ctxkey.SessionID), sessionID.String(), int(api.sessionTTL.Seconds()), api.isSecure, api.cookieDomain)
	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, dto.LoginResponse{
		Success: true,
		Message: "Сессия создана, session_id установлен в httpOnly cookie",
	})
}
