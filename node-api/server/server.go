// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"fmt"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/log/noop"
	"github.com/berachain/beacon-kit/node-api/backend"
	"github.com/berachain/beacon-kit/node-api/handlers"
	beaconapi "github.com/berachain/beacon-kit/node-api/handlers/beacon"
	builderapi "github.com/berachain/beacon-kit/node-api/handlers/builder"
	configapi "github.com/berachain/beacon-kit/node-api/handlers/config"
	debugapi "github.com/berachain/beacon-kit/node-api/handlers/debug"
	eventsapi "github.com/berachain/beacon-kit/node-api/handlers/events"
	nodeapi "github.com/berachain/beacon-kit/node-api/handlers/node"
	proofapi "github.com/berachain/beacon-kit/node-api/handlers/proof"
	"github.com/berachain/beacon-kit/node-api/middleware"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/types"
	cmtcfg "github.com/cometbft/cometbft/config"
)

// Server is the API Server service.
type Server struct {
	config     Config
	logger     log.Logger
	middleware *middleware.Middleware

	// exposed via getter for some tests.
	// TODO: find a way to avoid this
	b *backend.Backend
}

// New initializes a new API Server with the given config, engine, and logger.
// It will inject a noop logger into the API handlers and engine if logging is
// disabled.
func New(
	config Config,
	logger log.Logger,

	// attributes to build handlers backend
	storageBackend *storage.Backend,
	cs chain.Spec,
	cmtCfg *cmtcfg.Config,
	consensusService types.ConsensusService,
) (*Server, error) {
	apiLogger := logger
	if !config.Logging {
		apiLogger = noop.NewLogger[log.Logger]()
	}

	mware := middleware.NewDefaultMiddleware()

	// instantiate handlers and register their routes in the middleware
	b, err := backend.New(storageBackend, cs, cmtCfg)
	if err != nil {
		return nil, fmt.Errorf("failed instantiating handlers backend: %w", err)
	}
	b.AttachQueryBackend(consensusService)

	var handlers []handlers.Handlers
	handlers = append(handlers, beaconapi.NewHandler(b))
	handlers = append(handlers, builderapi.NewHandler())
	handlers = append(handlers, configapi.NewHandler(b))
	handlers = append(handlers, debugapi.NewHandler(b))
	handlers = append(handlers, eventsapi.NewHandler())
	handlers = append(handlers, nodeapi.NewHandler(b))
	handlers = append(handlers, proofapi.NewHandler(b))

	for _, handler := range handlers {
		handler.RegisterRoutes(apiLogger)
		mware.RegisterRoutes(handler.RouteSet(), apiLogger)
	}

	return &Server{
		config:     config,
		logger:     logger,
		middleware: mware,
		b:          b,
	}, nil
}

// Start starts the API Server at the configured address.
func (s *Server) Start(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}

	go s.start(ctx)
	return nil
}

func (s *Server) start(ctx context.Context) {
	errCh := make(chan error)
	go func() {
		errCh <- s.middleware.Run(s.config.Address)
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

func (s *Server) Stop() error {
	return nil
}

// Name returns the name of the API server service.
func (s *Server) Name() string {
	return "node-api-server"
}

func (s *Server) GetBackend() *backend.Backend {
	return s.b
}
