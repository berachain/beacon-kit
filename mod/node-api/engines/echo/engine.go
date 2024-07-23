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
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Engine struct {
	*echo.Echo
}

func New(e *echo.Echo) *Engine {
	return &Engine{
		Echo: e,
	}
}

func NewDefaultEngine() *Engine {
	engine := echo.New()
	engine.Use(middleware.CORSWithConfig(
		middleware.DefaultCORSConfig,
	))
	engine.Validator = &CustomValidator{
		Validator: ConstructValidator(),
	}
	return New(engine)
}

func (e *Engine) Run(addr string) error {
	return e.Echo.Start(addr)
}

func (e *Engine) RegisterRoutes(hs handlers.RouteSet[Context]) {
	group := e.Group(hs.BasePath)
	for _, route := range hs.Routes {
		group.Add(
			route.Method,
			route.Path,
			buildHandler(route),
		)
	}
}
