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

package storage

import (
	"context"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// Deposit is a struct that represents a deposit.
type Deposit interface {
	constraints.SSZMarshallable
	GetIndex() uint64
}

// DepositStore defines the interface for deposit storage.
type DepositStore[DepositT any] interface {
	// GetDepositsByIndex returns `numView` expected deposits.
	GetDepositsByIndex(
		startIndex uint64,
		numView uint64,
	) ([]DepositT, error)
	// Prune prunes the deposit store of [start, end)
	Prune(start, end uint64) error
	// EnqueueDeposits adds a list of deposits to the deposit store.
	EnqueueDeposits(deposits []DepositT) error
}

// KVStore is the interface for the key-value store holding the beacon state.
type KVStore[
	T,
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	ValidatorT any,
] interface {
	// Context returns the context of the key-value store.
	Context() context.Context
	// WithContext returns a new key-value store with the given context.
	WithContext(
		ctx context.Context,
	) T
	// Save saves the key-value store.
	Save()
	// Copy returns a copy of the key-value store.
	Copy() T
	// GetLatestExecutionPayloadHeader retrieves the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
	// SetLatestExecutionPayloadHeader sets the latest execution payload header.
	SetLatestExecutionPayloadHeader(
		payloadHeader ExecutionPayloadHeaderT,
	) error
	// GetEth1DepositIndex retrieves the eth1 deposit index.
	GetEth1DepositIndex() (uint64, error)
	// SetEth1DepositIndex sets the eth1 deposit index.
	SetEth1DepositIndex(
		index uint64,
	) error
	// GetBalance retrieves the balance of a validator.
	GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
	// SetBalance sets the balance of a validator.
	SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
	// GetSlot retrieves the current slot.
	GetSlot() (math.Slot, error)
	// SetSlot sets the current slot.
	SetSlot(slot math.Slot) error
	// GetFork retrieves the fork.
	GetFork() (ForkT, error)
	// SetFork sets the fork.
	SetFork(fork ForkT) error
	// GetGenesisValidatorsRoot retrieves the genesis validators root.
	GetGenesisValidatorsRoot() (common.Root, error)
	// SetGenesisValidatorsRoot sets the genesis validators root.
	SetGenesisValidatorsRoot(root common.Root) error
	// GetLatestBlockHeader retrieves the latest block header.
	GetLatestBlockHeader() (BeaconBlockHeaderT, error)
	// SetLatestBlockHeader sets the latest block header.
	SetLatestBlockHeader(header BeaconBlockHeaderT) error
	// GetBlockRootAtIndex retrieves the block root at the given index.
	GetBlockRootAtIndex(index uint64) (common.Root, error)
	// StateRootAtIndex retrieves the state root at the given index.
	StateRootAtIndex(index uint64) (common.Root, error)
	// GetEth1Data retrieves the eth1 data.
	GetEth1Data() (Eth1DataT, error)
	// SetEth1Data sets the eth1 data.
	SetEth1Data(data Eth1DataT) error
	// GetValidators retrieves all validators.
	GetValidators() ([]ValidatorT, error)
	// GetBalances retrieves all balances.
	GetBalances() ([]uint64, error)
	// GetNextWithdrawalIndex retrieves the next withdrawal index.
	GetNextWithdrawalIndex() (uint64, error)
	// SetNextWithdrawalIndex sets the next withdrawal index.
	SetNextWithdrawalIndex(index uint64) error
	// GetNextWithdrawalValidatorIndex retrieves the next withdrawal validator
	// index.
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	// SetNextWithdrawalValidatorIndex sets the next withdrawal validator index.
	SetNextWithdrawalValidatorIndex(index math.ValidatorIndex) error
	// GetTotalSlashing retrieves the total slashing.
	GetTotalSlashing() (math.Gwei, error)
	// SetTotalSlashing sets the total slashing.
	SetTotalSlashing(total math.Gwei) error
	// GetRandaoMixAtIndex retrieves the randao mix at the given index.
	GetRandaoMixAtIndex(index uint64) (common.Bytes32, error)
	// GetSlashings retrieves all slashings.
	GetSlashings() ([]uint64, error)
	// SetSlashingAtIndex sets the slashing at the given index.
	SetSlashingAtIndex(index uint64, amount math.Gwei) error
	// GetSlashingAtIndex retrieves the slashing at the given index.
	GetSlashingAtIndex(index uint64) (math.Gwei, error)
	// GetTotalValidators retrieves the total validators.
	GetTotalValidators() (uint64, error)
	// GetTotalActiveBalances retrieves the total active balances.
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	// ValidatorByIndex retrieves the validator at the given index.
	ValidatorByIndex(index math.ValidatorIndex) (ValidatorT, error)
	// UpdateBlockRootAtIndex updates the block root at the given index.
	UpdateBlockRootAtIndex(index uint64, root common.Root) error
	// UpdateStateRootAtIndex updates the state root at the given index.
	UpdateStateRootAtIndex(index uint64, root common.Root) error
	// UpdateRandaoMixAtIndex updates the randao mix at the given index.
	UpdateRandaoMixAtIndex(index uint64, mix common.Bytes32) error
	// UpdateValidatorAtIndex updates the validator at the given index.
	UpdateValidatorAtIndex(
		index math.ValidatorIndex,
		validator ValidatorT,
	) error
	// ValidatorIndexByPubkey retrieves the validator index by the given pubkey.
	ValidatorIndexByPubkey(pubkey crypto.BLSPubkey) (math.ValidatorIndex, error)
	// AddValidator adds a validator.
	AddValidator(val ValidatorT) error
	// AddValidatorBartio adds a validator to the Bartio chain.
	AddValidatorBartio(val ValidatorT) error
	// ValidatorIndexByCometBFTAddress retrieves the validator index by the
	// given comet BFT address.
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
	// GetValidatorsByEffectiveBalance retrieves validators by effective
	// balance.
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	// RemoveValidatorAtIndex removes the validator at the given index.
	RemoveValidatorAtIndex(idx math.ValidatorIndex) error
}

// Validator represents an interface for a validator with generic withdrawal
// credentials. WithdrawalCredentialsT is a type parameter that must implement
// the WithdrawalCredentials interface.
type Validator[WithdrawalCredentialsT WithdrawalCredentials] interface {
	// GetWithdrawalCredentials returns the withdrawal credentials of the
	// validator.
	GetWithdrawalCredentials() WithdrawalCredentialsT
	// IsFullyWithdrawable checks if the validator is fully withdrawable given a
	// certain Gwei amount and epoch.
	IsFullyWithdrawable(amount math.Gwei, epoch math.Epoch) bool
	// IsPartiallyWithdrawable checks if the validator is partially withdrawable
	// given two Gwei amounts.
	IsPartiallyWithdrawable(amount1 math.Gwei, amount2 math.Gwei) bool
}

// Withdrawal represents an interface for a withdrawal.
type Withdrawal[T any] interface {
	New(
		index math.U64,
		validator math.ValidatorIndex,
		address gethprimitives.ExecutionAddress,
		amount math.Gwei,
	) T
}

// WithdrawalCredentials represents an interface for withdrawal credentials.
type WithdrawalCredentials interface {
	// ToExecutionAddress converts the withdrawal credentials to an execution
	// address.
	ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
}
