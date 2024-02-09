// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package eth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"cosmossdk.io/log"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/ethclient"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
)

const (
	// jwtLength is the length of the JWT token.
	jwtLength = 32
)

// Eth1Client is a struct that holds the Ethereum 1 client and its configuration.
type Eth1Client struct {
	ctx    context.Context
	logger log.Logger
	*ethclient.Client

	connectedETH1       bool
	chainID             uint64
	jwtSecret           [32]byte
	jwtRefreshInterval  time.Duration
	healthCheckInterval time.Duration
	dialURL             *url.URL
}

// NewEth1Client creates a new Ethereum 1 client with the provided context and options.
func NewEth1Client(ctx context.Context, opts ...Option) (*Eth1Client, error) {
	c := &Eth1Client{
		ctx: ctx,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	c.Start() // TODO: move this so it is on the cmd.Context.
	return c, nil
}

// Start the powchain service's main event loop.
func (s *Eth1Client) Start() {
	for {
		if err := s.setupExecutionClientConnection(); err != nil {
			s.logger.Info("Waiting for connection to execution client...",
				"dial-url", s.dialURL.String(), "err", err)
			continue
		}
		break
	}

	// Start the health check & jwt refresh loop.
	go s.healthCheckLoop()
	go s.jwtRefreshLoop()
}

// setupExecutionClientConnections dials the execution client and ensures the chain ID is correct.
func (s *Eth1Client) setupExecutionClientConnection() error {
	// Dial the execution client.
	if err := s.dialExecutionRPCClient(); err != nil {
		return errors.Wrap(err, "could not dial execution node")
	}

	// Ensure we have the correct chain ID connected.
	if err := s.ensureCorrectExecutionChain(); err != nil {
		s.Client.Close()
		errStr := err.Error()
		if strings.Contains(errStr, "401 Unauthorized") {
			errStr = UnauthenticatedConnectionErrorStr
		}
		return errors.Wrap(err, errStr)
	}

	// Mark the client as connected.
	s.updateConnectedETH1(true)
	return nil
}

// DialExecutionRPCClient dials the execution client's RPC endpoint.
func (s *Eth1Client) dialExecutionRPCClient() error {
	var client *gethRPC.Client

	// Construct the headers for the execution client.
	// New headers must be constructed each time the client is dialed
	// to periodically generate a new JWT token, as the existing one will eventually expire.
	headers, err := s.buildHeaders()
	if err != nil {
		return err
	}

	// Dial the execution client based on the URL scheme.
	switch s.dialURL.Scheme {
	case "http", "https":
		client, err = gethRPC.DialOptions(
			s.ctx, s.dialURL.String(), gethRPC.WithHeaders(headers))
	case "", "ipc":
		client, err = gethRPC.DialIPC(s.ctx, s.dialURL.String())
	default:
		return fmt.Errorf("no known transport for URL scheme %q", s.dialURL.Scheme)
	}

	// Check for an error when dialing the execution client.
	if err != nil {
		s.logger.Error("could not dial execution client", "error", err)
		return err
	}

	s.Client = ethclient.NewClient(client)
	return nil
}
