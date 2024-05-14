package handlers

import (
	"net/http"
	types "github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)
type RouteHandlers struct {
	Backend types.BackendHandlers
}

func (rh RouteHandlers) NotImplemented(c echo.Context) (err error) {
	return c.JSON(http.StatusNotImplemented, "Not Implemented")
}
