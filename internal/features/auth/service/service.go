package auth_service

import (
	"context"
	"net"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/infrastructure/ipinfo/client"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	"github.com/google/uuid"
)

type Service struct {
	userPostgresRepo     UsersRepository
	settingsPostgresRepo SettingsRepository
	tokenPostgresRepo    TokenRepository
	ipinfoRepo           IPInfoClient
	sessionsRedisRepo    SessionsRepository

	tokenIssuer        core_jwt_security.IssuerMethods
	outboxPostgresRepo OutboxRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)

	GetUser(ctx context.Context, user domain.User) (domain.User, error)

	PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error)
}

type SettingsRepository interface {
	CreateUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)

	GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token domain.Token) (domain.Token, error)

	Get(ctx context.Context, token domain.Token) (domain.Token, error)

	Delete(ctx context.Context, id uuid.UUID) error
}

type IPInfoClient interface {
	Lookup(ctx context.Context, ip net.IP) (core_ipinfo_client.InfoModel, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, userID uuid.UUID, session domain.Session) (string, error)

	Get(ctx context.Context, userID uuid.UUID, sessionID string) (domain.Session, error)

	Delete(ctx context.Context, userID uuid.UUID, sessionID string) error
}

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

func NewService(
	userPostgresRepo UsersRepository,
	settingsPostgresRepo SettingsRepository,
	tokenPostgresRepo TokenRepository,
	ipinfoRepo IPInfoClient,
	sessionsRedisRepo SessionsRepository,
	tokenIssuer core_jwt_security.IssuerMethods,
	outboxPostgresRepo OutboxRepository,
) *Service {
	return &Service{
		userPostgresRepo:     userPostgresRepo,
		settingsPostgresRepo: settingsPostgresRepo,
		tokenPostgresRepo:    tokenPostgresRepo,
		ipinfoRepo:           ipinfoRepo,
		sessionsRedisRepo:    sessionsRedisRepo,
		tokenIssuer:          tokenIssuer,
		outboxPostgresRepo:   outboxPostgresRepo,
	}
}
