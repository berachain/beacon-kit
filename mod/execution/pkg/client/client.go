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

package client

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	ethclient2 "github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient2"
	ethclient2rpc "github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient2/rpc"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
)

// EngineClient is a struct that holds a pointer to an Eth1Client.
type EngineClient[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	PayloadAttributesT PayloadAttributes,
] struct {
	*ethclient2.EthRPC[ExecutionPayloadT]
	// cfg is the supplied configuration for the engine client.
	cfg *Config
	// logger is the logger for the engine client.
	logger log.Logger[any]
	// eth1ChainID is the chain ID of the execution client.
	eth1ChainID *big.Int
	// clientMetrics is the metrics for the engine client.
	metrics *clientMetrics
	// capabilities is a map of capabilities that the execution client has.
	capabilities map[string]struct{}
	// engineCache is an all-in-one cache for data
	// that are retrieved by the EngineClient.
	engineCache *cache.EngineCache
}

// New creates a new engine client EngineClient.
// It takes an Eth1Client as an argument and returns a pointer  to an
// EngineClient.
func New[
	ExecutionPayloadT constraints.EngineType[ExecutionPayloadT],
	PayloadAttributesT PayloadAttributes,
](
	cfg *Config,
	logger log.Logger[any],
	jwtSecret *jwt.Secret,
	telemetrySink TelemetrySink,
	eth1ChainID *big.Int,
) *EngineClient[
	ExecutionPayloadT, PayloadAttributesT,
] {
	return &EngineClient[ExecutionPayloadT, PayloadAttributesT]{
		cfg:    cfg,
		logger: logger,
		EthRPC: ethclient2.New[ExecutionPayloadT](
			ethclient2rpc.NewClient(
				cfg.RPCDialURL.String(),
				ethclient2rpc.WithJWTSecret(jwtSecret)),
		),
		capabilities: make(map[string]struct{}),
		engineCache:  cache.NewEngineCacheWithDefaultConfig(),
		eth1ChainID:  eth1ChainID,
		metrics:      newClientMetrics(telemetrySink, logger),
	}
}

// Name returns the name of the engine client.
func (s *EngineClient[
	_, _,
]) Name() string {
	return "engine-client"
}

// Start the engine client.
func (s *EngineClient[
	_, _,
]) Start(
	ctx context.Context,
) error {
	go s.EthRPC.Start(ctx)

	s.logger.Info(
		"Initializing connection to the execution client...",
		"dial_url", s.cfg.RPCDialURL.String(),
	)

	// If the connection connection succeeds, we can skip the
	// connection initializaation loop.
	if err := s.initializeConnection(ctx); err == nil {
		return nil
	}

	// Attempt to initialize the connection to the execution client.
	ticker := time.NewTicker(s.cfg.RPCStartupCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.logger.Info(
				"Waiting for execution client to start... 🍺🕔",
				"dial_url", s.cfg.RPCDialURL,
			)
			if err := s.initializeConnection(ctx); err != nil {
				if errors.Is(err, ErrMismatchedEth1ChainID) {
					s.logger.Error(err.Error())
				}
				continue
			}
			return nil
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Helpers                                  */
/* -------------------------------------------------------------------------- */

// setupConnection dials the execution client and
// ensures the chain ID is correct.
func (s *EngineClient[
	_, _,
]) initializeConnection(
	ctx context.Context,
) error {
	var (
		err     error
		chainID *big.Int
	)

	defer func() {
		if err != nil {
			s.EthRPC.Close()
		}
	}()

	// After the initial dial, check to make sure the chain ID is correct.
	chainID, err = s.EthRPC.ChainID(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			// We always log this error as it is a critical error.
			s.logger.Error(UnauthenticatedConnectionErrorStr)
		}
		return err
	}

	if chainID.Uint64() != s.eth1ChainID.Uint64() {
		err = errors.Wrapf(
			ErrMismatchedEth1ChainID,
			"wanted chain ID %d, got %d",
			s.eth1ChainID,
			chainID.Uint64(),
		)
		return err
	}

	// Log the chain ID.
	s.logger.Info(
		"Connected to execution client 🔌",
		"dial_url",
		s.cfg.RPCDialURL.String(),
		"chain_id",
		chainID.Uint64(),
		"required_chain_id",
		s.eth1ChainID,
	)

	// Exchange capabilities with the execution client.
	if _, err = s.ExchangeCapabilities(ctx); err != nil {
		s.logger.Error("failed to exchange capabilities", "err", err)
		return err
	}
	return nil
}
