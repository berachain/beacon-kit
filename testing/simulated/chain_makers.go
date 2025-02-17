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

package simulated

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// GenerateBeaconChain generates a beacon chain similar to geths chain generation utility.
// TODO: refactor this to be more flexible.
func GenerateBeaconChain(
	t *testing.T,
	numBlocks int,
	wrapELBlock func(latestSignedBeaconBlock *types.SignedBeaconBlock, block *gethprimitives.Block) (*types.SignedBeaconBlock, error),
) []*types.SignedBeaconBlock {
	t.Helper()
	genesis := &core.Genesis{
		Config:    params.AllEthashProtocolChanges,
		Alloc:     gethtypes.GenesisAlloc{},
		ExtraData: []byte("test genesis"),
		Timestamp: 1,
	}
	_, elBlocks, _ := core.GenerateChainWithGenesis(genesis, beacon.NewFaker(), numBlocks, func(_ int, b *core.BlockGen) {
		b.SetCoinbase(gethcommon.Address{0})
	})

	signedBeaconBlocks := []*types.SignedBeaconBlock{}
	var lastSignedBeaconBlock *types.SignedBeaconBlock
	for i := range elBlocks {
		signedBeaconBlock, err := wrapELBlock(lastSignedBeaconBlock, elBlocks[i])
		require.NoError(t, err)
		signedBeaconBlocks = append(signedBeaconBlocks, signedBeaconBlock)
		lastSignedBeaconBlock = signedBeaconBlock
	}

	return signedBeaconBlocks
}

// BlockToExecutionPayload TODO
func BlockToExecutionPayload(block *gethtypes.Block) *types.ExecutionPayload {
	payload := types.ExecutionPayload{
		ParentHash:    common.NewExecutionHashFromHex(block.ParentHash().Hex()),
		FeeRecipient:  common.ExecutionAddress{},
		StateRoot:     common.Bytes32{},
		ReceiptsRoot:  common.Bytes32{},
		LogsBloom:     bytes.B256{},
		Random:        common.Bytes32{},
		Number:        0,
		GasLimit:      0,
		GasUsed:       0,
		Timestamp:     0,
		ExtraData:     nil,
		BaseFeePerGas: nil,
		BlockHash:     common.NewExecutionHashFromHex(block.Hash().Hex()),
		Transactions:  nil,
		Withdrawals:   nil,
		BlobGasUsed:   0,
		ExcessBlobGas: 0,
		EpVersion:     common.Version{},
	}
	return &payload
}
