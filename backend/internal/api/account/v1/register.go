package account_v1

import (
	"encoding/json"
	"net/http"

	"mkk/internal/api/account/v1/dto"
	pkghttp "mkk/pkg/http"
)

// Register обрабатывает регистрацию пользователя.
// POST body: RegisterRequest (email, password, name).
// При успехе возвращает 201 и RegisterResponse с user_id.
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

	userID, err := api.accountService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(w, http.StatusCreated, dto.RegisterResponse{
		UserID:  userID,
		Message: "Пользователь успешно зарегистрирован",
	})
}
