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

	b *backend.Backend
	// exposed via getter for some tests.
	// TODO: consider extending this to other handlers
	beaconHandler *beaconapi.Handler
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

	// consensusService allows apis to access node state
	// and carry out all sorts of queries, including hystorical ones
	consensusService types.ConsensusService,
) *Server {
	apiLogger := logger
	if !config.Logging {
		apiLogger = noop.NewLogger[log.Logger]()
	}

	mware := middleware.NewDefaultMiddleware(apiLogger)

	// instantiate handlers and register their routes in the middleware
	b := backend.New(storageBackend, cs, cmtCfg, consensusService)
	beaconHandler := beaconapi.NewHandler(b, cs, apiLogger)
	mware.RegisterRoutes(beaconHandler.RouteSet())
	mware.RegisterRoutes(builderapi.NewHandler(apiLogger).RouteSet())
	mware.RegisterRoutes(configapi.NewHandler(b, apiLogger).RouteSet())
	mware.RegisterRoutes(debugapi.NewHandler(b, apiLogger).RouteSet())
	mware.RegisterRoutes(eventsapi.NewHandler(apiLogger).RouteSet())
	mware.RegisterRoutes(nodeapi.NewHandler(b, apiLogger).RouteSet())
	mware.RegisterRoutes(proofapi.NewHandler(b, apiLogger).RouteSet())

	return &Server{
		config:        config,
		logger:        logger,
		middleware:    mware,
		b:             b,
		beaconHandler: beaconHandler,
	}
}

// Start starts the API Server at the configured address.
func (s *Server) Start(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}

	// pre-load and cache all relevant node-api backend data
	if err := s.b.LoadData(); err != nil {
		return fmt.Errorf("failed loading api backend data: %w", err)
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

func (s *Server) GetBeaconHandler() *beaconapi.Handler {
	return s.beaconHandler
}
