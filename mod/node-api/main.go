// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"github.com/berachain/beacon-kit/mod/node-api/server"
	"github.com/berachain/beacon-kit/mod/node-api/server/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer(corsConfig middleware.CORSConfig,
	loggingConfig middleware.LoggerConfig) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler
	server.UseMiddlewares(e,
		middleware.CORSWithConfig(corsConfig),
		middleware.LoggerWithConfig(loggingConfig))
	server.AssignRoutes(e, handlers.RouteHandlers{})
	return e
}

func run() {
	e := NewServer(middleware.DefaultCORSConfig, middleware.DefaultLoggerConfig)
	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	run()
}
