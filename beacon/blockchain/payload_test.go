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
	"errors"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/node-api/backend/mocks"
	nodemetrics "github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/state-transition/core"
	depositstore "github.com/berachain/beacon-kit/storage/deposit"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// When we reject a block and we have optimistic payload building enabled
// We must make sure that Latest Execution Payload Header is duly pre-processed
// before building the block. (case for accepted block below)
func TestOptimisticBlockBuildingRejectedBlockStateChecks(t *testing.T) {
	t.Parallel()

	optimisticPayloadBuilds := true // key to this test
	cs, err := spec.MainnetChainSpec()
	require.NoError(t, err)

	sp, st, depStore, ctx, _, eng := statetransition.SetupTestState(t, cs)

	logger := log.NewNopLogger()
	ts := nodemetrics.NewNoOpTelemetrySink()
	sb := mocks.NewStorageBackend(t)
	sb.EXPECT().StateFromContext(mock.Anything).Return(st)
	sb.EXPECT().DepositStore().RunAndReturn(func() *depositstore.KVStore { return depStore })

	fb := &fakeBuilder{
		enabled: optimisticPayloadBuilds,
	}

	chain := blockchain.NewService(
		sb,
		nil, // blockchain.BlobProcessor
		nil, // deposit.Contract
		logger,
		cs,
		eng,
		fb,
		sp,
		ts,
		optimisticPayloadBuilds,
	)
	// Note: test avoid calling chain.Start since it only starts the deposits
	// goroutine which is not really relevant for this test

	// Before processing any block it is mandatory to handle genesis
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

	// Finally create a block that will be rejected and
	// verify the state on top of which is next payload built
	var (
		consensusTime   = math.U64(time.Now().Unix())
		proposerAddress = []byte{'d', 'u', 'm', 'm', 'y'} // this will err on purpose
	)

	// Since this is the first block called post genesis
	// forceSyncUponProcess will be called.
	// TODO: Make sure request matches expectations
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

	err = chain.VerifyIncomingBlock(
		ctx.ConsensusCtx(),
		invalidBlk,
		consensusTime,
		proposerAddress,
	)
	require.ErrorIs(t, err, core.ErrProposerMismatch)
}

var _ blockchain.LocalBuilder = (*fakeBuilder)(nil)

type fakeBuilder struct {
	enabled bool
}

func (fb *fakeBuilder) Enabled() bool { return fb.enabled }

func (fb *fakeBuilder) RequestPayloadAsync(
	_ context.Context,
	_ builder.ReadOnlyBeaconState,
	_ math.U64,
	_ common.ExecutionHash,
	_ common.ExecutionHash,
) (*engineprimitives.PayloadID, common.Version, error) {
	return nil, version.Electra(), errors.New("NOT IMPLEMENTED YET")
}
