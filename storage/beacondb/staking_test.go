// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package beacondb_test

// type mockValidator struct {
// 	*math.U64
// }

// func (m *mockValidator) IsActive(_ math.Epoch) bool {
// 	// Assuming a simple active status check based on a condition
// 	// This is a mock implementation and should be replaced with actual logic
// 	return true
// }

// func (m *mockValidator) GetPubkey() crypto.BLSPubkey {
// 	// Return a mock BLS public key
// 	// This is a mock implementation and should be replaced with actual logic
// 	return crypto.BLSPubkey{}
// }

// func (m *mockValidator) GetEffectiveBalance() math.Gwei {
// 	// Return a mock effective balance
// 	// This is a mock implementation and should be replaced with actual logic
// 	return 1000000 // 1 million Gwei as a placeholder
// }

// func testFactory() *math.U64 {
// 	return (*math.U64)(nil)
// }

// func TestDeposits(t *testing.T) {
// testName := "test"
// logger := log.NewTestLogger(t)
// keys := storetypes.NewKVStoreKeys(testName)
// cms := integration.CreateMultiStore(keys, logger)
// ctx := sdk.NewContext(cms, true, logger)
// storeKey := keys[testName]

// sdb := beacondb.New[
// 	*math.U64, *math.U64, *math.U64, *math.U64, *mockValidator,
// ](
// 	sdkruntime.NewKVStoreService(storeKey),
// 	testFactory,
// )
// _ = sdb.WithContext(ctx)
// t.Run("should work with deposit", func(t *testing.T) {
// fakeDeposit := primitives.U64(69420)
// err := sdb.EnqueueDeposits([]*primitives.U64{&fakeDeposit})
// require.NoError(t, err)
// deposits, err := sdb.DequeueDeposits(1)
// require.NoError(t, err)
// require.Equal(t, fakeDeposit, *deposits[0])
// })
// }
