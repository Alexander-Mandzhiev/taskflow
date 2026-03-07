package team_v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Invite приглашает пользователя в команду. POST /api/v1/teams/{id}/invite. Только owner/admin. inviterUserID из JWT, проверка прав в сервисе.
func (api *API) Invite(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	teamID, err := util.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		mapError(w, r, model.ErrInvalidID)
		return
	}
	var req dto.InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	invitedMember, err := api.teamService.InviteByEmail(r.Context(), teamID, userID, req.Email, req.Role)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusCreated, dto.InviteResponse{
		Success: true,
		Message: "Пользователь приглашён в команду",
		Member:  converter.MemberToResponse(invitedMember),
	})
}
