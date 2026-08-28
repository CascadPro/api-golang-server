package app

import (
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

func (a *App) initHttp() error {
	a.logger.Debug("initializing http server")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		a.logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(a.logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
		core_http_middleware.IP(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(
		a.features.Auth,
		a.features.Clients,
		a.features.Requests,
		a.features.Sessions,
		a.features.Settings,
		a.features.Users,
	)

	httpServer.RegisterRouters(
		a.features.Root,
		a.features.Media,
	)

	httpServer.RegisterApiRouters(apiVersionRouter)

	httpServer.RegisterSwagger()

	a.httpServer = httpServer

	return nil
}
