package handlers

import (
	"net/http"
	echo "github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetSpecConfig(c echo.Context) (err error) {
	rh.Backend.GetSpecConfig()
	return c.String(http.StatusOK, "GetSpecConfig")
}
