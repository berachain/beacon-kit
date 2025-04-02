//go:build test
// +build test

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

package core_test

import (
	"strconv"
	"testing"

	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/node-core/components"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/require"
)

func setupChain(t *testing.T) chain.Spec {
	t.Helper()

	t.Setenv(components.ChainSpecTypeEnvVar, components.DevnetChainSpecType)
	cs, err := components.ProvideChainSpec()
	require.NoError(t, err)

	return cs
}

//nolint:unused // may be used in the future.
func progressStateToSlot(
	t *testing.T,
	beaconState *statetransition.TestBeaconStateT,
	slot math.U64,
) {
	t.Helper()

	if slot == math.U64(0) {
		t.Fatal("for genesis slot, use InitializePreminedBeaconStateFromEth1")
	}

	err := beaconState.SetSlot(slot)
	require.NoError(t, err)
	err = beaconState.SetLatestBlockHeader(types.NewBeaconBlockHeader(
		slot,
		math.U64(0),
		common.Root{},
		common.Root{},
		common.Root{},
	))
	require.NoError(t, err)
}

func buildNextBlock(
	t *testing.T,
	beaconState *statetransition.TestBeaconStateT,
	eth1Data *types.Eth1Data,
	timestamp math.U64,
	blockDeposits types.Deposits,
	withdrawals ...*engineprimitives.Withdrawal,
) *types.BeaconBlock {
	t.Helper()

	// first update state root, similarly to what we do in processSlot
	parentBlkHeader, err := beaconState.GetLatestBlockHeader()
	require.NoError(t, err)
	root := beaconState.HashTreeRoot()
	parentBlkHeader.SetStateRoot(root)

	// build the block
	fv := version.Deneb1()
	versionable := types.NewVersionable(fv)
	blk, err := types.NewBeaconBlockWithVersion(
		parentBlkHeader.GetSlot()+1,
		parentBlkHeader.GetProposerIndex(),
		parentBlkHeader.HashTreeRoot(),
		fv,
	)
	require.NoError(t, err)

	// build the payload
	payload := &types.ExecutionPayload{
		Versionable:   versionable,
		Timestamp:     timestamp,
		ExtraData:     []byte("testing"),
		Transactions:  [][]byte{},
		Withdrawals:   withdrawals,
		BaseFeePerGas: math.NewU256(0),
	}
	parentBeaconBlockRoot := parentBlkHeader.HashTreeRoot()
	ethBlk, _, err := types.MakeEthBlock(payload, &parentBeaconBlockRoot)
	require.NoError(t, err)
	payload.BlockHash = common.ExecutionHash(ethBlk.Hash())

	require.NoError(t, err)
	blk.Body = &types.BeaconBlockBody{
		Versionable:      versionable,
		ExecutionPayload: payload,
		Eth1Data:         eth1Data,
		Deposits:         blockDeposits,
	}
	return blk
}

func generateTestExecutionAddress(
	t *testing.T,
	rndSeed int,
) (types.WithdrawalCredentials, int) {
	t.Helper()

	addrStr := strconv.Itoa(rndSeed)
	addrBytes := bytes.ExtendToSize([]byte(addrStr), bytes.B20Size)
	execAddr, err := bytes.ToBytes20(addrBytes)
	require.NoError(t, err)
	rndSeed++
	return types.NewCredentialsFromExecutionAddress(
		common.ExecutionAddress(execAddr),
	), rndSeed
}

func generateTestPK(t *testing.T, rndSeed int) (bytes.B48, int) {
	t.Helper()
	keyStr := strconv.Itoa(rndSeed)
	keyBytes := bytes.ExtendToSize([]byte(keyStr), bytes.B48Size)
	key, err := bytes.ToBytes48(keyBytes)
	require.NoError(t, err)
	rndSeed++
	return key, rndSeed
}

func moveToEndOfEpoch(
	t *testing.T,
	tip *types.BeaconBlock,
	cs chain.Spec,
	sp *statetransition.TestStateProcessorT,
	st *statetransition.TestBeaconStateT,
	ctx core.ReadOnlyContext,
	depRoot common.Root,
) *types.BeaconBlock {
	t.Helper()
	blk := tip
	currEpoch := cs.SlotToEpoch(blk.GetSlot())
	for currEpoch == cs.SlotToEpoch(blk.GetSlot()+1) {
		blk = buildNextBlock(
			t,
			st,
			types.NewEth1Data(depRoot),
			blk.Body.ExecutionPayload.Timestamp+1,
			[]*types.Deposit{},
			st.EVMInflationWithdrawal(blk.GetSlot()+1),
		)

		vals, err := sp.Transition(ctx, st, blk)
		require.NoError(t, err)
		require.Empty(t, vals) // no vals changes expected before next epoch
	}
	return blk
}
