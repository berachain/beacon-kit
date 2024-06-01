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
	"reflect"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// StateDB is the underlying struct behind the BeaconState interface.
//
//nolint:revive // todo fix somehow
type StateDB[
	BeaconStateT any,
	KVStoreT KVStore[
		KVStoreT,
		ForkT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ValidatorT,
	],
	ForkT any,
	BeaconBlockHeaderT any,
	Eth1DataT any,
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	KVStore[
		KVStoreT,
		ForkT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ValidatorT,
	]
	cs primitives.ChainSpec
}

// NewBeaconState creates a new beacon state from an underlying state db.
func NewBeaconStateFromDB[
	BeaconStateT any,
	KVStoreT KVStore[
		KVStoreT,
		ForkT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ValidatorT,
	],
	ForkT any,
	BeaconBlockHeaderT any,
	Eth1DataT any,
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalCredentialsT WithdrawalCredentials,
](
	bdb KVStore[
		KVStoreT,
		ForkT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ValidatorT,
	],
	cs primitives.ChainSpec,
) BeaconStateT {
	result := &StateDB[
		BeaconStateT,
		KVStoreT,
		ForkT,
		BeaconBlockHeaderT,
		Eth1DataT,
		ValidatorT,
		WithdrawalCredentialsT,
	]{
		KVStore: bdb,
		cs:      cs,
	}

	// TODO: Fix this is hood as fuck.
	return reflect.ValueOf(result).Interface().(BeaconStateT)
}

// Copy returns a copy of the beacon state.
func (s *StateDB[
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
]) Copy() BeaconStateT {
	return NewBeaconStateFromDB[BeaconStateT](
		s.KVStore.Copy(),
		s.cs,
	)
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB[
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
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
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
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
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
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
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
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

	// Iterate through indicies to find the next validators to withdraw.
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

// Store is the interface for the beacon store.
//
//nolint:funlen,gocognit // todo fix somehow
func (s *StateDB[
	BeaconStateT, KVStoreT, ForkT,
	BeaconBlockHeaderT, Eth1DataT, ValidatorT, WithdrawalCredentialsT,
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

	blockRoots := make([]primitives.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		blockRoots[i], err = s.GetBlockRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
	}

	stateRoots := make([]primitives.Root, s.cs.SlotsPerHistoricalRoot())
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

	randaoMixes := make([]primitives.Bytes32, s.cs.EpochsPerHistoricalVector())
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

	activeFork := s.cs.ActiveForkVersionForSlot(slot)
	switch activeFork {
	case version.Deneb:
		executionPayloadHeader, ok :=
			latestExecutionPayloadHeader.(*types.ExecutionPayloadHeaderDeneb)
		if !ok {
			return [32]byte{}, errors.New(
				"latest execution payload is not of type ExecutableDataDeneb")
		}

		// TODO: Use New() on BeaconState to prevent reflection usage.
		return (&deneb.BeaconState{
			Slot:                  slot,
			GenesisValidatorsRoot: genesisValidatorsRoot,
			Fork: reflect.ValueOf(fork).
				Interface().(*types.Fork),
			LatestBlockHeader: reflect.ValueOf(latestBlockHeader).
				Interface().(*types.BeaconBlockHeader),
			BlockRoots:                   blockRoots,
			StateRoots:                   stateRoots,
			LatestExecutionPayloadHeader: executionPayloadHeader,
			Eth1Data: reflect.ValueOf(eth1Data).
				Interface().(*types.Eth1Data),
			Eth1DepositIndex: eth1DepositIndex,
			Validators: reflect.ValueOf(validators).
				Interface().([]*types.Validator),
			Balances:                     balances,
			RandaoMixes:                  randaoMixes,
			NextWithdrawalIndex:          nextWithdrawalIndex,
			NextWithdrawalValidatorIndex: nextWithdrawalValidatorIndex,
			Slashings:                    slashings,
			TotalSlashing:                totalSlashings,
		}).HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown fork version")
	}
}
