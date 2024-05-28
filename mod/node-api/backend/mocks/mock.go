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

package mocks

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
}

func NewMockBackend() *backend.Backend {
	sdb := &StateDB{}
	b := backend.New(func(context.Context, string) backend.StateDB {
		return sdb
	})
	setReturnValues(sdb)
	return b
}

func setReturnValues(sdb *StateDB) {
	sdb.EXPECT().GetGenesisValidatorsRoot().Return(primitives.Root{0x01}, nil)
	sdb.EXPECT().GetSlot().Return(1, nil)
	sdb.EXPECT().GetLatestExecutionPayloadHeader().Return(nil, nil)
	sdb.EXPECT().SetLatestExecutionPayloadHeader(mock.Anything).Return(nil)
	sdb.EXPECT().GetEth1DepositIndex().Return(0, nil)
	sdb.EXPECT().SetEth1DepositIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetBalance(mock.Anything).Return(1, nil)
	sdb.EXPECT().SetBalance(mock.Anything, mock.Anything).Return(nil)
	sdb.EXPECT().SetSlot(mock.Anything).Return(nil)
	sdb.EXPECT().GetFork().Return(nil, nil)
	sdb.EXPECT().SetFork(mock.Anything).Return(nil)
	sdb.EXPECT().GetLatestBlockHeader().Return(nil, nil)
	sdb.EXPECT().SetLatestBlockHeader(mock.Anything).Return(nil)
	sdb.EXPECT().
		GetBlockRootAtIndex(mock.Anything).
		Return(primitives.Root{0x01}, nil)
	sdb.EXPECT().
		StateRootAtIndex(mock.Anything).
		Return(primitives.Root{0x01}, nil)
	sdb.EXPECT().GetEth1Data().Return(nil, nil)
	sdb.EXPECT().SetEth1Data(mock.Anything).Return(nil)
	sdb.EXPECT().GetValidators().Return(nil, nil)
	sdb.EXPECT().GetBalances().Return(nil, nil)
	sdb.EXPECT().GetNextWithdrawalIndex().Return(0, nil)
	sdb.EXPECT().SetNextWithdrawalIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetNextWithdrawalValidatorIndex().Return(0, nil)
	sdb.EXPECT().SetNextWithdrawalValidatorIndex(mock.Anything).Return(nil)
	sdb.EXPECT().GetTotalSlashing().Return(0, nil)
	sdb.EXPECT().SetTotalSlashing(mock.Anything).Return(nil)
	sdb.EXPECT().
		GetRandaoMixAtIndex(mock.Anything).
		Return(primitives.Bytes32{0x01}, nil)
	sdb.EXPECT().GetSlashings().Return(nil, nil)
	sdb.EXPECT().SetSlashingAtIndex(mock.Anything, mock.Anything).Return(nil)
	sdb.EXPECT().GetSlashingAtIndex(mock.Anything).Return(0, nil)
	sdb.EXPECT().GetTotalValidators().Return(0, nil)
	sdb.EXPECT().GetTotalActiveBalances(mock.Anything).Return(0, nil)
	sdb.EXPECT().ValidatorByIndex(mock.Anything).Return(&types.Validator{
		Pubkey:                     crypto.BLSPubkey{0x01},
		WithdrawalCredentials:      types.WithdrawalCredentials{0x01},
		EffectiveBalance:           0,
		Slashed:                    false,
		ActivationEligibilityEpoch: 0,
		ActivationEpoch:            0,
		ExitEpoch:                  0,
		WithdrawableEpoch:          0,
	}, nil)
	sdb.EXPECT().
		UpdateBlockRootAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateStateRootAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateRandaoMixAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().
		UpdateValidatorAtIndex(mock.Anything, mock.Anything).
		Return(nil)
	sdb.EXPECT().ValidatorIndexByPubkey(mock.Anything).Return(0, nil)
	sdb.EXPECT().AddValidator(mock.Anything).Return(nil)
	sdb.EXPECT().GetValidatorsByEffectiveBalance().Return(nil, nil)
}
