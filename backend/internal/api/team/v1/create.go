package team_v1

import (
	"encoding/json"
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/validation"
)

// Create создаёт команду. POST /api/v1/teams. Body: CreateTeamRequest. Пользователь из JWT становится owner.
func (api *API) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mapError(w, r, err)
		return
	}

	if err := validation.Validator.Struct(req); err != nil {
		mapError(w, r, err)
		return
	}

	input := converter.CreateTeamRequestToDomain(req)
	team, err := api.teamService.Create(r.Context(), input, userID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusCreated, converter.TeamToResponse(team))
}
