package team_v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/util"
)

// GetByID возвращает команду с участниками по id. GET /api/v1/teams/{id}. Доступ только для участников команды.
func (api *API) GetByID(w http.ResponseWriter, r *http.Request) {
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
	team, err := api.teamService.GetByID(r.Context(), teamID, userID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, converter.TeamWithMembersToResponse(team))
}
