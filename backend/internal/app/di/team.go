package di

import (
	"context"
	"fmt"

	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	teamNotificationV1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc/notification/v1"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	teamRepoAdapter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter"
	teamRepoInvitationReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation/reader"
	teamRepoInvitationWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/invitation/writer"
	teamRepoMemberReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member/reader"
	teamRepoMemberWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/member/writer"
	teamRepoTeamReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team/reader"
	teamRepoTeamWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/team/writer"
	teamServiceDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service"
	teamServiceImpl "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/service/service"
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
	repo, err := d.TeamAdapter(ctx)
	if err != nil {
		return nil, fmt.Errorf("team adapter: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx manager: %w", err)
	}
	userRepo, err := d.UserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user repository: %w", err)
	}
	notifier := teamNotificationV1.NewClient()
	d.teamService = teamServiceImpl.NewTeamService(repo, txMgr, userRepo, notifier)
	return d.teamService, nil
}

// TeamAdapter возвращает адаптер доступа к данным команд (team + member + invitation reader/writer).
func (d *Container) TeamAdapter(ctx context.Context) (teamRepoDef.TeamAdapter, error) {
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
	invitationReader := teamRepoInvitationReader.NewRepository(db)
	invitationWriter := teamRepoInvitationWriter.NewRepository(db)
	d.teamRepo = teamRepoAdapter.NewAdapter(teamReader, teamWriter, memberReader, memberWriter, invitationReader, invitationWriter)
	return d.teamRepo, nil
}
