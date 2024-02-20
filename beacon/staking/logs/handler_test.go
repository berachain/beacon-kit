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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
	stakingabi "github.com/itsdevbear/bolaris/beacon/staking/abi"
	"github.com/itsdevbear/bolaris/beacon/staking/logs"
	"github.com/itsdevbear/bolaris/runtime/service"
	"github.com/itsdevbear/bolaris/types/consensus"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

func Test_CallbackHandler(t *testing.T) {
	// Setup
	ctx := context.Background()
	stakingService := &mockStakingService{}
	callbackHandler, err := newCallbackHandler(stakingService)
	require.NoError(t, err)

	events, err := depositContractEvents()
	require.NoError(t, err)

	depositEventName := "Deposit"
	depositEvent := events[depositEventName]

	withdrawalEventName := "Withdrawal"
	withdrawalEvent := events[withdrawalEventName]

	t.Run("should add correct deposits and withdrawals into staking service", func(t *testing.T) {
		deposit := consensus.NewDeposit(
			[]byte("pubkey"),
			10000,
			[]byte("12345678901234567890"),
		)
		log, err := newLogFromDeposit(depositEvent, deposit)
		require.NoError(t, err)

		err = callbackHandler.HandleLog(ctx, &log)
		require.NoError(t, err)

		latestDeposit, err := stakingService.mostRecentDeposit()
		require.NoError(t, err)
		require.Equal(t, deposit, latestDeposit)

		withdrawal := engine.NewWithdrawal(
			[]byte("pubkey"),
			10000,
		)
		log, err = newLogFromWithdrawal(withdrawalEvent, withdrawal)
		require.NoError(t, err)

		err = callbackHandler.HandleLog(ctx, &log)
		require.NoError(t, err)

		latestWithdrawal, err := stakingService.mostRecentWithdrawal()
		require.NoError(t, err)
		require.Equal(t, withdrawal, latestWithdrawal)
	})
}

func newLogFromDeposit(event abi.Event, deposit *consensusv1.Deposit) (coretypes.Log, error) {
	data, error := event.Inputs.Pack(
		deposit.Pubkey,
		[20]byte(deposit.WithdrawalCredentials),
		deposit.Amount,
	)
	if error != nil {
		return coretypes.Log{}, error
	}
	return coretypes.Log{
		Topics: []common.Hash{event.ID},
		Data:   data,
	}, nil
}

func newLogFromWithdrawal(event abi.Event, withdrawal *enginev1.Withdrawal) (coretypes.Log, error) {
	data, error := event.Inputs.Pack(
		[]byte{},
		[20]byte{},
		withdrawal.Amount,
	)
	if error != nil {
		return coretypes.Log{}, error
	}
	return coretypes.Log{
		Topics: []common.Hash{event.ID},
		Data:   data,
	}, nil
}

func depositContractEvents() (map[string]abi.Event, error) {
	stakingAbi, err := stakingabi.StakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return stakingAbi.Events, nil
}

func newCallbackHandler(stakingService logs.StakingService) (callback.LogHandler, error) {
	logHander := service.New[logs.Handler](
		logs.WithStakingService(stakingService),
	)
	callbackHandler, err := callback.NewFrom(logHander)
	if err != nil {
		return nil, err
	}
	return callbackHandler, nil
}
