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

package state

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// StateDB is the underlying struct behind the BeaconState interface.
//
//nolint:revive // todo fix somehow
type StateDB[
	BeaconBlockHeaderT any,
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	],
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	StateStoreT StateStore[
		StateStoreT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	StateStore[
		StateStoreT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	]
	cs common.ChainSpec
}

// NewFromDB creates a new beacon state from an underlying state db.
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) NewFromDB(
	bdb StateStoreT,
	cs common.ChainSpec,
) *StateDB[
	BeaconBlockHeaderT,
	BeaconStateMarshallableT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	StateStoreT,
	ValidatorT,
	WithdrawalCredentialsT,
] {
	result := &StateDB[
		BeaconBlockHeaderT,
		BeaconStateMarshallableT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		StateStoreT,
		ValidatorT,
		WithdrawalCredentialsT,
	]{
		StateStore: bdb,
		cs:         cs,
	}

	// TODO: Fix this is hood as fuck.
	// return reflect.ValueOf(result).Interface().(BeaconStateT)
	return result
}

// Copy returns a copy of the beacon state.
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) Copy() *StateDB[
	BeaconBlockHeaderT,
	BeaconStateMarshallableT,
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT,
	StateStoreT,
	ValidatorT,
	WithdrawalCredentialsT,
] {
	return s.NewFromDB(
		s.StateStore.Copy(),
		s.cs,
	)
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) IncreaseBalance(
	idx math.ValidatorIndex,
	delta math.Gwei,
) error {
	balance, err := s.GetBalance(idx)
	if err != nil {
		return err
	}
	return s.SetBalance(idx, balance+delta)
}

// DecreaseBalance decreases the balance of a validator.
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) DecreaseBalance(
	idx math.ValidatorIndex,
	delta math.Gwei,
) error {
	balance, err := s.GetBalance(idx)
	if err != nil {
		return err
	}
	return s.SetBalance(idx, balance-min(balance, delta))
}

// UpdateSlashingAtIndex sets the slashing amount in the store.
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) UpdateSlashingAtIndex(
	index uint64,
	amount math.Gwei,
) error {
	// Update the total slashing amount before overwriting the old amount.
	total, err := s.GetTotalSlashing()
	if err != nil {
		return err
	}

	oldValue, err := s.GetSlashingAtIndex(index)
	if err != nil {
		return err
	}

	// Defensive check but total - oldValue should never underflow.
	if oldValue > total {
		return errors.New("count of total slashing is not up to date")
	} else if err = s.SetTotalSlashing(
		total - oldValue + amount,
	); err != nil {
		return err
	}

	return s.SetSlashingAtIndex(index, amount)
}

// ExpectedWithdrawals as defined in the Ethereum 2.0 Specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-get_expected_withdrawals
//
//nolint:lll
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) ExpectedWithdrawals() ([]*engineprimitives.Withdrawal, error) {
	var (
		validator         ValidatorT
		balance           math.Gwei
		withdrawalAddress common.ExecutionAddress
		withdrawals       = make([]*engineprimitives.Withdrawal, 0)
	)

	slot, err := s.GetSlot()
	if err != nil {
		return nil, err
	}

	epoch := math.Epoch(uint64(slot) / s.cs.SlotsPerEpoch())

	withdrawalIndex, err := s.GetNextWithdrawalIndex()
	if err != nil {
		return nil, err
	}

	validatorIndex, err := s.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return nil, err
	}

	totalValidators, err := s.GetTotalValidators()
	if err != nil {
		return nil, err
	}

	// Iterate through indices to find the next validators to withdraw.
	for range min(
		s.cs.MaxValidatorsPerWithdrawalsSweep(), totalValidators,
	) {
		validator, err = s.ValidatorByIndex(validatorIndex)
		if err != nil {
			return nil, err
		}

		balance, err = s.GetBalance(validatorIndex)
		if err != nil {
			return nil, err
		}

		withdrawalAddress, err = validator.
			GetWithdrawalCredentials().ToExecutionAddress()
		if err != nil {
			return nil, err
		}

		// These fields are the same for both partial and full withdrawals.
		withdrawal := &engineprimitives.Withdrawal{
			Index:     math.U64(withdrawalIndex),
			Validator: validatorIndex,
			Address:   withdrawalAddress,
		}

		// Set the amount of the withdrawal depending on the balance of the
		// validator.
		if validator.IsFullyWithdrawable(balance, epoch) {
			withdrawal.Amount = balance
		} else if validator.IsPartiallyWithdrawable(
			balance, math.Gwei(s.cs.MaxEffectiveBalance()),
		) {
			withdrawal.Amount = balance - math.Gwei(s.cs.MaxEffectiveBalance())
		}
		withdrawals = append(withdrawals, withdrawal)

		// Increment the withdrawal index to process the next withdrawal.
		withdrawalIndex++

		// Cap the number of withdrawals to the maximum allowed per payload.
		//#nosec:G701 // won't overflow in practice.
		if len(withdrawals) == int(s.cs.MaxWithdrawalsPerPayload()) {
			break
		}

		// Increment the validator index to process the next validator.
		validatorIndex = (validatorIndex + 1) % math.ValidatorIndex(
			totalValidators,
		)
	}

	return withdrawals, nil
}

// HashTreeRoot is the interface for the beacon store.
//
//nolint:funlen,gocognit // todo fix somehow
func (s *StateDB[
	BeaconBlockHeaderT, BeaconStateMarshallableT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, StateStoreT,
	ValidatorT, WithdrawalCredentialsT,
]) HashTreeRoot() ([32]byte, error) {
	slot, err := s.GetSlot()
	if err != nil {
		return [32]byte{}, err
	}

	fork, err := s.GetFork()
	if err != nil {
		return [32]byte{}, err
	}

	genesisValidatorsRoot, err := s.GetGenesisValidatorsRoot()
	if err != nil {
		return [32]byte{}, err
	}

	latestBlockHeader, err := s.GetLatestBlockHeader()
	if err != nil {
		return [32]byte{}, err
	}

	blockRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		blockRoots[i], err = s.GetBlockRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
	}

	stateRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		stateRoots[i], err = s.StateRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
	}

	latestExecutionPayloadHeader, err := s.GetLatestExecutionPayloadHeader()
	if err != nil {
		return [32]byte{}, err
	}

	eth1Data, err := s.GetEth1Data()
	if err != nil {
		return [32]byte{}, err
	}

	eth1DepositIndex, err := s.GetEth1DepositIndex()
	if err != nil {
		return [32]byte{}, err
	}

	validators, err := s.GetValidators()
	if err != nil {
		return [32]byte{}, err
	}

	balances, err := s.GetBalances()
	if err != nil {
		return [32]byte{}, err
	}

	randaoMixes := make([]common.Bytes32, s.cs.EpochsPerHistoricalVector())
	for i := range s.cs.EpochsPerHistoricalVector() {
		randaoMixes[i], err = s.GetRandaoMixAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
	}

	nextWithdrawalIndex, err := s.GetNextWithdrawalIndex()
	if err != nil {
		return [32]byte{}, err
	}

	nextWithdrawalValidatorIndex, err := s.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return [32]byte{}, err
	}

	slashings, err := s.GetSlashings()
	if err != nil {
		return [32]byte{}, err
	}

	totalSlashings, err := s.GetTotalSlashing()
	if err != nil {
		return [32]byte{}, err
	}

	// TODO: Properly move BeaconState into full generics.
	st, err := (*new(BeaconStateMarshallableT)).New(
		s.cs.ActiveForkVersionForSlot(slot),
		genesisValidatorsRoot,
		slot,
		fork,
		latestBlockHeader,
		blockRoots,
		stateRoots,
		eth1Data,
		eth1DepositIndex,
		latestExecutionPayloadHeader,
		validators,
		balances,
		randaoMixes,
		nextWithdrawalIndex,
		nextWithdrawalValidatorIndex,
		slashings,
		totalSlashings,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return st.HashTreeRoot()
}
