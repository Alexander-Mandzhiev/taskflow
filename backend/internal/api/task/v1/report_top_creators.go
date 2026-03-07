package task_v1

import (
	"net/http"
	"time"

	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const (
	defaultTopCreatorsSinceDays = 30
	defaultTopCreatorsLimit     = 3
)

// ReportTopCreators возвращает топ-N пользователей по созданным задачам в каждой команде за период.
// GET /api/v1/reports/top-creators?since_days=30&limit=3.
func (api *API) ReportTopCreators(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	sinceDays := defaultTopCreatorsSinceDays
	if d := r.URL.Query().Get("since_days"); d != "" {
		n, err := parseIntPositive(d)
		if err != nil {
			mapError(w, r, err)
			return
		}
		sinceDays = n
	}
	since := time.Now().AddDate(0, 0, -sinceDays)

	limit := defaultTopCreatorsLimit
	if l := r.URL.Query().Get("limit"); l != "" {
		n, err := parseIntPositive(l)
		if err != nil {
			mapError(w, r, err)
			return
		}
		limit = n
	}
	if err := checkReportLimit(limit); err != nil {
		mapError(w, r, err)
		return
	}

	items, err := api.reportService.TopCreatorsByTeam(r.Context(), userID, since, limit)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, taskconverter.TopCreatorsListToResponse(items))
}

func checkReportLimit(limit int) error {
	if limit <= 0 {
		return model.ErrInvalidLimit
	}
	return nil
}
