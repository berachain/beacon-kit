// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package client

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	eth "github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/log"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// EngineClient is a struct that holds a pointer to an Eth1Client.
type EngineClient[
	ExecutionPayloadDenebT engineprimitives.ExecutionPayload,
] struct {
	*eth.Eth1Client[ExecutionPayloadDenebT]

	cfg          *Config
	capabilities map[string]struct{}
	logger       log.Logger[any]
	jwtSecret    *jwt.Secret

	// engineCache is an all-in-one cache for data
	// that are retrieved by the EngineClient.
	engineCache *cache.EngineCache

	statusErrCond *sync.Cond
	statusErrMu   *sync.RWMutex
	statusErr     error
}

// New creates a new engine client EngineClient.
// It takes an Eth1Client as an argument and returns a pointer  to an
// EngineClient.
func New[ExecutionPayloadDenebT engineprimitives.ExecutionPayload](
	cfg *Config,
	logger log.Logger[any],
	jwtSecret *jwt.Secret,
) *EngineClient[ExecutionPayloadDenebT] {
	ec := &EngineClient[ExecutionPayloadDenebT]{
		cfg:          cfg,
		logger:       logger,
		jwtSecret:    jwtSecret,
		Eth1Client:   new(eth.Eth1Client[ExecutionPayloadDenebT]),
		capabilities: make(map[string]struct{}),
		statusErrMu:  new(sync.RWMutex),
	}
	ec.statusErrCond = sync.NewCond(ec.statusErrMu)

	// If the engine cache is not set, we create a new one.
	if ec.engineCache == nil {
		ec.engineCache = cache.NewEngineCacheWithDefaultConfig()
	}

	return ec
}

// Start starts the engine client.
func (s *EngineClient[ExecutionPayloadDenebT]) Start(ctx context.Context) {
	for {
		s.logger.Info("waiting for execution client to start 🍺🕔",
			"dial-url", s.cfg.RPCDialURL)
		if err := s.setupExecutionClientConnection(ctx); err != nil {
			s.statusErrMu.Lock()
			s.statusErr = err
			s.statusErrMu.Unlock()
			time.Sleep(s.cfg.RPCStartupCheckInterval)
			continue
		}
		break
	}

	// Get the chain ID from the execution client.
	chainID, err := s.ChainID(ctx)
	if err != nil {
		s.logger.Error("failed to get chain ID", "err", err)
		return
	}

	// Log the chain ID.
	s.logger.Info(
		"connected to execution client 🔌",
		"dial-url",
		s.cfg.RPCDialURL.String(),
		"chain-id",
		chainID.Uint64(),
		"required-chain-id",
		s.cfg.RequiredChainID,
	)

	// Exchange capabilities with the execution client.
	if _, err = s.ExchangeCapabilities(ctx); err != nil {
		s.logger.Error("failed to exchange capabilities", "err", err)
	}

	// If we reached this point, the execution client is connected so we can
	// start the jwt refresh loop.
	go s.jwtRefreshLoop(ctx)
}

// Status verifies the chain ID via JSON-RPC. By proxy
// we will also verify the connection to the execution client.
func (s *EngineClient[ExecutionPayloadDenebT]) Status() error {
	s.statusErrMu.RLock()
	defer s.statusErrMu.RUnlock()
	return s.status(context.Background())
}

// status returns the status of the engine client.
func (s *EngineClient[ExecutionPayloadDenebT]) status(
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

// WaitForHealthy waits for the engine client to be healthy.
func (s *EngineClient[ExecutionPayloadDenebT]) WaitForHealthy(
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

// refreshUntilHealthy refreshes the engine client until it is healthy.
// TODO: remove after hack testing done.
func (s *EngineClient[ExecutionPayloadDenebT]) refreshUntilHealthy(
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

// Checks the chain ID of the execution client to ensure
// it matches local parameters of what Prysm expects.
func (s *EngineClient[ExecutionPayloadDenebT]) VerifyChainID(
	ctx context.Context,
) error {
	chainID, err := s.Client.ChainID(ctx)
	if err != nil {
		return err
	}

	if chainID.Uint64() != s.cfg.RequiredChainID {
		return errors.Newf(
			"wanted chain ID %d, got %d",
			s.cfg.RequiredChainID,
			chainID.Uint64(),
		)
	}

	return nil
}

// jwtRefreshLoop refreshes the JWT token for the execution client.
func (s *EngineClient[ExecutionPayloadDenebT]) jwtRefreshLoop(
	ctx context.Context,
) {
	ticker := time.NewTicker(s.cfg.RPCJWTRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.statusErrMu.Lock()
			if err := s.dialExecutionRPCClient(ctx); err != nil {
				s.logger.Error("failed to refresh JWT token", "err", err)
				//#nosec:G703 wtf is even this problem here.
				s.statusErr = errors.Newf(
					"%w: failed to refresh JWT token",
					err,
				)
			} else {
				s.statusErr = nil
			}
			s.statusErrMu.Unlock()
		}
	}
}

// setupExecutionClientConnections dials the execution client and
// ensures the chain ID is correct.
func (s *EngineClient[ExecutionPayloadDenebT]) setupExecutionClientConnection(
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

// DialExecutionRPCClient dials the execution client's RPC endpoint.
func (s *EngineClient[ExecutionPayloadDenebT]) dialExecutionRPCClient(
	ctx context.Context,
) error {
	var (
		client *rpc.Client
		err    error
	)

	// Build an http.Header with the JWT token attached.
	header, err := s.buildJWTHeader()
	if err != nil {
		return err
	}

	// Dial the execution client based on the URL scheme.
	switch s.cfg.RPCDialURL.Scheme {
	case "http", "https":
		client, err = rpc.DialOptions(
			ctx, s.cfg.RPCDialURL.String(), rpc.WithHeaders(header),
		)
	case "", "ipc":
		client, err = rpc.DialIPC(ctx, s.cfg.RPCDialURL.String())
	default:
		return errors.Newf(
			"no known transport for URL scheme %q",
			s.cfg.RPCDialURL.Scheme,
		)
	}

	// Check for an error when dialing the execution client.
	if err != nil {
		return err
	}

	s.Client = ethclient.NewClient(client)
	return nil
}

// buildJWTHeader builds an http.Header that has the JWT token
// attached for authorization.
//
//nolint:lll
func (s *EngineClient[ExecutionPayloadDenebT]) buildJWTHeader() (http.Header, error) {
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
