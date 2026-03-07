package team_v1

import (
	"net/http"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1/dto"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

// List возвращает список команд, где пользователь участник. GET /api/v1/teams.
func (api *API) List(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	teams, err := api.teamService.ListByUserID(r.Context(), userID)
	if err != nil {
		mapError(w, r, err)
		return
	}

	resp := make([]dto.TeamWithRoleResponse, 0, len(teams))
	for _, t := range teams {
		resp = append(resp, converter.TeamWithRoleToResponse(t))
	}
	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, resp)
}
