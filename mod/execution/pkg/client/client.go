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
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// EngineClient is a struct that holds a pointer to an Eth1Client.
type EngineClient[
	ExecutionPayloadT interface {
		Empty(uint32) ExecutionPayloadT
		Version() uint32
		json.Marshaler
		json.Unmarshaler
	},
] struct {
	// Eth1Client is a struct that holds the Ethereum 1 client and
	// its configuration.
	*ethclient.Eth1Client[ExecutionPayloadT]
	// cfg is the supplied configuration for the engine client.
	cfg *Config
	// logger is the logger for the engine client.
	logger log.Logger[any]
	// jwtSecret is the JWT secret for the execution client.
	jwtSecret *jwt.Secret
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
func New[ExecutionPayloadT interface {
	Empty(uint32) ExecutionPayloadT
	Version() uint32
	json.Marshaler
	json.Unmarshaler
}](
	cfg *Config,
	logger log.Logger[any],
	jwtSecret *jwt.Secret,
	telemetrySink TelemetrySink,
	eth1ChainID *big.Int,
) *EngineClient[ExecutionPayloadT] {
	return &EngineClient[ExecutionPayloadT]{
		cfg:          cfg,
		logger:       logger,
		jwtSecret:    jwtSecret,
		Eth1Client:   new(ethclient.Eth1Client[ExecutionPayloadT]),
		capabilities: make(map[string]struct{}),
		engineCache:  cache.NewEngineCacheWithDefaultConfig(),
		eth1ChainID:  eth1ChainID,
		metrics:      newClientMetrics(telemetrySink, logger),
	}
}

// Name returns the name of the engine client.
func (s *EngineClient[ExecutionPayloadT]) Name() string {
	return "engine-client"
}

// Start the engine client.
func (s *EngineClient[ExecutionPayloadT]) Start(
	ctx context.Context,
) error {
	if s.cfg.RPCDialURL.IsHTTP() || s.cfg.RPCDialURL.IsHTTPS() {
		// If we are dialing with HTTP(S), start the JWT refresh loop.
		defer func() {
			if s.jwtSecret == nil {
				s.logger.Warn(
					"JWT secret not provided for http(s) connection" +
						" - please verify your configuration settings",
				)
				return
			}
			go s.jwtRefreshLoop(ctx)
		}()
	}

	s.logger.Info(
		"initializing connection to the execution client...",
		"dial_url", s.cfg.RPCDialURL.String(),
	)

	// If the connection connection succeeds, we can skip the
	// connection initialization loop.
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
				"waiting for execution client to start... 🍺🕔",
				"dial_url", s.cfg.RPCDialURL,
			)
			if err := s.initializeConnection(ctx); err != nil {
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
func (s *EngineClient[ExecutionPayloadT]) initializeConnection(
	ctx context.Context,
) error {
	var (
		err     error
		chainID *big.Int
	)

	defer func() {
		if err != nil {
			s.Client.Close()
		}
	}()

	// Dial the execution client.
	if err = s.dialExecutionRPCClient(ctx); err != nil {
		return err
	}

	// After the initial dial, check to make sure the chain ID is correct.
	chainID, err = s.Client.ChainID(ctx)
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
		"connected to execution client 🔌",
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

/* -------------------------------------------------------------------------- */
/*                                   Dialing                                  */
/* -------------------------------------------------------------------------- */

// dialExecutionRPCClient dials the execution client's RPC endpoint.
func (s *EngineClient[ExecutionPayloadT]) dialExecutionRPCClient(
	ctx context.Context,
) error {
	var (
		client *ethrpc.Client
		err    error
	)

	// Dial the execution client based on the URL scheme.
	switch {
	case s.cfg.RPCDialURL.IsHTTP(), s.cfg.RPCDialURL.IsHTTPS():
		// Build an http.Header with the JWT token attached.
		if s.jwtSecret != nil {
			var header http.Header
			if header, err = s.buildJWTHeader(); err != nil {
				return err
			}
			if client, err = ethrpc.DialOptions(
				ctx, s.cfg.RPCDialURL.String(), ethrpc.WithHeaders(header),
			); err != nil {
				return err
			}
		} else {
			if client, err = ethrpc.DialContext(
				ctx, s.cfg.RPCDialURL.String()); err != nil {
				return err
			}
		}
	case s.cfg.RPCDialURL.IsIPC():
		if client, err = ethrpc.DialIPC(
			ctx, s.cfg.RPCDialURL.Path); err != nil {
			s.logger.Error("failed to dial IPC", "err", err)
			return err
		}
	default:
		return errors.Newf(
			"no known transport for URL scheme %q",
			s.cfg.RPCDialURL.Scheme,
		)
	}

	// Refresh the execution client with the new client.
	s.Eth1Client, err = ethclient.NewFromRPCClient[ExecutionPayloadT](
		client,
	)
	return err
}
