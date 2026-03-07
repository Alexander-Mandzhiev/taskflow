package team_v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
)

// GetByID возвращает команду с участниками по id. GET /api/v1/teams/{id}.
func (api *API) GetByID(w http.ResponseWriter, r *http.Request) {
	teamID, err := util.ParseUUID(chi.URLParam(r, "id"))
	if err != nil {
		mapError(w, r, model.ErrInvalidID)
		return
	}
	team, err := api.teamService.GetByID(r.Context(), teamID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, converter.TeamWithMembersToResponse(team))
}
