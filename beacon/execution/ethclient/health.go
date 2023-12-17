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
	"fmt"
	"time"

	"github.com/prysmaticlabs/prysm/v4/io/logs"
)

// healthCheckPeriod defines the time interval for periodic health checks.
var healthCheckPeriod = 5 * time.Second

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
