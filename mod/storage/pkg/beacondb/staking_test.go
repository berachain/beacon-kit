// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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
