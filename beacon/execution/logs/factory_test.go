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
	"errors"
	"reflect"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	"github.com/itsdevbear/bolaris/contracts/abi"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/types/consensus"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	"github.com/stretchr/testify/require"
)

func TestLogFactory(t *testing.T) {
	contractAddress := ethcommon.HexToAddress("0x1234")
	stakingAbi, err := abi.StakingMetaData.GetAbi()
	require.NoError(t, err)

	depositSig := ethcrypto.Keccak256Hash(
		[]byte("Deposit(bytes,bytes,uint64)"),
	)
	depositName := "Deposit"
	depositType := reflect.TypeOf(consensusv1.Deposit{})

	withdrawalSig := ethcrypto.Keccak256Hash(
		[]byte("Withdrawal(bytes,bytes,uint64)"),
	)
	withdrawalName := "Withdrawal"
	withdrawalType := reflect.TypeOf(enginev1.Withdrawal{})

	allocator := logs.New[logs.TypeAllocator](
		logs.WithABI(stakingAbi),
		logs.WithNameAndType(depositSig, depositName, depositType),
		logs.WithNameAndType(withdrawalSig, withdrawalName, withdrawalType),
	)

	factory, err := logs.NewFactory(
		logs.WithTypeAllocator(contractAddress, allocator),
	)
	require.NoError(t, err)

	deposit := consensus.NewDeposit(
		[]byte("pubkey"),
		10000,
		[]byte("12345678901234567890"),
	)
	log, err := newLogFromDeposit(stakingAbi.Events[depositName], deposit)
	require.NoError(t, err)
	log.Address = contractAddress

	val, err := factory.UnmarshalEthLog(log)
	require.NoError(t, err)

	valType := reflect.TypeOf(val.Interface())
	require.NotNil(t, valType)
	require.Equal(t, reflect.Ptr, valType.Kind())
	require.Equal(t, depositType, valType.Elem())

	newDeposit, ok := val.Interface().(*consensusv1.Deposit)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, deposit, newDeposit)

	withdrawal := enginetypes.NewWithdrawal([]byte("pubkey"), 10000)
	log, err = newLogFromWithdrawal(
		stakingAbi.Events[withdrawalName],
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

// newLog creates a new log of an event from the given arguments.
func newLog(event ethabi.Event, args ...interface{}) (*coretypes.Log, error) {
	if len(event.Inputs) != len(args) {
		return nil, errors.New("mismatched number of arguments")
	}
	data, err := event.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return &coretypes.Log{
		Topics: []ethcommon.Hash{event.ID},
		Data:   data,
	}, nil
}

// NewLogFromDeposit creates a new log from the given deposit.
func newLogFromDeposit(
	event ethabi.Event,
	deposit *consensusv1.Deposit,
) (*coretypes.Log, error) {
	return newLog(event,
		deposit.GetValidatorPubkey(),
		deposit.GetWithdrawalCredentials(),
		deposit.GetAmount(),
	)
}

// newLogFromWithdrawal creates a new log from the given withdrawal.
func newLogFromWithdrawal(
	event ethabi.Event,
	withdrawal *enginev1.Withdrawal,
) (*coretypes.Log, error) {
	return newLog(event,
		[]byte{},
		[]byte{},
		withdrawal.GetAmount(),
	)
}
