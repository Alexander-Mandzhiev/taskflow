package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/model"
	"github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/task/repository/resources"
)

// ToTeamTaskStats преобразует строку отчёта в model.TeamTaskStats.
func ToTeamTaskStats(r resources.TeamTaskStatsRow) (model.TeamTaskStats, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamTaskStats{}, err
	}
	return model.TeamTaskStats{
		TeamID:         teamID,
		TeamName:       r.TeamName,
		MemberCount:    r.MemberCount,
		DoneTasksCount: r.DoneTasksCount,
	}, nil
}

// ToTeamTopCreator преобразует строку отчёта в model.TeamTopCreator.
func ToTeamTopCreator(r resources.TeamTopCreatorRow) (model.TeamTopCreator, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return model.TeamTopCreator{}, err
	}
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.TeamTopCreator{}, err
	}
	return model.TeamTopCreator{
		TeamID:       teamID,
		UserID:       userID,
		Rank:         r.Rank,
		CreatedCount: r.CreatedCount,
	}, nil
}
