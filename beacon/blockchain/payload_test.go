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

package blockchain_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	bcmocks "github.com/berachain/beacon-kit/beacon/blockchain/mocks"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	gethprimitives "github.com/berachain/beacon-kit/geth-primitives"
	bemocks "github.com/berachain/beacon-kit/node-api/backend/mocks"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core"
	stmocks "github.com/berachain/beacon-kit/state-transition/core/mocks"
	"github.com/berachain/beacon-kit/state-transition/core/state"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// When we reject a block and we have optimistic payload building enabled
// we must make sure that a few beacon state quantities are duly pre-processed
// before building the block.
func TestOptimisticBlockBuildingRejectedBlockStateChecks(t *testing.T) {
	t.Parallel()

	optimisticPayloadBuilds := true // key to this test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	chain, st, _, ctx, _, b, sb, eng, depStore := setupOptimisticPayloadTests(t, cs, optimisticPayloadBuilds)
	sb.EXPECT().StateFromContext(mock.Anything).Return(st)
	sb.EXPECT().DepositStore().RunAndReturn(func() *depositstore.KVStore { return depStore })
	b.EXPECT().Enabled().Return(optimisticPayloadBuilds)

	// Note: test avoid calling chain.Start since it only starts the deposits
	// goroutine which is not really relevant for this test

	// Before processing any block it is mandatory to handle genesis
	genesisData := testProcessGenesis(t, cs, chain, ctx)

	// Finally create a block that will be rejected and
	// verify the state on top of which is next payload built
	var (
		consensusTime   = math.U64(time.Now().Unix())
		proposerAddress = []byte{'d', 'u', 'm', 'm', 'y'} // this will err on purpose
	)

	// Since this is the first block called post genesis
	// forceSyncUponProcess will be called.
	dummyPayloadID := &engineprimitives.PayloadID{1, 2, 3}
	eng.EXPECT().NotifyForkchoiceUpdate(mock.Anything, mock.Anything).Return(dummyPayloadID, nil)

	// we set just enough data in invalid block to let it pass
	// the first validations in chain before state processor is invoked
	invalidBlk := &ctypes.BeaconBlock{
		Slot: 1, // first block after genesis
		Body: &ctypes.BeaconBlockBody{
			ExecutionPayload: &ctypes.ExecutionPayload{
				Timestamp: math.U64(cs.GenesisTime() + 1),
			},
		},
	}

	// register async call to block building
	var wg sync.WaitGroup          // useful to make test wait on async checks
	var ch = make(chan struct{})   // useful to serialize build block goroutine and avoid data races
	stateRoot := st.HashTreeRoot() // track state root before the changes done by optimistic build
	latestHeader, err := st.GetLatestBlockHeader()
	require.NoError(t, err)
	latestHeader.SetStateRoot(st.HashTreeRoot())
	expectedParentBlockRoot := latestHeader.HashTreeRoot()

	b.EXPECT().RequestPayloadAsync(
		mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Run(
		func(
			_ context.Context,
			st *state.StateDB,
			slot, timestamp math.U64,
			parentBlockRoot common.Root,
			headEth1BlockHash, finalEth1BlockHash common.ExecutionHash,
		) {
			defer wg.Done()
			<-ch // wait for block verification to finish. This avoids data races over state reads
			genesisHeader := genesisData.ExecutionPayloadHeader
			genesisBlkHeader := core.GenesisBlockHeader(cs.GenesisForkVersion())
			genesisBlkHeader.SetStateRoot(stateRoot)

			require.Equal(t, timestamp, consensusTime+1)

			require.Equal(t, genesisHeader.GetBlockHash(), headEth1BlockHash)

			require.Equal(t, expectedParentBlockRoot, parentBlockRoot)

			require.Empty(t, finalEth1BlockHash)          // this is first block post genesis
			require.Equal(t, constants.GenesisSlot, slot) // genesis slot in state
			var stateSlot math.Slot
			stateSlot, err = st.GetSlot()
			require.NoError(t, err)
			require.Equal(t, constants.GenesisSlot, stateSlot)
		},
	).Return(nil, common.Version{0xff}, errors.New("does not matter")) // return values do not really matter in this test
	wg.Add(1)

	err = chain.VerifyIncomingBlock(
		ctx.ConsensusCtx(),
		invalidBlk,
		consensusTime,
		proposerAddress,
	)
	require.ErrorIs(t, err, core.ErrProposerMismatch)

	// unlock checks on block building goroutine and
	// wait for it to carry out all the checks
	ch <- struct{}{}
	wg.Wait()
}

// When we verify successfully a block and we have optimistic payload building enabled
// we must make sure that a few beacon state quantities are duly pre-processed
// before building the block.
func TestOptimisticBlockBuildingVerifiedBlockStateChecks(t *testing.T) {
	t.Parallel()

	optimisticPayloadBuilds := true // key to this test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	chain, st, cms, ctx, sp, b, sb, eng, depStore := setupOptimisticPayloadTests(t, cs, optimisticPayloadBuilds)
	sb.EXPECT().StateFromContext(mock.Anything).Return(st).Times(1) // only for genesis
	sb.EXPECT().DepositStore().RunAndReturn(func() *depositstore.KVStore { return depStore })
	b.EXPECT().Enabled().Return(optimisticPayloadBuilds)

	// Before processing any block it is mandatory to handle genesis
	genesisData := testProcessGenesis(t, cs, chain, ctx)

	// write genesis changes to make them available for next block
	//nolint:errcheck // false positive as this has no return value
	ctx.ConsensusCtx().(sdk.Context).MultiStore().(storetypes.CacheMultiStore).Write()

	// Finally create a block that will be rejected and
	// verify the state on top of which is next payload built
	var (
		consensusTime = math.U64(time.Now().Unix())
		proposer      = ctx.ProposerAddress()
	)

	// Since this is the first block called post genesis
	// forceSyncUponProcess will be called.
	dummyPayloadID := &engineprimitives.PayloadID{1, 2, 3}
	eng.EXPECT().NotifyForkchoiceUpdate(mock.Anything, mock.Anything).Return(dummyPayloadID, nil)

	// BUILD A VALID BLOCK (without polluting state st)
	sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
	buildState := state.NewBeaconStateFromDB(
		st.KVStore.WithContext(sdkCtx), cs, sdkCtx.Logger(), metrics.NewNoOpTelemetrySink(),
	)

	nextBlkTimestamp := math.U64(cs.GenesisTime() + 1)
	_, err = sp.ProcessSlots(buildState, constants.GenesisSlot+1)
	require.NoError(t, err)

	depositsRoot := ctypes.Deposits(genesisData.Deposits).HashTreeRoot()

	validBlk := buildNextBlock(
		t,
		cs,
		buildState,
		ctypes.NewEth1Data(depositsRoot),
		nextBlkTimestamp,
	)
	stateRoot, err := computeStateRoot( // fix state root in block
		ctx.ConsensusCtx(),
		proposer,
		consensusTime,
		sp,
		buildState,
		validBlk,
	)
	require.NoError(t, err)
	validBlk.SetStateRoot(stateRoot)
	// end of BUILD A VALID BLOCK

	// register async call to block building
	var wg sync.WaitGroup        // useful to make test wait on async checks
	var ch = make(chan struct{}) // useful to serialize build block goroutine and avoid data races
	b.EXPECT().RequestPayloadAsync(
		mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Run(
		func(
			_ context.Context,
			st *state.StateDB,
			slot, timestamp math.U64,
			parentBlockRoot common.Root,
			headEth1BlockHash, finalEth1BlockHash common.ExecutionHash,
		) {
			defer wg.Done()
			<-ch // wait for block verification to finish. This avoids data races over state reads
			require.Equal(t, timestamp, consensusTime+1)

			require.Equal(
				t,
				validBlk.GetBody().GetExecutionPayload().GetBlockHash(),
				headEth1BlockHash,
			)

			genesisHeader := genesisData.ExecutionPayloadHeader.GetBlockHash()
			require.Equal(t, genesisHeader, finalEth1BlockHash)

			require.Equal(t, validBlk.HashTreeRoot(), parentBlockRoot)

			var stateSlot math.Slot
			stateSlot, err = st.GetSlot()
			require.NoError(t, err)
			require.Equal(t, validBlk.Slot+1, stateSlot)
			require.Equal(t, slot, stateSlot)
		},
	).Return(nil, common.Version{0xff}, errors.New("does not matter")) // return values do not really matter in this test
	wg.Add(1)

	eng.EXPECT().NotifyNewPayload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
	sb.EXPECT().StateFromContext(mock.Anything).Return(st).Times(1)
	err = chain.VerifyIncomingBlock(
		ctx.ConsensusCtx(),
		validBlk,
		consensusTime,
		ctx.ProposerAddress(),
	)
	require.NoError(t, err)

	// unlock checks on block building goroutine and
	// wait for it to carry out all the checks
	ch <- struct{}{}
	wg.Wait()
}

func setupOptimisticPayloadTests(t *testing.T, cs chain.Spec, optimisticPayloadBuilds bool) (
	*blockchain.Service,
	*statetransition.TestBeaconStateT,
	storetypes.CommitMultiStore,
	core.ReadOnlyContext,
	*statetransition.TestStateProcessorT,
	*bcmocks.LocalBuilder,
	*bemocks.StorageBackend,
	*stmocks.ExecutionEngine,
	*depositstore.KVStore,
) {
	t.Helper()
	sp, st, depStore, ctx, cms, eng := statetransition.SetupTestState(t, cs)

	logger := log.NewNopLogger()
	ts := metrics.NewNoOpTelemetrySink()
	sb := bemocks.NewStorageBackend(t)
	b := bcmocks.NewLocalBuilder(t)

	chain := blockchain.NewService(
		sb,
		nil, // blockchain.BlobProcessor unused in this test
		nil, // deposit.Contract unused in this test
		logger,
		cs,
		eng,
		b,
		sp,
		ts,
		optimisticPayloadBuilds,
	)
	return chain, st, cms, ctx, sp, b, sb, eng, depStore
}

func testProcessGenesis(
	t *testing.T,
	cs chain.Spec,
	chain *blockchain.Service,
	ctx core.ReadOnlyContext,
) *ctypes.Genesis {
	t.Helper()

	// TODO: I had to manually align default genesis and cs specs
	// Check if this is correct/necessary
	genesisData := ctypes.DefaultGenesis(cs.GenesisForkVersion())
	genesisData.ExecutionPayloadHeader.Timestamp = math.U64(cs.GenesisTime())
	genesisData.Deposits = []*ctypes.Deposit{
		{
			Pubkey: [48]byte{0x01},
			Amount: cs.MaxEffectiveBalance(),
			Credentials: ctypes.NewCredentialsFromExecutionAddress(
				common.ExecutionAddress{0x01},
			),
			Index: uint64(0),
		},
	}
	genBytes, err := json.Marshal(genesisData)
	require.NoError(t, err)
	_, err = chain.ProcessGenesisData(ctx.ConsensusCtx(), genBytes)
	require.NoError(t, err)
	return genesisData
}

func buildNextBlock(
	t *testing.T,
	cs chain.Spec,
	st *state.StateDB,
	eth1Data *ctypes.Eth1Data,
	timestamp math.U64,
) *ctypes.BeaconBlock {
	t.Helper()
	require.NotNil(t, cs)

	parentBlkHeader, err := st.GetLatestBlockHeader()
	require.NoError(t, err)
	nextBlockSlot := parentBlkHeader.GetSlot() + 1
	nextBlockEpoch := cs.SlotToEpoch(nextBlockSlot)

	randaoMix, err := st.GetRandaoMixAtIndex(nextBlockEpoch.Unwrap() % cs.EpochsPerHistoricalVector())
	require.NoError(t, err)

	// build the block
	fv := cs.ActiveForkVersionForTimestamp(timestamp)
	versionable := ctypes.NewVersionable(fv)
	blk, err := ctypes.NewBeaconBlockWithVersion(
		nextBlockSlot,
		parentBlkHeader.GetProposerIndex(),
		parentBlkHeader.HashTreeRoot(),
		fv,
	)
	require.NoError(t, err)

	// build the payload
	lph, err := st.GetLatestExecutionPayloadHeader()
	require.NoError(t, err)

	// Check chain canonicity
	payload := &ctypes.ExecutionPayload{
		Versionable: versionable,
		Timestamp:   timestamp,
		ParentHash:  lph.GetBlockHash(),
		Random:      randaoMix,

		ExtraData:     []byte("testing"),
		Transactions:  [][]byte{},
		Withdrawals:   []*engineprimitives.Withdrawal{st.EVMInflationWithdrawal(timestamp)},
		BaseFeePerGas: math.NewU256(0),
	}
	parentBeaconBlockRoot := parentBlkHeader.HashTreeRoot()

	var (
		ethBlk    *gethprimitives.Block
		noExecReq = &ctypes.ExecutionRequests{}
	)
	if version.IsBefore(fv, version.Electra()) {
		ethBlk, _, err = ctypes.MakeEthBlock(payload, &parentBeaconBlockRoot)
		require.NoError(t, err)
	} else {
		encodedER, erErr := ctypes.GetExecutionRequestsList(noExecReq)
		require.NoError(t, erErr)
		require.NotNil(t, encodedER)
		ethBlk, _, err = ctypes.MakeEthBlockWithExecutionRequests(payload, &parentBeaconBlockRoot, encodedER)
		require.NoError(t, err)
	}
	payload.BlockHash = common.ExecutionHash(ethBlk.Hash())

	require.NoError(t, err)
	blk.Body = &ctypes.BeaconBlockBody{
		Versionable:      versionable,
		ExecutionPayload: payload,
		Eth1Data:         eth1Data,
	}
	if version.EqualsOrIsAfter(fv, version.Electra()) {
		err = blk.Body.SetExecutionRequests(noExecReq)
		require.NoError(t, err)
	}
	return blk
}

func computeStateRoot(
	ctx context.Context,
	proposerAddress []byte,
	consensusTime math.U64,
	sp *statetransition.TestStateProcessorT,
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
) (common.Root, error) {
	txCtx := transition.NewTransitionCtx(
		ctx,
		consensusTime,
		proposerAddress,
	).
		WithVerifyPayload(false).
		WithVerifyRandao(false).
		WithVerifyResult(false).
		WithMeterGas(false)

	//nolint:contextcheck // we need txCtx
	if _, err := sp.Transition(txCtx, st, blk); err != nil {
		return common.Root{}, err
	}

	return st.HashTreeRoot(), nil
}
