package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func assignDebugRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.GET("/eth/v2/debug/beacon/states/:state_id",
		h.NotImplemented)
	e.GET("/eth/v2/debug/beacon/states/heads",
		h.NotImplemented)
	e.GET("/eth/v1/debug/fork_choice",
		h.NotImplemented)
}
