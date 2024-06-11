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
	"sync"
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
	// statusErrCond is a condition variable for the status error.
	statusErrCond *sync.Cond
	// statusErrMu is a mutex for the status error.
	statusErrMu *sync.RWMutex
	// statusErr is the status error of the engine client.
	statusErr error
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
	statusErrMu := new(sync.RWMutex)
	return &EngineClient[ExecutionPayloadT]{
		cfg:           cfg,
		logger:        logger,
		jwtSecret:     jwtSecret,
		Eth1Client:    new(ethclient.Eth1Client[ExecutionPayloadT]),
		capabilities:  make(map[string]struct{}),
		statusErrMu:   statusErrMu,
		statusErrCond: sync.NewCond(statusErrMu),
		engineCache:   cache.NewEngineCacheWithDefaultConfig(),
		eth1ChainID:   eth1ChainID,
		metrics:       newClientMetrics(telemetrySink, logger),
	}
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
	return s.initializeConnection(ctx)
}

// Status verifies the chain ID via JSON-RPC. By proxy
// we will also verify the connection to the execution client.
func (s *EngineClient[ExecutionPayloadT]) Status() error {
	s.statusErrMu.RLock()
	defer s.statusErrMu.RUnlock()
	return s.status(context.Background())
}

// WaitForHealthy waits for the engine client to be healthy.
func (s *EngineClient[ExecutionPayloadT]) WaitForHealthy(
	ctx context.Context,
) {
	s.statusErrMu.Lock()
	defer s.statusErrMu.Unlock()

	for s.status(ctx) != nil {
		go s.refreshUntilHealthy(ctx)
		select {
		case <-ctx.Done():
			return
		default:
			// Then we wait until we are blessed tf up.
			s.statusErrCond.Wait()
		}
	}
}

// VerifyChainID Checks the chain ID of the execution client to ensure
// it matches local parameters of what Prysm expects.
func (s *EngineClient[ExecutionPayloadT]) VerifyChainID(
	ctx context.Context,
) error {
	chainID, err := s.Client.ChainID(ctx)
	if err != nil {
		return err
	}

	if chainID.Uint64() != s.eth1ChainID.Uint64() {
		return errors.Newf(
			"wanted chain ID %d, got %d",
			s.eth1ChainID,
			chainID.Uint64(),
		)
	}

	return nil
}

// ============================== HELPERS ==============================

func (s *EngineClient[ExecutionPayloadT]) initializeConnection(
	ctx context.Context,
) error {
	// Initialize the connection to the execution client.
	var (
		err     error
		chainID *big.Int
	)
	for {
		s.logger.Info(
			"waiting for execution client to start 🍺🕔",
			"dial_url", s.cfg.RPCDialURL,
		)
		if err = s.setupExecutionClientConnection(ctx); err != nil {
			s.statusErrMu.Lock()
			s.statusErr = err
			s.statusErrMu.Unlock()
			time.Sleep(s.cfg.RPCStartupCheckInterval)
			s.logger.Error("failed to setup execution client", "err", err)
			continue
		}
		break
	}
	// Get the chain ID from the execution client.
	chainID, err = s.ChainID(ctx)
	if err != nil {
		s.logger.Error("failed to get chain ID", "err", err)
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

// setupExecutionClientConnections dials the execution client and
// ensures the chain ID is correct.
func (s *EngineClient[ExecutionPayloadT]) setupExecutionClientConnection(
	ctx context.Context,
) error {
	// Dial the execution client.
	if err := s.dialExecutionRPCClient(ctx); err != nil {
		return err
	}

	// Ensure the execution client is connected to the correct chain.
	if err := s.VerifyChainID(ctx); err != nil {
		s.Client.Close()
		if strings.Contains(err.Error(), "401 Unauthorized") {
			// We always log this error as it is a critical error.
			s.logger.Error(UnauthenticatedConnectionErrorStr)
		}
		return err
	}
	return nil
}

// ================================ Dialing ================================

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

// ================================ JWT ================================

// jwtRefreshLoop refreshes the JWT token for the execution client.
func (s *EngineClient[ExecutionPayloadT]) jwtRefreshLoop(
	ctx context.Context,
) {
	s.logger.Info("starting JWT refresh loop 🔄")
	ticker := time.NewTicker(s.cfg.RPCJWTRefreshInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			s.statusErrMu.Lock()
			if err := s.dialExecutionRPCClient(ctx); err != nil {
				s.logger.Error(
					"failed to refresh engine auth token",
					"err",
					err,
				)
				s.statusErr = ErrFailedToRefreshJWT
			} else {
				s.statusErr = nil
				s.logger.Info("successfully refreshed engine auth token")
			}
			s.statusErrMu.Unlock()
		}
	}
}

// buildJWTHeader builds an http.Header that has the JWT token
// attached for authorization.
//
//nolint:lll
func (s *EngineClient[ExecutionPayloadT]) buildJWTHeader() (http.Header, error) {
	header := make(http.Header)

	// Build the JWT token.
	token, err := buildSignedJWT(s.jwtSecret)
	if err != nil {
		s.logger.Error("failed to build JWT token", "err", err)
		return header, err
	}

	// Add the JWT token to the headers.
	header.Set("Authorization", "Bearer "+token)
	return header, nil
}

// Name returns the name of the engine client.
func (s *EngineClient[ExecutionPayloadT]) Name() string {
	return "engine-client"
}

// ================================ Info ================================

// status returns the status of the engine client.
func (s *EngineClient[ExecutionPayloadT]) status(
	ctx context.Context,
) error {
	// If the client is not started, we return an error.
	if s.Eth1Client.Client == nil {
		return ErrNotStarted
	}

	if s.statusErr == nil {
		// If we have an error, we will attempt
		// to verify the chain ID again.
		//#nosec:G703 wtf is even this problem here.
		s.statusErr = s.VerifyChainID(ctx)
	}

	if s.statusErr == nil {
		s.statusErrCond.Broadcast()
	}

	return s.statusErr
}

// refreshUntilHealthy refreshes the engine client until it is healthy.
// TODO: remove after hack testing done.
func (s *EngineClient[ExecutionPayloadT]) refreshUntilHealthy(
	ctx context.Context,
) {
	ticker := time.NewTicker(s.cfg.RPCStartupCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.status(ctx); err == nil {
				return
			}
		}
	}
}
