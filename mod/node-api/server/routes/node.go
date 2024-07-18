package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func aasignNodeRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.GET("/eth/v1/node/identity",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/:peer_id",
		h.NotImplemented)
	e.GET("/eth/v1/node/peers/peer_count",
		h.NotImplemented)
	e.GET("/eth/v1/node/version",
		h.NotImplemented)
	e.GET("/eth/v1/node/syncing",
		h.NotImplemented)
	e.GET("/eth/v1/node/health",
		h.NotImplemented)
}
