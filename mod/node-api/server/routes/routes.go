package routes

import (
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
)

// Assign assigns all routes to the echo instance.
func Assign[ValidatorT any](
	e *echo.Echo,
	handler handlers.RouteHandlers[ValidatorT],
) {
	assignBeaconRoutes(e, handler)
	assignBuilderRoutes(e, handler)
	assignConfigRoutes(e, handler)
	assignDebugRoutes(e, handler)
	assignEventsRoutes(e, handler)
	aasignNodeRoutes(e, handler)
	assignValidatorRoutes(e, handler)
	assignRewardsRoutes(e, handler)
}
