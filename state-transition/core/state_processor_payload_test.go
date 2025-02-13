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

// // TestPayloadTimestampVerification ensures that payload timestamp
// // is properly validated
// func TestPayloadTimestampVerification(t *testing.T) {
// 	// Create state processor to test
// 	cs := setupChain(t)
// 	sp, st, ds, ctx, cms, mockEngine := statetransition.SetupTestState(t, cs)

// 	// process genesis before any other block
// 	genesisTime := time.Now().Truncate(time.Second)
// 	var (
// 		genDeposits = types.Deposits{
// 			{
// 				Pubkey:      [48]byte{0x00},
// 				Credentials: types.NewCredentialsFromExecutionAddress(common.ExecutionAddress{}),
// 				Amount:      math.Gwei(cs.MaxEffectiveBalance()),
// 				Index:       0,
// 			},
// 		}
// 		genPayloadHeader = new(types.ExecutionPayloadHeader).Empty()
// 		genVersion       = version.Deneb()
// 	)
// 	genPayloadHeader.Timestamp = math.U64(genesisTime.Unix())

// 	_, err := sp.InitializePreminedBeaconStateFromEth1(
// 		st, genDeposits, genPayloadHeader, genVersion,
// 	)
// 	require.NoError(t, err)
// 	require.NoError(t, ds.EnqueueDeposits(ctx.ConsensusCtx(), genDeposits))

// 	// Test cases
// 	consensusBlkTime := genesisTime.Add(time.Second)
// 	tests := []struct {
// 		name        string
// 		setupMocksF func()
// 		payloadTime time.Time
// 		expectedErr error
// 	}{
// 		{
// 			name: "Payload timestamp < consensus timestamp",
// 			setupMocksF: func() {
// 				// we don't really verify payloads here, so just mock a positive result
// 				mockEngine.EXPECT().VerifyAndNotifyNewPayload(mock.Anything, mock.Anything).Return(nil)
// 			},
// 			payloadTime: consensusBlkTime.Add(-10 * time.Second),
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "Payload timestamp == consensus timestamp",
// 			setupMocksF: func() {
// 				// we don't really verify payloads here, so just mock a positive result
// 				mockEngine.EXPECT().VerifyAndNotifyNewPayload(mock.Anything, mock.Anything).Return(nil)
// 			},
// 			payloadTime: consensusBlkTime,
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "Payload timestamp > consensus timestamp",
// 			setupMocksF: func() {
// 				// we don't really verify payloads here, so just mock a positive result
// 				mockEngine.EXPECT().VerifyAndNotifyNewPayload(mock.Anything, mock.Anything).Return(nil)
// 			},
// 			payloadTime: consensusBlkTime.Add(time.Second),
// 			expectedErr: nil,
// 		},
// 		{
// 			name: "Payload timestamp >> consensus timestamp",
// 			setupMocksF: func() {
// 				// no mock here, since timestamp validation fails
// 			},
// 			payloadTime: consensusBlkTime.Add(2 * time.Second),
// 			expectedErr: payloadtime.ErrTooFarInTheFuture,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// these test cases should not be run in parallel
// 			// since state processors is shared by them
// 			tt.setupMocksF()

// 			// create independent states per each test
// 			sdkCtx := statetransition.NewSDKContext(cms)
// 			testSt := st.Copy(sdkCtx)

// 			// DEBUG CHECKS. Since we copy testSt Slot must be equal to st's one, i.e. zero
// 			// There is currently a bug in the test where this happens only for the fist test
// 			// while subsequent ones show increasing slots, as if we don't properly isolate tests
// 			// (which on the other hand works all right in isolation)
// 			checkSlot, err := testSt.GetSlot()
// 			require.NoError(t, err)
// 			require.Equal(t, math.Slot(0), checkSlot)

// 			ctx = transition.NewTransitionCtx(
// 				sdkCtx.Context(),
// 				math.U64(consensusBlkTime.Unix()),
// 				statetransition.DummyProposerAddr,
// 			).
// 				WithVerifyPayload(true).
// 				WithVerifyRandao(false).
// 				WithVerifyResult(false).
// 				WithMeterGas(false)

// 			blk := buildNextBlock(
// 				t,
// 				testSt,
// 				&types.BeaconBlockBody{
// 					ExecutionPayload: testPayload(
// 						math.U64(tt.payloadTime.Unix()),
// 						testSt.EVMInflationWithdrawal(constants.GenesisSlot+1),
// 					),
// 					Eth1Data: types.NewEth1Data(genDeposits.HashTreeRoot()),
// 				},
// 			)

// 			_, err = sp.Transition(ctx, testSt, blk)
// 			if tt.expectedErr == nil {
// 				require.NoError(t, err)
// 			} else {
// 				require.ErrorIs(t, err, tt.expectedErr)
// 			}
// 		})
// 	}
// }
