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
	"reflect"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	beacontypes "github.com/itsdevbear/bolaris/beacon/core/types"
	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	loghandler "github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/beacon/staking/logs"
	"github.com/itsdevbear/bolaris/beacon/staking/logs/mocks"
	"github.com/itsdevbear/bolaris/contracts/abi"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/stretchr/testify/require"
)

func TestLogFactory(t *testing.T) {
	contractAddress := ethcommon.HexToAddress("0x1234")
	stakingAbi, err := abi.StakingMetaData.GetAbi()
	require.NoError(t, err)

	stakingLogRequest, err := logs.NewStakingRequest(
		contractAddress,
	)
	require.NoError(t, err)
	factory, err := loghandler.NewFactory(
		loghandler.WithRequest(stakingLogRequest),
	)
	require.NoError(t, err)

	deposit := beacontypes.NewDeposit(
		[]byte("pubkey"),
		10000,
		[]byte("12345678901234567890"),
	)
	log, err := mocks.NewLogFromDeposit(
		stakingAbi.Events[logs.DepositName],
		deposit,
	)
	require.NoError(t, err)
	log.Address = contractAddress

	val, err := factory.UnmarshalEthLog(log)
	require.NoError(t, err)

	valType := reflect.TypeOf(val.Interface())
	require.NotNil(t, valType)
	require.Equal(t, reflect.Ptr, valType.Kind())
	require.Equal(t, logs.DepositType, valType.Elem())

	newDeposit, ok := val.Interface().(*beacontypesv1.Deposit)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, deposit, newDeposit)

	withdrawal := enginetypes.NewWithdrawal([]byte("pubkey"), 10000)
	log, err = mocks.NewLogFromWithdrawal(
		stakingAbi.Events[logs.WithdrawalName],
		withdrawal,
	)
	require.NoError(t, err)
	log.Address = contractAddress

	_, err = factory.UnmarshalEthLog(log)
	// An error is expected because the event type in ABI and
	// withdrawalType are mismatched,
	// (no validatorPubkey in withdrawalType currently).
	require.Error(t, err)
}
