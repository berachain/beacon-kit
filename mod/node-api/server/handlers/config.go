package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (rh RouteHandlers) GetConfigSpec(c echo.Context) (err error) {
	rh.Backend.GetConfigSpec(context.TODO())
	return c.String(http.StatusOK, "GetConfigSpec")
}
