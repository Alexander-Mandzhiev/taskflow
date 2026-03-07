package di

import (
	"context"
	"fmt"

	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository"
	teamRepoAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/adapter"
	teamRepoMemberReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/member/reader"
	teamRepoMemberWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/member/writer"
	teamRepoTeamReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/team/reader"
	teamRepoTeamWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/repository/team/writer"
	teamServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/service"
	teamServiceImpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/team/service/service"
)

// TeamV1API возвращает HTTP API team v1 (create, list, get, invite).
func (d *Container) TeamV1API(ctx context.Context) (*team_v1.API, error) {
	if d.teamAPI != nil {
		return d.teamAPI, nil
	}
	teamSvc, err := d.TeamService(ctx)
	if err != nil {
		return nil, fmt.Errorf("team service: %w", err)
	}
	d.teamAPI = team_v1.NewAPI(teamSvc)
	return d.teamAPI, nil
}

// TeamService возвращает сервис команд. Явно прокидывается адаптер пользователей (userRepo) для invite по email.
func (d *Container) TeamService(ctx context.Context) (teamServiceDef.TeamService, error) {
	if d.teamService != nil {
		return d.teamService, nil
	}
	repo, err := d.TeamRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("team repository: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx manager: %w", err)
	}
	userRepo, err := d.UserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user repository: %w", err)
	}
	d.teamService = teamServiceImpl.NewTeamService(repo, txMgr, userRepo)
	return d.teamService, nil
}

// TeamRepository возвращает адаптер репозитория команд (team reader/writer + member reader/writer).
func (d *Container) TeamRepository(ctx context.Context) (teamRepoDef.TeamRepository, error) {
	if d.teamRepo != nil {
		return d.teamRepo, nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlx db: %w", err)
	}
	teamReader := teamRepoTeamReader.NewRepository(db)
	teamWriter := teamRepoTeamWriter.NewRepository(db)
	memberReader := teamRepoMemberReader.NewRepository(db)
	memberWriter := teamRepoMemberWriter.NewRepository(db)
	d.teamRepo = teamRepoAdapter.NewRepository(teamReader, teamWriter, memberReader, memberWriter)
	return d.teamRepo, nil
}
