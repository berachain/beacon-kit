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

package handlers

import (
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-api/types/context"
)

// handlerFn enforces a signature for all handler functions.
type handlerFn[ContextT context.Context] func(c ContextT) (any, error)

// Handlers is an interface that all handlers must implement.
type Handlers[ContextT context.Context] interface {
	// RegisterRoutes is a method that registers the routes for the handler.
	RegisterRoutes()
	RouteSet() RouteSet[ContextT]
}

// BaseHandler is a base handler for all handlers. It abstracts the route set
// and logger from the handler.
type BaseHandler[ContextT context.Context] struct {
	routes RouteSet[ContextT]
	logger log.APILogger[any]
}

// NewBaseHandler initializes a new base handler with the given routes and
// logger.
func NewBaseHandler[ContextT context.Context](
	routes RouteSet[ContextT],
	logger log.APILogger[any],
) *BaseHandler[ContextT] {
	return &BaseHandler[ContextT]{
		routes: routes,
		logger: logger,
	}
}

// RouteSet returns the route set for the base handler.
func (b *BaseHandler[ContextT]) RouteSet() RouteSet[ContextT] {
	return b.routes
}

// Logger is used to access the logger for the base handler.
func (b *BaseHandler[ContextT]) Logger() log.APILogger[any] {
	return b.logger
}

// AddRoutes adds the given slice of routes to the base handler.
func (b *BaseHandler[ContextT]) AddRoutes(
	routes []Route[ContextT],
) {
	b.routes.Routes = append(b.routes.Routes, routes...)
}
