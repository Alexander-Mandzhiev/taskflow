package account_v1

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/account/v1/dto"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
)

// Register обрабатывает регистрацию пользователя.
// POST body: RegisterRequest (email, password, name).
// При успехе возвращает 201 и RegisterResponse (success + message).
func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := api.accountService.Register(r.Context(), req.Email, req.Password, req.Name); err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(w, http.StatusCreated, dto.RegisterResponse{
		Success: true,
		Message: "Пользователь успешно зарегистрирован",
	})
}
