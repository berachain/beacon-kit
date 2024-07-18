package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

func assignConfigRoutes[ValidatorT any](
	e *echo.Echo,
	h handlers.RouteHandlers[ValidatorT],
) {
	e.GET("/eth/v1/config/fork_schedule",
		h.NotImplemented)
	e.GET("/eth/v1/config/spec",
		h.NotImplemented)
	e.GET("/eth/v1/config/deposit_contract",
		h.NotImplemented)
}
