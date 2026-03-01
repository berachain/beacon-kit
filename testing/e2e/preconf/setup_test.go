// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

//go:build e2e

package preconf_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/testing/e2e/suite"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	// Kurtosis service names.
	loadTestPreconfRPCEL   = "el-preconf-rpc-reth-0"
	loadTestPreconfRPCPort = "eth-json-rpc"
)

// PreconfE2ESuite tests the preconf system by sending real ETH
// transactions through the preconf RPC node, measuring flashblock
// latency, and verifying state consistency with the standard RPC.
type PreconfE2ESuite struct {
	suite.KurtosisE2ESuite
	preconfClient *ethclient.Client
	chainID       *big.Int
}

// TestPreconfE2ESuite runs the preconf e2e test suite.
func TestPreconfE2ESuite(t *testing.T) {
	suite.Run(t, new(PreconfE2ESuite))
}

// SetupSuite initializes the network with a dedicated sequencer and
// preconf RPC node, then discovers the preconf RPC endpoint.
func (s *PreconfE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions(suite.WithPreconfLoadConfig())

	// Discover preconf RPC EL node via Kurtosis port mapping.
	sCtx, err := s.Enclave().GetServiceContext(loadTestPreconfRPCEL)
	s.Require().NoError(err, "Should get preconf RPC EL service context")

	port, ok := sCtx.GetPublicPorts()[loadTestPreconfRPCPort]
	s.Require().True(ok, "Preconf RPC EL should expose eth-json-rpc port")

	preconfURL := fmt.Sprintf("http://0.0.0.0:%d", port.GetNumber())
	s.T().Logf("Preconf RPC EL URL: %s", preconfURL)

	s.preconfClient, err = types.DialWithPooling(preconfURL)
	s.Require().NoError(err, "Should connect to preconf RPC EL")
	s.T().Cleanup(func() { s.preconfClient.Close() })

	s.chainID, err = s.RPCClient().ChainID(s.Ctx())
	s.Require().NoError(err, "Should get chain ID")

	// Brief warmup: confirm network is producing blocks after funding.
	err = s.WaitForNBlockNumbers(1)
	s.Require().NoError(err, "Network should produce warmup blocks")
}
