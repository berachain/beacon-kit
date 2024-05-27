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

package state

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type KVStore[
	KVStoreT any,
	ForkT any,
	BeaconBlockHeaderT any,
	Eth1DataT any,
	ValidatorT any,
] interface {
	Context() context.Context
	WithContext(
		ctx context.Context,
	) KVStoreT
	Save()
	GetLatestExecutionPayloadHeader() (
		engineprimitives.ExecutionPayloadHeader, error,
	)
	SetLatestExecutionPayloadHeader(
		payloadHeader engineprimitives.ExecutionPayloadHeader,
	) error
	GetEth1DepositIndex() (uint64, error)
	SetEth1DepositIndex(
		index uint64,
	) error
	GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
	SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
	Copy() KVStoreT
	GetSlot() (math.Slot, error)
	SetSlot(slot math.Slot) error
	GetFork() (ForkT, error)
	SetFork(fork ForkT) error
	GetGenesisValidatorsRoot() (common.Root, error)
	SetGenesisValidatorsRoot(root common.Root) error
	GetLatestBlockHeader() (BeaconBlockHeaderT, error)
	SetLatestBlockHeader(header BeaconBlockHeaderT) error
	GetBlockRootAtIndex(index uint64) (primitives.Root, error)
	StateRootAtIndex(index uint64) (primitives.Root, error)
	GetEth1Data() (Eth1DataT, error)
	SetEth1Data(data Eth1DataT) error
	GetValidators() ([]ValidatorT, error)
	GetBalances() ([]uint64, error)
	GetNextWithdrawalIndex() (uint64, error)
	SetNextWithdrawalIndex(index uint64) error
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	SetNextWithdrawalValidatorIndex(index math.ValidatorIndex) error
	GetTotalSlashing() (math.Gwei, error)
	SetTotalSlashing(total math.Gwei) error
	GetRandaoMixAtIndex(index uint64) (primitives.Bytes32, error)
	GetSlashings() ([]uint64, error)
	SetSlashingAtIndex(index uint64, amount math.Gwei) error
	GetSlashingAtIndex(index uint64) (math.Gwei, error)
	GetTotalValidators() (uint64, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	ValidatorByIndex(index math.ValidatorIndex) (ValidatorT, error)
	UpdateBlockRootAtIndex(index uint64, root primitives.Root) error
	UpdateStateRootAtIndex(index uint64, root primitives.Root) error
	UpdateRandaoMixAtIndex(index uint64, mix primitives.Bytes32) error
	UpdateValidatorAtIndex(
		index math.ValidatorIndex,
		validator ValidatorT,
	) error
	ValidatorIndexByPubkey(pubkey crypto.BLSPubkey) (math.ValidatorIndex, error)
	AddValidator(
		val ValidatorT,
	) error
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	RemoveValidatorAtIndex(idx math.ValidatorIndex) error
}

type Validator interface {
	GetWithdrawalCredentials() types.WithdrawalCredentials
	IsFullyWithdrawable(math.Gwei, math.Epoch) bool
	IsPartiallyWithdrawable(math.Gwei, math.Gwei) bool
}
