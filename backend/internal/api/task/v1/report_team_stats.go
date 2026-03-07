package task_v1

import (
	"net/http"
	"time"

	taskconverter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/converter"
	pkghttp "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/http"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/metadata"
)

const defaultReportSinceDays = 7

// ReportTeamStats возвращает для каждой команды: название, кол-во участников, кол-во задач done за период.
// GET /api/v1/reports/team-stats?since_days=7 (по умолчанию 7 дней).
func (api *API) ReportTeamStats(w http.ResponseWriter, r *http.Request) {
	userID, err := metadata.UserID(r.Context())
	if err != nil {
		mapError(w, r, err)
		return
	}

	sinceDays := defaultReportSinceDays
	if d := r.URL.Query().Get("since_days"); d != "" {
		n, err := parseIntPositive(d)
		if err != nil {
			mapError(w, r, err)
			return
		}
		sinceDays = n
	}
	since := time.Now().AddDate(0, 0, -sinceDays)

	items, err := api.reportService.TeamTaskStats(r.Context(), userID, since)
	if err != nil {
		mapError(w, r, err)
		return
	}

	pkghttp.WriteJSON(r.Context(), w, http.StatusOK, taskconverter.TeamTaskStatsListToResponse(items))
}
