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

package logs_test

import (
	"context"
	"strconv"
	"testing"

	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/beacon/staking/logs/mocks"
	"github.com/itsdevbear/bolaris/types/consensus"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/stretchr/testify/require"
)

func FuzzHandlerSimple(f *testing.F) {
	// Setup
	ctx := context.Background()
	stakingService := &mocks.StakingService{}
	callbackHandler, err := mocks.NewCallbackHandler(stakingService)
	require.NoError(f, err)

	events, err := mocks.DepositContractEvents()
	require.NoError(f, err)

	depositEventName := "Deposit"
	depositEvent := events[depositEventName]

	// withdrawalEventName := "Withdrawal"
	// withdrawalEvent := events[withdrawalEventName]

	f.Add([]byte("pubkey"), uint64(1))
	f.Fuzz(func(t *testing.T, pubKey []byte, amount uint64) {
		var (
			deposit       *consensusv1.Deposit
			latestDeposit *consensusv1.Deposit
			log           *coretypes.Log
			// We don't fuzz withdrawalCredentials because it's a fixed length
			// that prevents us from generating a variety of inputs.
			withdrawalCredentials = []byte("12345678901234567890")
		)

		// Deposit
		deposit = consensus.NewDeposit(pubKey, amount, withdrawalCredentials)
		log, err = mocks.NewLogFromDeposit(depositEvent, deposit)
		require.NoError(t, err)

		err = callbackHandler.HandleLog(ctx, log)
		require.NoError(t, err)

		latestDeposit, err = stakingService.MostRecentDeposit()
		require.NoError(t, err)
		require.Equal(t, deposit, latestDeposit)
		require.Equal(t, deposit.GetAmount(), latestDeposit.GetAmount())
		require.Equal(
			t,
			deposit.GetValidatorPubkey(),
			latestDeposit.GetValidatorPubkey(),
		)
		require.Equal(t,
			deposit.GetWithdrawalCredentials(),
			latestDeposit.GetWithdrawalCredentials())

		// err = stakingService.PersistDeposits(ctx)
		// require.NoError(t, err)
	})
}

func FuzzHandlerMulti(f *testing.F) {
	// Setup
	ctx := context.Background()
	stakingService := &mocks.StakingService{}
	callbackHandler, err := mocks.NewCallbackHandler(stakingService)
	require.NoError(f, err)

	events, err := mocks.DepositContractEvents()
	require.NoError(f, err)

	depositEventName := "Deposit"
	depositEvent := events[depositEventName]

	// withdrawalEventName := "Withdrawal"
	// withdrawalEvent := events[withdrawalEventName]

	f.Add(uint64(100), []byte("pubkey"), uint64(1))
	f.Fuzz(
		func(t *testing.T, nDeposits uint64, seekPubKey []byte, initAmount uint64) {
			var (
				deposit       *consensusv1.Deposit
				latestDeposit *consensusv1.Deposit
				log           *coretypes.Log
				// We don't fuzz withdrawalCredentials because it's a fixed
				// length
				// that prevents us from generating a variety of inputs.
				withdrawalCredentials = []byte("12345678901234567890")
			)

			for i := uint64(0); i < nDeposits; i++ {
				i := i
				// Deposit
				var pubKey []byte
				pubKey = append(pubKey, seekPubKey...)
				pubKey = append(pubKey, []byte(strconv.Itoa(int(i)))...)
				deposit = consensus.NewDeposit(
					pubKey,
					initAmount+i,
					withdrawalCredentials,
				)
				require.Equal(t, initAmount+i, deposit.GetAmount())
				require.Equal(t, pubKey, deposit.GetValidatorPubkey())

				log, err = mocks.NewLogFromDeposit(depositEvent, deposit)
				require.NoError(t, err)
				err = callbackHandler.HandleLog(ctx, log)
				require.NoError(t, err)

				latestDeposit, err = stakingService.MostRecentDeposit()
				require.NoError(t, err)
				require.Equal(t, deposit, latestDeposit)

				require.Equal(t, int(i+1), stakingService.NumPendingDeposits())
			}
			err = stakingService.ApplyDeposits(ctx)
			require.NoError(t, err)
		},
	)
}
