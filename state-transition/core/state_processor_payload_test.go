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
	"testing"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	statetransition "github.com/berachain/beacon-kit/testing/state-transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestPayloadTimestampVerification ensures that payload timestamp
// is properly validated
//
//nolint:paralleltest // uses envars
func TestPayloadTimestampVerification(t *testing.T) {
	// Create state processor to test
	cs := setupChain(t)
	sp, st, ds, ctx, cms, mockEngine := statetransition.SetupTestState(t, cs)

	// process genesis before any other block
	genesisTime := time.Now().Truncate(time.Second)
	genesisFork := cs.ActiveForkVersionForTimestamp(math.U64(genesisTime.Unix()))
	require.Equal(t, genesisFork, cs.GenesisForkVersion())
	var (
		genDeposits = types.Deposits{
			{
				Pubkey:      [48]byte{0x00},
				Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{}),
				Amount:      math.Gwei(cs.MaxEffectiveBalance()),
				Index:       0,
			},
		}
		genPayloadHeader = &types.ExecutionPayloadHeader{
			Versionable: types.NewVersionable(genesisFork),
		}
	)
	genPayloadHeader.Timestamp = math.U64(genesisTime.Unix())

	_, err := sp.InitializeBeaconStateFromEth1(st, genDeposits, genPayloadHeader, genesisFork)
	require.NoError(t, err)
	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

	// write genesis changes to make them available for next blocks
	//nolint:errcheck // false positive as this has no return value
	ctx.ConsensusCtx().(sdk.Context).MultiStore().(storetypes.CacheMultiStore).Write()

	// Test cases
	consensusBlkTime := genesisTime.Add(time.Second)
	tests := []struct {
		name        string
		setupMocksF func()
		payloadTime time.Time
		expectedErr error
	}{
		{
			name: "Payload timestamp < consensus timestamp",
			setupMocksF: func() {
				mockEngine.EXPECT().NotifyNewPayload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			payloadTime: consensusBlkTime.Add(-10 * time.Second),
			expectedErr: nil,
		},
		{
			name: "Payload timestamp == consensus timestamp",
			setupMocksF: func() {
				mockEngine.EXPECT().NotifyNewPayload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			payloadTime: consensusBlkTime,
			expectedErr: nil,
		},
		{
			name: "Payload timestamp > consensus timestamp",
			setupMocksF: func() {
				mockEngine.EXPECT().NotifyNewPayload(mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			payloadTime: consensusBlkTime.Add(time.Second),
			expectedErr: nil,
		},
		{
			name: "Payload timestamp >> consensus timestamp",
			setupMocksF: func() {
				// no mock here, since timestamp validation fails
			},
			payloadTime: consensusBlkTime.Add(2 * time.Second),
			expectedErr: payloadtime.ErrTooFarInTheFuture,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// these test cases should not be run in parallel
			// since state processors is shared by them
			tt.setupMocksF()

			// create independent states per each test
			sdkCtx := sdk.NewContext(cms.CacheMultiStore(), true, log.NewNopLogger())
			testSt := statedb.NewBeaconStateFromDB(st.KVStore.WithContext(sdkCtx), cs)

			tCtx := transition.NewTransitionCtx(
				sdkCtx,
				math.U64(consensusBlkTime.Unix()),
				statetransition.DummyProposerAddr,
			).
				WithVerifyPayload(true).
				WithVerifyRandao(false).
				WithVerifyResult(false).
				WithMeterGas(false)

			blk := buildNextBlock(
				t,
				cs,
				testSt,
				types.NewEth1Data(genDeposits.HashTreeRoot()),
				math.U64(tt.payloadTime.Unix()),
				nil,
				testSt.EVMInflationWithdrawal(math.U64(tt.payloadTime.Unix())),
			)

			_, err = sp.Transition(tCtx, testSt, blk)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}
