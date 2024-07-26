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

package server

import (
	"context"

	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	apicontext "github.com/berachain/beacon-kit/mod/node-api/server/context"
)

// Server is the API Server service.
type Server[
	ContextT apicontext.Context,
	EngineT Engine[ContextT, EngineT],
] struct {
	engine EngineT
	config Config
	logger log.Logger[any]
}

// New initializes a new API Server with the given config, engine, and logger.
// It will inject a noop logger into the API handlers and engine if logging is
// disabled.
func New[
	ContextT apicontext.Context,
	EngineT Engine[ContextT, EngineT],
](
	config Config,
	engine EngineT,
	logger log.Logger[any],
	handlers ...handlers.Handlers[ContextT],
) *Server[ContextT, EngineT] {
	apiLogger := logger
	if !config.Logging {
		apiLogger = noop.NewLogger[log.Logger[any]]()
	}
	for _, handler := range handlers {
		handler.RegisterRoutes(apiLogger)
		engine.RegisterRoutes(handler.RouteSet(), apiLogger)
	}
	return &Server[ContextT, EngineT]{
		engine: engine,
		config: config,
		logger: logger,
	}
}

// Start starts the API Server at the configured address.
func (s *Server[_, _]) Start(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}
	go s.start(ctx)
	return nil
}

func (s *Server[_, _]) start(ctx context.Context) {
	errCh := make(chan error)
	go func() {
		errCh <- s.engine.Run(s.config.Address)
	}()
	for {
		select {
		case err := <-errCh:
			s.logger.Error(err.Error())
		case <-ctx.Done():
			return
		}
	}
}

// Name returns the name of the API server service.
func (s *Server[_, _]) Name() string {
	return "node-api-server"
}
