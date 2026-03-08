package di

import (
	"context"
	"fmt"

	team_v1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/api/team/v1"
	teamClientGrpc "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc"
	teamClientCB "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc/circuitbreaker"
	teamNotificationV1 "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/client/grpc/notification/v1"
	teamRepoDef "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository"
	teamRepoAdapterInvitation "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/invitation"
	teamRepoAdapterMember "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/member"
	teamRepoAdapterTeam "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/adapter/team"
	teamRepoInvitationReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/invitation/reader"
	teamRepoInvitationWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/invitation/writer"
	teamRepoMemberReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/member/reader"
	teamRepoMemberWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/member/writer"
	teamRepoTeamReader "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/team/reader"
	teamRepoTeamWriter "github.com/Alexander-Mandzhiev/taskflow/backend/internal/module/workspace/team/repository/repository/team/writer"
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

// TeamService возвращает сервис команд.
func (d *Container) TeamService(ctx context.Context) (teamServiceDef.TeamService, error) {
	if d.teamService != nil {
		return d.teamService, nil
	}
	teamRepo, err := d.TeamRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("team repo: %w", err)
	}
	memberRepo, err := d.MemberRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("member repo: %w", err)
	}
	invitationRepo, err := d.InvitationRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("invitation repo: %w", err)
	}
	txMgr, err := d.UserTxManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("tx manager: %w", err)
	}
	userRepo, err := d.UserRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("user repository: %w", err)
	}
	notifier, err := d.teamNotifierOrInit(ctx)
	if err != nil {
		return nil, fmt.Errorf("team notifier: %w", err)
	}
	d.teamService = teamServiceImpl.NewTeamService(teamRepo, memberRepo, invitationRepo, txMgr, userRepo, notifier)
	return d.teamService, nil
}

// teamNotifierOrInit возвращает кешированный notifier для команды (с circuit breaker); при первом вызове создаёт и кеширует.
func (d *Container) teamNotifierOrInit(_ context.Context) (teamClientGrpc.Notification, error) {
	if d.teamNotifier != nil {
		return d.teamNotifier, nil
	}
	notificationClient := teamNotificationV1.NewClient()
	d.teamNotifier = teamClientCB.NewNotificationWithCircuitBreaker(notificationClient, teamClientCB.DefaultNotificationCBSettings())
	return d.teamNotifier, nil
}

// initTeamRepos создаёт team/member/invitation адаптеры при первом обращении к любому из репозиториев.
func (d *Container) initTeamRepos(ctx context.Context) error {
	if d.teamRepo != nil {
		return nil
	}
	db, err := d.SqlxDB(ctx)
	if err != nil {
		return fmt.Errorf("sqlx db: %w", err)
	}
	teamReader := teamRepoTeamReader.NewRepository(db)
	teamWriter := teamRepoTeamWriter.NewRepository(db)
	memberReader := teamRepoMemberReader.NewRepository(db)
	memberWriter := teamRepoMemberWriter.NewRepository(db)
	invitationReader := teamRepoInvitationReader.NewRepository(db)
	invitationWriter := teamRepoInvitationWriter.NewRepository(db)
	d.teamRepo = teamRepoAdapterTeam.NewAdapter(teamReader, teamWriter)
	d.memberRepo = teamRepoAdapterMember.NewAdapter(memberReader, memberWriter)
	d.invitationRepo = teamRepoAdapterInvitation.NewAdapter(invitationReader, invitationWriter)
	return nil
}

// TeamRepository возвращает репозиторий команд (таблица teams).
func (d *Container) TeamRepository(ctx context.Context) (teamRepoDef.TeamRepository, error) {
	if err := d.initTeamRepos(ctx); err != nil {
		return nil, err
	}
	return d.teamRepo, nil
}

// MemberRepository возвращает репозиторий участников команд (таблица team_members).
func (d *Container) MemberRepository(ctx context.Context) (teamRepoDef.MemberRepository, error) {
	if err := d.initTeamRepos(ctx); err != nil {
		return nil, err
	}
	return d.memberRepo, nil
}

// InvitationRepository возвращает репозиторий приглашений (таблица team_invitations).
func (d *Container) InvitationRepository(ctx context.Context) (teamRepoDef.InvitationRepository, error) {
	if err := d.initTeamRepos(ctx); err != nil {
		return nil, err
	}
	return d.invitationRepo, nil
}
