package app

import (
	"fmt"
	"time"

	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/infrastructure/ipinfo/client"
	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	core_postgres_token "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/token"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	auth_service "github.com/CascadePro/api-golang-server/internal/features/auth/service"
	auth_transport_http "github.com/CascadePro/api-golang-server/internal/features/auth/transport/http"
	client_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/client/repository/postgres"
	client_service "github.com/CascadePro/api-golang-server/internal/features/client/service"
	client_transport_http "github.com/CascadePro/api-golang-server/internal/features/client/transport/http"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	media_transport_http "github.com/CascadePro/api-golang-server/internal/features/media/transport/http"
	requests_mongo_repository "github.com/CascadePro/api-golang-server/internal/features/requests/repository/mongo"
	requests_service "github.com/CascadePro/api-golang-server/internal/features/requests/service"
	requests_transport_http "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http"
	root_transport_http "github.com/CascadePro/api-golang-server/internal/features/root/transport/http"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
	session_service "github.com/CascadePro/api-golang-server/internal/features/sessions/service"
	sessions_transport_http "github.com/CascadePro/api-golang-server/internal/features/sessions/transport/http"
	settings_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/settings/repository/postgres"
	settings_service "github.com/CascadePro/api-golang-server/internal/features/settings/service"
	settings_transport_http "github.com/CascadePro/api-golang-server/internal/features/settings/transport/http"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	users_service "github.com/CascadePro/api-golang-server/internal/features/users/service"
	users_transport_http "github.com/CascadePro/api-golang-server/internal/features/users/transport/http"
	"go.uber.org/zap"
)

type Features struct {
	Root     *core_http_server.Router
	Media    *core_http_server.Router
	Users    *core_http_server.Router
	Settings *core_http_server.Router
	Clients  *core_http_server.Router
	Requests *core_http_server.Router
	Sessions *core_http_server.Router
	Auth     *core_http_server.Router
}

func (a *App) initFeatures() error {
	tokenPostgresRepository := core_postgres_token.NewRepository(a.infrastructure.Postgres)
	outboxPostgresRepository := core_postgres_outbox.NewRepository(a.infrastructure.Postgres)
	ipInfoRepository := core_ipinfo_client.NewRepository(a.infrastructure.IPInfo)

	logFeatureInit := func(name, path string) {
		a.logger.Debug(
			"initializing feature",
			zap.String("feature", name),
			zap.String("path", fmt.Sprintf("%s/*", path)),
		)
	}

	// Root routes
	logFeatureInit("root", "")

	rootRateLimit := core_http_middleware.NewRateLimitConfig(100, time.Minute*5)
	rootRouter := core_http_server.NewRouter("", rootRateLimit.Middleware())

	rootHttpHandler := root_transport_http.NewHttpHandler()
	rootRouter.RegisterRoutes(rootHttpHandler.Routes()...)

	// Media routes
	logFeatureInit("media", "/media")
	mediaRateLimit := core_http_middleware.NewRateLimitConfig(25, time.Minute*10)
	mediaRouter := core_http_server.NewRouter("/media", mediaRateLimit.Middleware())

	mediaPostgresRepo := media_postgres_repository.NewRepository(a.infrastructure.Postgres)
	mediaService := media_service.NewService(mediaPostgresRepo, outboxPostgresRepository, a.infrastructure.S3)
	mediaHttpHandler := media_transport_http.NewHttpHandler(mediaService)

	mediaRouter.RegisterRoutes(mediaHttpHandler.Routes()...)

	// Users route
	logFeatureInit("users", "/api/v1/users")
	usersRouter := core_http_server.NewRouter("/users",
		rootRateLimit.Middleware(), core_http_middleware.Authorization(a.infrastructure.TokenIssuer))

	usersPostgresRepository := users_postgres_repository.NewRepository(a.infrastructure.Postgres)
	usersService := users_service.NewService(mediaService, usersPostgresRepository, outboxPostgresRepository)
	usersHttpHandler := users_transport_http.NewHttpHandler(usersService)

	usersRouter.RegisterRoutes(usersHttpHandler.Routes()...)

	// Settings route
	logFeatureInit("settings", "/api/v1/settings")
	settingsRouter := core_http_server.NewRouter("/settings",
		rootRateLimit.Middleware(), core_http_middleware.Authorization(a.infrastructure.TokenIssuer))

	settingsPostgresRepository := settings_postgres_repository.NewRepository(a.infrastructure.Postgres)
	settingsService := settings_service.NewService(settingsPostgresRepository)
	settingsHttpHandler := settings_transport_http.NewHttpHandler(settingsService)

	settingsRouter.RegisterRoutes(settingsHttpHandler.Routes()...)

	// Client route
	logFeatureInit("client", "/api/v1/client")
	clientRouter := core_http_server.NewRouter("/clients", rootRateLimit.Middleware())

	clientPostgresRepository := client_postgres_repository.NewRepository(a.infrastructure.Postgres)
	clientService := client_service.NewService(clientPostgresRepository)
	clientHttpHandler := client_transport_http.NewHttpHandler(clientService, a.infrastructure.TokenIssuer)

	clientRouter.RegisterRoutes(clientHttpHandler.Routes()...)

	// Requests route
	logFeatureInit("requests", "/api/v1/requests")
	requestsRouter := core_http_server.NewRouter("/requests", rootRateLimit.Middleware())

	requestsMongoRepository := requests_mongo_repository.NewRepository(a.infrastructure.Mongo)
	requestsService := requests_service.NewService(mediaService, clientService, usersPostgresRepository, requestsMongoRepository)
	requestsHttpHandler := requests_transport_http.NewHttpHandler(requestsService, a.infrastructure.TokenIssuer)

	requestsRouter.RegisterRoutes(requestsHttpHandler.Routes()...)

	// Sessions route
	logFeatureInit("sessions", "/api/v1/sessions")
	sessionsRouter := core_http_server.NewRouter("/sessions", core_http_middleware.Authorization(a.infrastructure.TokenIssuer))

	sessionsRedisRepo := sessions_redis_repository.NewRepository(a.infrastructure.Redis)
	sessionsService := session_service.NewService(sessionsRedisRepo)
	sessionsHttpHandler := sessions_transport_http.NewHttpHandler(sessionsService)

	sessionsRouter.RegisterRoutes(sessionsHttpHandler.Routes()...)

	// Auth routes
	logFeatureInit("authentication", "/api/v1/auth")
	authRouter := core_http_server.NewRouter("/auth")

	authService := auth_service.NewService(usersPostgresRepository, settingsPostgresRepository,
		tokenPostgresRepository, ipInfoRepository, sessionsRedisRepo, a.infrastructure.TokenIssuer)
	authHttpHandler := auth_transport_http.NewHttpHandler(authService, a.infrastructure.TokenIssuer)

	authRouter.RegisterRoutes(authHttpHandler.Routes()...)

	a.features = &Features{
		Root:     rootRouter,
		Media:    mediaRouter,
		Users:    usersRouter,
		Settings: settingsRouter,
		Clients:  clientRouter,
		Requests: requestsRouter,
		Sessions: sessionsRouter,
		Auth:     authRouter,
	}

	return nil
}
