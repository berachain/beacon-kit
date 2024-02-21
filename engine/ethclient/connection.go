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

package ethclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// setupExecutionClientConnections dials the execution client and ensures the chain ID is correct.
func (s *Eth1Client) setupExecutionClientConnection(ctx context.Context) {
	// Dial the execution client.
	if err := s.dialExecutionRPCClient(ctx); err != nil {
		// This log gets spammy, we only log it when we first lose connection.
		if s.isConnected.Load() {
			s.logger.Error("could not dial execution client", "error", err)
		}
		s.isConnected.Store(false)
		return
	}

	// Ensure the execution client is connected to the correct chain.
	if err := s.ensureCorrectExecutionChain(ctx); err != nil {
		s.Client.Close()
		if strings.Contains(err.Error(), "401 Unauthorized") {
			// We always log this error as it is a critical error.
			s.logger.Error(UnauthenticatedConnectionErrorStr)
		} else if s.isConnected.Load() {
			// This log gets spammy, we only log it when we first lose connection.
			s.logger.Error("could not dial execution client", "error", err)
		}

		s.isConnected.Store(false)
		return
	}

	// If we reached here the client is connected and we mark as such.
	s.isConnected.Store(true)
}

// DialExecutionRPCClient dials the execution client's RPC endpoint.
func (s *Eth1Client) dialExecutionRPCClient(ctx context.Context) error {
	var client *rpc.Client

	// Construct the headers for the execution client.
	// New headers must be constructed each time the client is dialed
	// to periodically generate a new JWT token, as the existing one will eventually expire.
	headers, err := s.BuildHeaders()
	if err != nil {
		return err
	}

	// Dial the execution client based on the URL scheme.
	switch s.dialURL.Scheme {
	case "http", "https":
		client, err = rpc.DialOptions(
			ctx, s.dialURL.String(), rpc.WithHeaders(headers))
	case "", "ipc":
		client, err = rpc.DialIPC(ctx, s.dialURL.String())
	default:
		return fmt.Errorf("no known transport for URL scheme %q", s.dialURL.Scheme)
	}

	// Check for an error when dialing the execution client.
	if err != nil {
		return err
	}

	s.Client = ethclient.NewClient(client)
	s.GethRPCClient = client
	return nil
}

// tryConnectionAfter attempts a connection after a given interval.
func (s *Eth1Client) tryConnectionAfter(ctx context.Context, interval time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(interval):
		s.setupExecutionClientConnection(ctx)
	}
}
