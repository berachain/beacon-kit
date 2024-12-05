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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
)

// handlerFn enforces a signature for all handler functions.
type handlerFn[ContextT any] func(c ContextT) (any, error)

// Handlers is an interface that all handlers must implement.
type Handlers[ContextT any] interface {
	// RegisterRoutes is a method that registers the routes for the handler.
	RegisterRoutes(logger log.Logger)
	RouteSet() *RouteSet[ContextT]
}

// BaseHandler is a base handler for all handlers. It abstracts the route set
// and logger from the handler.
type BaseHandler[ContextT any] struct {
	routes *RouteSet[ContextT]
	logger log.Logger
}

// NewBaseHandler initializes a new base handler with the given routes and
// logger.
func NewBaseHandler[ContextT any](
	routes *RouteSet[ContextT],
) *BaseHandler[ContextT] {
	return &BaseHandler[ContextT]{
		routes: routes,
	}
}

// NotImplemented is a placeholder for the beacon API.
func (b *BaseHandler[ContextT]) NotImplemented(ContextT) (any, error) {
	return nil, errors.New("not implemented")
}

// RouteSet returns the route set for the base handler.
func (b *BaseHandler[ContextT]) RouteSet() *RouteSet[ContextT] {
	return b.routes
}

// Logger is used to access the logger for the base handler.
func (b *BaseHandler[ContextT]) Logger() log.Logger {
	return b.logger
}

func (b *BaseHandler[ContextT]) SetLogger(logger log.Logger) {
	b.logger = logger
}

// AddRoutes adds the given slice of routes to the base handler.
func (b *BaseHandler[ContextT]) AddRoutes(
	routes []*Route[ContextT],
) {
	b.routes.Routes = append(b.routes.Routes, routes...)
}
