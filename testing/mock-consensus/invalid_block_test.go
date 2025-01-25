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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package mock_consensus_test

import (
	"context"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	comettypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/berachain/beacon-kit/log/phuslu"
)

func TestInitChainRequestsInvalidChainID(t *testing.T) {
	// Create a test node that
	node := newTestNode(t)

	var cometService *cometbft.Service[*phuslu.Logger]
	err := node.FetchService(&cometService)
	require.NoError(t, err)
	require.NotNil(t, cometService)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	request := &comettypes.InitChainRequest{
		ChainId: "80090",
	}
	_, err = cometService.InitChain(ctx, request)
	require.Error(t, err, "invalid chain-id on InitChain; expected: beacond-2061, got: 80090")
}

// TestProcessProposalRequestInvalidBlock tests the scenario where a peer sends us a block with an invalid timestamp
func TestProcessProposalRequestInvalidBlock(t *testing.T) {
	
}
