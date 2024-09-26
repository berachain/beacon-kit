// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package echo

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/engines/echo/validator"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Engine is an implementation of the API engine interface using Echo.
type Engine struct {
	*echo.Echo
	logger log.Logger
}

// New initializes a new API engine with the given Echo instance.
func New(e *echo.Echo) *Engine {
	return &Engine{
		Echo: e,
	}
}

// NewDefaultEngine returns a new default Echo Engine instance.
func NewDefaultEngine() *Engine {
	engine := echo.New()
	engine.Use(middleware.CORSWithConfig(
		middleware.DefaultCORSConfig,
	))
	engine.Validator = &validator.CustomValidator{
		Validator: validator.ConstructValidator(),
	}
	engine.HideBanner = true
	return New(engine)
}

// Run starts the Echo engine at the given address.
func (e *Engine) Run(addr string) error {
	return e.Echo.Start(addr)
}

// RegisterRoutes registers the given route set with the Echo engine.
func (e *Engine) RegisterRoutes(
	hs *handlers.RouteSet[Context],
	logger log.Logger,
) {
	e.logger = logger
	group := e.Group(hs.BasePath)
	for _, route := range hs.Routes {
		route.DecorateWithLogs(e.logger)
		group.Add(
			route.Method,
			route.Path,
			responseMiddleware(route),
		)
	}
}
