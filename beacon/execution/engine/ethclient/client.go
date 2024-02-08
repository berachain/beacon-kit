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

	"github.com/pkg/errors"

	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/ethclient"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
)

const (
	// jwtLength is the length of the JWT token.
	jwtLength = 32
	// backOffPeriod is the time to wait before trying to reconnect with the eth1 node.
	backOffPeriod = 5
)

// eth1ClientConfig is a struct that holds the configuration for the Ethereum 1 client.
type eth1ClientConfig struct {
	chainID   uint64
	jwtSecret []byte
	headers   []string
	dialURL   *url.URL
}

// Eth1Client is a struct that holds the Ethereum 1 client and its configuration.
type Eth1Client struct {
	*ethclient.Client
	connectedETH1 bool
	cfg           *eth1ClientConfig
	ctx           context.Context
	logger        log.Logger
}

// NewEth1Client creates a new Ethereum 1 client with the provided context and options.
func NewEth1Client(ctx context.Context, opts ...Option) (*Eth1Client, error) {
	c := &Eth1Client{
		ctx: ctx,
		cfg: &eth1ClientConfig{},
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
	for {
		if err := s.setupExecutionClientConnections(); err != nil {
			s.logger.Info("Waiting for connection to execution client...",
				"dial-url", s.cfg.dialURL.String(), "err", err)
			time.Sleep(backOffPeriod * time.Second)
			continue
		}
		break
	}

	// Start the health check loop.
	go s.connectionHealthLoop(ctx)
}

func (s *Eth1Client) setupExecutionClientConnections() error {
	// Dial the execution client.
	if err := s.dialExecutionRPCClient(); err != nil {
		return errors.Wrap(err, "could not dial execution node")
	}

	// Ensure we have the correct chain ID connected.
	if err := s.ensureCorrectExecutionChain(); err != nil {
		s.Client.Close()
		errStr := err.Error()
		if strings.Contains(errStr, "401 Unauthorized") {
			errStr = "could not verify execution chain ID as your " +
				"connection is not authenticated. " +
				"If connecting to your execution client " +
				"via HTTP, you will need to set up JWT authentication..."
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

	// Build the headers for the execution client.
	// We have to build new headers every time we dial the client
	// since we need to periodically create a new JWT token since
	// the current one will eventually expire.
	headers, err := s.buildHeaders()
	if err != nil {
		return err
	}

	// Dial the execution client based on the URL scheme.
	switch s.cfg.dialURL.Scheme {
	case "http", "https":
		client, err = gethRPC.DialOptions(
			s.ctx, s.cfg.dialURL.String(), gethRPC.WithHeaders(headers))
	case "", "ipc":
		client, err = gethRPC.DialIPC(s.ctx, s.cfg.dialURL.String())
	default:
		return fmt.Errorf("no known transport for URL scheme %q", s.cfg.dialURL.Scheme)
	}

	// Check for an error when dialing the execution client.
	if err != nil {
		return err
	}

	// Attach the client to the struct.
	s.Client = ethclient.NewClient(client)
	return nil
}

// Every N seconds, defined as a backoffPeriod, attempts to re-establish an execution client
// connection and if this does not work, we fallback to the next endpoint if defined.
func (s *Eth1Client) pollConnectionStatus() {
	ticker := time.NewTicker(backOffPeriod * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.logger.Info("Trying to dial endpoint...", "dial-url", s.cfg.dialURL.String())
			currClient := s.Client.Client()
			if err := s.setupExecutionClientConnections(); err != nil {
				s.logger.Error("Could not connect to execution client endpoint", "error", err)
				continue
			}
			// Close previous client, if connection was successful.
			if currClient != nil {
				currClient.Close()
			}
			s.logger.Info("Connected to new endpoint", "dial-url", s.cfg.dialURL.String())
			return
		case <-s.ctx.Done():
			s.logger.Info("Received cancelled context,closing existing powchain service")
			return
		}
	}
}
