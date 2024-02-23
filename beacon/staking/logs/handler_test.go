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
	"testing"

	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/itsdevbear/bolaris/beacon/staking/logs/mocks"
	"github.com/itsdevbear/bolaris/types/consensus"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
	"github.com/stretchr/testify/require"
)

func Test_CallbackHandler(t *testing.T) {
	// Setup
	ctx := context.Background()
	stakingService := &mocks.StakingService{}
	callbackHandler, err := mocks.NewCallbackHandler(stakingService)
	require.NoError(t, err)

	events, err := mocks.DepositContractEvents()
	require.NoError(t, err)

	depositEventName := "Deposit"
	depositEvent := events[depositEventName]

	withdrawalEventName := "Withdrawal"
	withdrawalEvent := events[withdrawalEventName]

	t.Run(
		"should add correct deposits and withdrawals into staking service",
		func(t *testing.T) {
			var (
				deposit          *consensusv1.Deposit
				latestDeposit    *consensusv1.Deposit
				withdrawal       *enginev1.Withdrawal
				latestWithdrawal *enginev1.Withdrawal
				log              *coretypes.Log
			)

			deposit = consensus.NewDeposit(
				[]byte("pubkey"),
				10000,
				[]byte("12345678901234567890"),
			)
			log, err = mocks.NewLogFromDeposit(depositEvent, deposit)
			require.NoError(t, err)

			err = callbackHandler.HandleLog(ctx, log)
			require.NoError(t, err)

			latestDeposit, err = stakingService.MostRecentDeposit()
			require.NoError(t, err)
			require.Equal(t, deposit, latestDeposit)

			withdrawal = engine.NewWithdrawal(
				[]byte("pubkey"),
				10000,
			)
			log, err = mocks.NewLogFromWithdrawal(withdrawalEvent, withdrawal)
			require.NoError(t, err)

			err = callbackHandler.HandleLog(ctx, log)
			require.NoError(t, err)

			latestWithdrawal, err = stakingService.MostRecentWithdrawal()
			require.NoError(t, err)
			require.Equal(t, withdrawal, latestWithdrawal)
		},
	)
}
