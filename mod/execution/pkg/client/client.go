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
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/cache"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// EngineClient is a struct that holds a pointer to an Eth1Client.
type EngineClient[
	ExecutionPayloadDenebT engineprimitives.ExecutionPayload,
] struct {
	// Eth1Client is a struct that holds the Ethereum 1 client and
	// its configuration.
	*ethclient.Eth1Client[ExecutionPayloadDenebT]
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
	// IPC
	ipcListener net.Listener
}

// New creates a new engine client EngineClient.
// It takes an Eth1Client as an argument and returns a pointer  to an
// EngineClient.
func New[ExecutionPayloadDenebT engineprimitives.ExecutionPayload](
	cfg *Config,
	logger log.Logger[any],
	jwtSecret *jwt.Secret,
	telemetrySink TelemetrySink,
	eth1ChainID *big.Int,
) *EngineClient[ExecutionPayloadDenebT] {
	statusErrMu := new(sync.RWMutex)
	return &EngineClient[ExecutionPayloadDenebT]{
		cfg:           cfg,
		logger:        logger,
		jwtSecret:     jwtSecret,
		Eth1Client:    new(ethclient.Eth1Client[ExecutionPayloadDenebT]),
		capabilities:  make(map[string]struct{}),
		statusErrMu:   statusErrMu,
		statusErrCond: sync.NewCond(statusErrMu),
		engineCache:   cache.NewEngineCacheWithDefaultConfig(),
		eth1ChainID:   eth1ChainID,
		metrics:       newClientMetrics(telemetrySink, logger),
	}
}

func (s *EngineClient[ExecutionPayloadDenebT]) StartWithIPC(
	ctx context.Context,
) error {
	if err := s.initializeConnection(ctx); err != nil {
		return err
	}
	if s.cfg.RPCDialURL.IsIPC() {
		s.startIPCServer(ctx)
	}
	return nil
}

// StartWithHTTP starts the engine client.
func (s *EngineClient[ExecutionPayloadDenebT]) Start(
	ctx context.Context,
) error {
	// This is not required for IPC connections.
	if s.cfg.RPCDialURL.IsHTTP() || s.cfg.RPCDialURL.IsHTTPS() {
		// If we are in a JWT mode, we will start the JWT refresh loop.
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
	go s.syncCheck(ctx)
	return s.initializeConnection(ctx)
}

// Status verifies the chain ID via JSON-RPC. By proxy
// we will also verify the connection to the execution client.
func (s *EngineClient[ExecutionPayloadDenebT]) Status() error {
	s.statusErrMu.RLock()
	defer s.statusErrMu.RUnlock()
	return s.status()
}

// WaitForHealthy waits for the engine client to be healthy.
func (s *EngineClient[ExecutionPayloadDenebT]) WaitForHealthy(
	ctx context.Context,
) {
	s.statusErrMu.Lock()
	defer s.statusErrMu.Unlock()

	for s.status() != nil {
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

// Checks the chain ID of the execution client to ensure
// it matches local parameters of what Prysm expects.
func (s *EngineClient[ExecutionPayloadDenebT]) VerifyChainID(
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

// syncCheck checks the sync status of the execution client.
func (s *EngineClient[ExecutionPayloadDenebT]) syncCheck(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.SyncCheckInterval)
	defer ticker.Stop()
	s.logger.Info(
		"starting sync check routine",
		"interval",
		s.cfg.SyncCheckInterval,
	)
	for {
		select {
		case <-ticker.C:
			syncProgress, err := s.SyncProgress(ctx)
			if err != nil {
				s.logger.Error("failed to get sync progress", "err", err)
				continue
			}

			s.statusErrMu.Lock()
			if syncProgress == nil || syncProgress.Done() {
				s.logger.Info("execution client is in sync 🍻")
				s.statusErr = nil
			} else {
				s.logger.Warn("execution client is syncing", "sync_progress", syncProgress)
				s.statusErr = engineerrors.ErrExecutionClientIsSyncing
			}
			s.statusErrMu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (s *EngineClient[ExecutionPayloadDenebT]) initializeConnection(
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

// ================================ Dialing ================================

// dialExecutionRPCClient dials the execution client's RPC endpoint.
func (s *EngineClient[ExecutionPayloadDenebT]) dialExecutionRPCClient(
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
			ctx, s.cfg.RPCDialURL.String()); err != nil {
			return err
		}
	default:
		return errors.Newf(
			"no known transport for URL scheme %q",
			s.cfg.RPCDialURL.Scheme,
		)
	}

	// Refresh the execution client with the new client.
	s.Eth1Client, err = ethclient.NewFromRPCClient[ExecutionPayloadDenebT](
		client,
	)
	return err
}

// ================================ JWT ================================

// jwtRefreshLoop refreshes the JWT token for the execution client.
func (s *EngineClient[ExecutionPayloadDenebT]) jwtRefreshLoop(
	ctx context.Context,
) {
	s.logger.Info("starting JWT refresh loop 🔄")
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

func (s *EngineClient[ExecutionPayloadDenebT]) Name() string {
	return "EngineClient"
}

// ================================ IPC ================================

//

func (s *EngineClient[ExecutionPayloadDenebT]) startIPCServer(
	ctx context.Context,
) {
	if s.cfg.RPCDialURL == nil || !s.cfg.RPCDialURL.IsIPC() {
		s.logger.Error("IPC server not started, invalid IPC URL")
		return
	}
	// remove existing socket file if exists
	// alternatively we can use existing one by checking for os.IsNotExist(err)
	if _, err := os.Stat(s.cfg.RPCDialURL.Path); err != nil {
		s.logger.Info(
			"Removing existing IPC file",
			"path",
			s.cfg.RPCDialURL.Path,
		)

		if err = os.Remove(s.cfg.RPCDialURL.Path); err != nil {
			s.logger.Error("failed to remove existing IPC file", "err", err)
			return
		}
	}

	// use UDS for IPC
	listener, err := net.Listen("unix", s.cfg.RPCDialURL.Path)
	if err != nil {
		s.logger.Error("failed to listen on IPC socket", "err", err)
		return
	}
	s.ipcListener = listener

	// register the RPC server
	server := rpc.NewServer()
	if err = server.Register(s); err != nil {
		s.logger.Error("failed to register RPC server", "err", err)
		return
	}
	s.logger.Info("IPC server started", "path", s.cfg.RPCDialURL.Path)

	// start server in a goroutine
	go func() {
		for {
			// continuously accept incoming connections until context is
			// cancelled
			select {
			case <-ctx.Done():
				s.logger.Info("shutting down IPC server")
				return
			default:
				var conn net.Conn
				conn, err = listener.Accept()
				if err != nil {
					s.logger.Error(
						"failed to accept IPC connection",
						"err",
						err,
					)
					continue
				}
				go server.ServeConn(conn)
			}
		}
	}()
}

// ================================ Info ================================

// status returns the status of the engine client.
func (s *EngineClient[ExecutionPayloadDenebT]) status() error {
	// If the client is not started, we return an error.
	if s.Eth1Client.Client == nil {
		return ErrNotStarted
	}

	if s.statusErr == nil {
		s.statusErrCond.Broadcast()
	}

	return s.statusErr
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
			if err := s.status(); err == nil {
				return
			}
		}
	}
}
