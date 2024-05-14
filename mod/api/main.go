package main

import (
	backend "github.com/berachain/beacon-kit/mod/api/backend"
	server "github.com/berachain/beacon-kit/mod/api/server"
	handlers "github.com/berachain/beacon-kit/mod/api/server/handlers"
	validator "github.com/go-playground/validator/v10"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func NewServer(corsConfig middleware.CORSConfig, loggingConfig middleware.LoggerConfig) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	e.Validator = &handlers.CustomValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}
	server.UseMiddlewares(e, middleware.CORSWithConfig(corsConfig), middleware.LoggerWithConfig(loggingConfig))
	server.AssignRoutes(e, handlers.RouteHandlers{Backend: backend.Backend{}})
	return e
}

func run() {
	e := NewServer(middleware.DefaultCORSConfig, middleware.DefaultLoggerConfig)
	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	run()
}


