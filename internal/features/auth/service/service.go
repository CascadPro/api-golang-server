package auth_service

import (
	"context"
	"net"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/client"
	core_postgres_token "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/token"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
	settings_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/settings/repository/postgres"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
)

type Service struct {
	userPostgresRepo     users_postgres_repository.RepositoryMethods
	settingsPostgresRepo settings_postgres_repository.RepositoryMethods
	tokenPostgresRepo    core_postgres_token.RepositoryMethods
	ipinfoRepo           core_ipinfo_client.RepositoryMethods
	sessionsRedisRepo    sessions_redis_repository.RepositoryMethods
}

type ServiceMethods interface {
	CreateRegisterToken(context.Context, domain.User) (domain.Token, error)
	Register(context.Context, domain.UserPatch, string) error
	Login(context.Context, domain.User, net.IP, *core_http_request.UserAgent) (string, string, error)
	GetNewTokens(context.Context, string) (string, error)
	Logout(context.Context) error
}

func NewService(
	userPostgresRepo users_postgres_repository.RepositoryMethods,
	settingsPostgresRepo settings_postgres_repository.RepositoryMethods,
	tokenPostgresRepo core_postgres_token.RepositoryMethods,
	ipinfoRepo core_ipinfo_client.RepositoryMethods,
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods,
) *Service {
	return &Service{
		userPostgresRepo:     userPostgresRepo,
		settingsPostgresRepo: settingsPostgresRepo,
		tokenPostgresRepo:    tokenPostgresRepo,
		ipinfoRepo:           ipinfoRepo,
		sessionsRedisRepo:    sessionsRedisRepo,
	}
}
