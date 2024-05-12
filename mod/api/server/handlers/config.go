package handlers

import (
	"net/http"
	// types "github.com/berachain/beacon-kit/mod/api/server/types"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetSpecConfig(c echo.Context) (err error) {
	rh.Backend.GetSpecConfig()
	return c.String(http.StatusOK, "GetSpecConfig")
}
