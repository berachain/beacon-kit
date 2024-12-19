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
	"sync"
	"time"

	"github.com/berachain/beacon-kit/errors"
	ethclient "github.com/berachain/beacon-kit/execution/client/ethclient"
	ethclientrpc "github.com/berachain/beacon-kit/execution/client/ethclient/rpc"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
)

// EngineClient is a struct that holds a pointer to an Eth1Client.
type EngineClient struct {
	*ethclient.Client
	// cfg is the supplied configuration for the engine client.
	cfg *Config
	// logger is the logger for the engine client.
	logger log.Logger
	// eth1ChainID is the chain ID of the execution client.
	eth1ChainID *big.Int
	// clientMetrics is the metrics for the engine client.
	metrics *clientMetrics
	// capabilities is a map of capabilities that the execution client has.
	capabilities map[string]struct{}
	// connected will be set to true when we have successfully connected
	// to the execution client.
	connectedMu sync.RWMutex
	connected   bool
}

// New creates a new engine client EngineClient.
// It takes an Eth1Client as an argument and returns a pointer  to an
// EngineClient.
func New(
	cfg *Config,
	logger log.Logger,
	jwtSecret *jwt.Secret,
	telemetrySink TelemetrySink,
	eth1ChainID *big.Int,
) *EngineClient {
	return &EngineClient{
		cfg:    cfg,
		logger: logger,
		Client: ethclient.New(
			ethclientrpc.NewClient(
				cfg.RPCDialURL.String(),
				ethclientrpc.WithJWTSecret(jwtSecret),
				ethclientrpc.WithJWTRefreshInterval(
					cfg.RPCJWTRefreshInterval,
				),
			)),
		capabilities: make(map[string]struct{}),
		eth1ChainID:  eth1ChainID,
		metrics:      newClientMetrics(telemetrySink, logger),
		connected:    false,
	}
}

// Name returns the name of the engine client.
func (s *EngineClient) Name() string {
	return "engine-client"
}

// Start the engine client.
func (s *EngineClient) Start(
	ctx context.Context,
) error {
	// Start the Client.
	go s.Client.Start(ctx)

	s.logger.Info(
		"Initializing connection to the execution client...",
		"dial_url", s.cfg.RPCDialURL.String(),
	)

	// If the connection connection succeeds, we can skip the
	// connection initialization loop.
	if err := s.verifyChainIDAndConnection(ctx); err == nil {
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
			if err := s.verifyChainIDAndConnection(ctx); err != nil {
				if errors.Is(err, ErrMismatchedEth1ChainID) {
					s.logger.Error(err.Error())
				}
				continue
			}
			s.connectedMu.Lock()
			s.connected = true
			s.connectedMu.Unlock()
			return nil
		}
	}
}

func (s *EngineClient) Stop() error {
	return nil
}

func (s *EngineClient) IsConnected() bool {
	s.connectedMu.RLock()
	defer s.connectedMu.RUnlock()
	return s.connected
}

func (s *EngineClient) HasCapability(capability string) bool {
	_, ok := s.capabilities[capability]
	return ok
}

/* -------------------------------------------------------------------------- */
/*                                   Helpers                                  */
/* -------------------------------------------------------------------------- */

// verifyChainID dials the execution client and
// ensures the chain ID is correct.
func (s *EngineClient) verifyChainIDAndConnection(
	ctx context.Context,
) error {
	var (
		err     error
		chainID math.U64
	)

	defer func() {
		if err != nil {
			err = s.Client.Close()
		}
	}()

	// After the initial dial, check to make sure the chain ID is correct.
	chainID, err = s.Client.ChainID(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			// We always log this error as it is a critical error.
			s.logger.Error(UnauthenticatedConnectionErrorStr)
		}
		return err
	}

	// TODO: consider validating once when config is set or
	// client is initialized
	if !s.eth1ChainID.IsUint64() {
		err = errors.Wrapf(
			errors.New("provided chain ID is not uint64"),
			s.eth1ChainID.String(),
		)
		return err
	}
	if chainID.Unwrap() != s.eth1ChainID.Uint64() {
		err = errors.Wrapf(
			ErrMismatchedEth1ChainID,
			"wanted chain ID %d, got %d",
			s.eth1ChainID,
			chainID,
		)
		return err
	}

	// Log the chain ID.
	s.logger.Info(
		"Connected to execution client 🔌",
		"dial_url",
		s.cfg.RPCDialURL.String(),
		"chain_id",
		chainID.Unwrap(),
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
