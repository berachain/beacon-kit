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
	"time"

	"github.com/prysmaticlabs/prysm/v4/io/logs"
)

// healthCheckPeriod defines the time interval for periodic health checks.
const healthCheckPeriod = 5 * time.Second

// ConnectedETH1 returns the connection status of the Ethereum 1 client.
func (s *Eth1Client) ConnectedETH1() bool {
	// Return the connection status of the Ethereum 1 client.
	return s.connectedETH1
}

// updateConnectedETH1 updates the connection status of the Ethereum 1 client.
func (s *Eth1Client) updateConnectedETH1(state bool) {
	// Update the connection status of the Ethereum 1 client.
	s.connectedETH1 = state
}

// Checks the chain ID of the execution client to ensure
// it matches local parameters of what Prysm expects.
func (s *Eth1Client) ensureCorrectExecutionChain(ctx context.Context) error {
	chainID, err := s.Client.ChainID(ctx)
	if err != nil {
		return err
	}

	if chainID.Uint64() != s.cfg.chainID {
		return fmt.Errorf("wanted chain ID %d, got %d", s.cfg.chainID, chainID.Uint64())
	}
	return nil
}

// connectionHealthLoop periodically checks the connection health of the execution client.
func (s *Eth1Client) connectionHealthLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := s.Client.ChainID(ctx); err != nil {
				s.logger.Error("eth1 connection health check failed",
					"dial-url", logs.MaskCredentialsLogging(s.cfg.currHTTPEndpoint.Url),
					"err", err,
				)
				s.pollConnectionStatus(ctx)
			}
			time.Sleep(healthCheckPeriod)
		}
	}
}
