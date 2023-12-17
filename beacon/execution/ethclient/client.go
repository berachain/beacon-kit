// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package eth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/io/logs"
	"github.com/prysmaticlabs/prysm/v4/network"
	"github.com/prysmaticlabs/prysm/v4/network/authorization"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/ethclient"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
)

// backOffPeriod is the time to wait before trying to reconnect with the eth1 node.
var (
	backOffPeriod = 10 * time.Second
)

const (
	// jwtLength is the length of the JWT token.
	jwtLength = 32
)

// Eth1Client is a struct that holds the Ethereum 1 client and its configuration.
type Eth1Client struct {
	*ethclient.Client
	connectedETH1 bool
	cfg           *Eth1ClientConfig
	ctx           context.Context
	rpcClient     *gethRPC.Client
	logger        log.Logger
}

// Eth1ClientConfig is a struct that holds the configuration for the Ethereum 1 client.
type Eth1ClientConfig struct {
	chainID          uint64
	headers          []string
	currHTTPEndpoint network.Endpoint
}

// NewEth1Client creates a new Ethereum 1 client with the provided context and options.
func NewEth1Client(ctx context.Context, opts ...Option) (*Eth1Client, error) {
	c := &Eth1Client{
		ctx: ctx,
		cfg: &Eth1ClientConfig{},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.Start(ctx)
	return c, nil
}

// Start the powchain service's main event loop.
func (s *Eth1Client) Start(ctx context.Context) {
	if err := s.setupExecutionClientConnections(s.ctx, s.cfg.currHTTPEndpoint); err != nil {
		s.logger.Error("could not connect to execution endpoint", "err", err)
	}

	// Start the health check loop.
	go s.connectionHealthLoop(ctx)
}

func (s *Eth1Client) setupExecutionClientConnections(
	ctx context.Context, currEndpoint network.Endpoint,
) error {
	client, err := s.newRPCClientWithAuth(ctx, currEndpoint)
	if err != nil {
		return errors.Wrap(err, "could not dial execution node")
	}
	// Attach the clients to the service struct.
	s.Client = ethclient.NewClient(client)
	s.rpcClient = client

	// Ensure we have the correct chain ID connected.
	if err = s.ensureCorrectExecutionChain(ctx); err != nil {
		client.Close()
		errStr := err.Error()
		if strings.Contains(errStr, "401 Unauthorized") {
			errStr = "could not verify execution chain ID as your " +
				"connection is not authenticated. " +
				"If connecting to your execution client " +
				"via HTTP, you will need to set up JWT authentication..."
		}
		return errors.Wrap(err, errStr)
	}

	s.updateConnectedETH1(true)
	return nil
}

// Every N seconds, defined as a backoffPeriod, attempts to re-establish an execution client
// connection and if this does not work, we fallback to the next endpoint if defined.
func (s *Eth1Client) pollConnectionStatus(ctx context.Context) {
	ticker := time.NewTicker(backOffPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.logger.Info("Trying to dial endpoint...", "dial-url",
				logs.MaskCredentialsLogging(s.cfg.currHTTPEndpoint.Url))
			currClient := s.rpcClient
			if err := s.setupExecutionClientConnections(ctx, s.cfg.currHTTPEndpoint); err != nil {
				s.logger.Error("Could not connect to execution client endpoint", "err", err)
				continue
			}
			// Close previous client, if connection was successful.
			if currClient != nil {
				currClient.Close()
			}
			s.logger.Info("Connected to new endpoint", "dial-url",
				logs.MaskCredentialsLogging(s.cfg.currHTTPEndpoint.Url))
			return
		case <-s.ctx.Done():
			s.logger.Info("Received cancelled context,closing existing powchain service")
			return
		}
	}
}

// Forces to retry an execution client connection.
func (s *Eth1Client) RetryExecutionClientConnection(ctx context.Context, _ error) {
	// s.runError = errors.Wrap(err, "retryExecutionClientConnection")
	s.updateConnectedETH1(false)
	// Back off for a while before redialing.
	time.Sleep(backOffPeriod)
	currClient := s.rpcClient
	if newErr := s.setupExecutionClientConnections(ctx, s.cfg.currHTTPEndpoint); newErr != nil {
		// s.runError = errors.Wrap(err, "setupExecutionClientConnections")
		return
	}
	// Close previous client, if connection was successful.
	if currClient != nil {
		currClient.Close()
	}
	// Reset run error in the event of a successful connection.
	// s.runError = nil
}

// Initializes an RPC connection with authentication headers.
func (s *Eth1Client) newRPCClientWithAuth(
	ctx context.Context, endpoint network.Endpoint,
) (*gethRPC.Client, error) {
	headers := http.Header{}
	if endpoint.Auth.Method != authorization.None {
		header, err := endpoint.Auth.ToHeaderValue()
		if err != nil {
			return nil, err
		}
		headers.Set("Authorization", header)
	}
	for _, h := range s.cfg.headers {
		if h == "" {
			continue
		}
		keyValue := strings.Split(h, "=")
		if len(keyValue) < 2 { //nolint:gomnd // it's okay.
			s.logger.Error("Incorrect HTTP header flag format. Skipping %v", keyValue[0])
			continue
		}
		headers.Set(keyValue[0], strings.Join(keyValue[1:], "="))
	}

	s.logger.Info("connecting to execution client...", "dial-url",
		logs.MaskCredentialsLogging(endpoint.Url))
	return network.NewExecutionRPCClient(ctx, endpoint, headers)
}
