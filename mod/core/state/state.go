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
	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/statedb"
)

// beaconState is a wrapper around the state db that implements the BeaconState
// interface.
type beaconState struct {
	*statedb.StateDB
	cfg *params.BeaconChainConfig
}

// NewBeaconState creates a new beacon state from an underlying state db.
func NewBeaconStateFromDB(
	sdb *statedb.StateDB,
	cfg *params.BeaconChainConfig,
) BeaconState {
	return &beaconState{
		StateDB: sdb,
		cfg:     cfg,
	}
}

// Copy returns a copy of the beacon state.
func (s *beaconState) Copy() BeaconState {
	return NewBeaconStateFromDB(s.StateDB.Copy(), s.cfg)
}

// ExpectedWithdrawals as defined in the Ethereum 2.0 Specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-get_expected_withdrawals
//
//nolint:lll
func (s *beaconState) ExpectedWithdrawals() ([]*primitives.Withdrawal, error) {
	var (
		validator         *types.Validator
		balance           primitives.Gwei
		withdrawalAddress primitives.ExecutionAddress
		withdrawals       = make([]*primitives.Withdrawal, 0)
	)

	slot, err := s.GetSlot()
	if err != nil {
		return nil, err
	}

	epoch := primitives.Epoch(uint64(slot) / s.cfg.SlotsPerEpoch)

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
	for i := uint64(0); i < min(
		s.cfg.MaxValidatorsPerWithdrawalsSweep, totalValidators,
	); i++ {
		validator, err = s.ValidatorByIndex(validatorIndex)
		if err != nil {
			return nil, err
		}

		balance, err = s.GetBalance(validatorIndex)
		if err != nil {
			return nil, err
		}

		withdrawalAddress, err = validator.
			WithdrawalCredentials.ToExecutionAddress()
		if err != nil {
			return nil, err
		}

		// These fields are the same for both partial and full withdrawals.
		withdrawal := &primitives.Withdrawal{
			Index:     withdrawalIndex,
			Validator: validatorIndex,
			Address:   withdrawalAddress,
		}

		// Set the amount of the withdrawal depending on the balance of the
		// validator.
		if validator.IsFullyWithdrawable(balance, epoch) {
			withdrawal.Amount = balance
		} else if validator.IsPartiallyWithdrawable(balance, primitives.Gwei(s.cfg.MaxEffectiveBalance)) {
			withdrawal.Amount = balance - primitives.Gwei(s.cfg.MaxEffectiveBalance)
		}
		withdrawals = append(withdrawals, withdrawal)

		// Increment the withdrawal index to process the next withdrawal.
		withdrawalIndex++

		// Cap the number of withdrawals to the maximum allowed per payload.
		//#nosec:G701 // won't overflow in practice.
		if len(withdrawals) == int(s.cfg.MaxWithdrawalsPerPayload) {
			break
		}

		// Increment the validator index to process the next validator.
		validatorIndex = (validatorIndex + 1) % primitives.ValidatorIndex(
			totalValidators,
		)
	}

	return withdrawals, nil
}

// Store is the interface for the beacon store.
//
//nolint:funlen // todo fix somehow
func (s *beaconState) HashTreeRoot() ([32]byte, error) {
	slot, err := s.GetSlot()
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

	var blockRoot [32]byte
	blockRoots := make([][32]byte, s.cfg.SlotsPerHistoricalRoot)
	for i := uint64(0); i < s.cfg.SlotsPerHistoricalRoot; i++ {
		blockRoot, err = s.GetBlockRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		blockRoots[i] = blockRoot
	}

	var stateRoot [32]byte
	stateRoots := make([][32]byte, s.cfg.SlotsPerHistoricalRoot)
	for i := uint64(0); i < s.cfg.SlotsPerHistoricalRoot; i++ {
		stateRoot, err = s.StateRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		stateRoots[i] = stateRoot
	}

	eth1BlockHash, err := s.GetEth1BlockHash()
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

	var randaoMix [32]byte
	randaoMixes := make([][32]byte, s.cfg.EpochsPerHistoricalVector)
	for i := uint64(0); i < s.cfg.EpochsPerHistoricalVector; i++ {
		randaoMix, err = s.GetRandaoMixAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		randaoMixes[i] = randaoMix
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

	// TODO: handle hardforks.
	return (&BeaconStateDeneb{
		Slot:                         slot,
		GenesisValidatorsRoot:        genesisValidatorsRoot,
		LatestBlockHeader:            latestBlockHeader,
		BlockRoots:                   blockRoots,
		StateRoots:                   stateRoots,
		Eth1BlockHash:                eth1BlockHash,
		Eth1DepositIndex:             eth1DepositIndex,
		Validators:                   validators,
		Balances:                     balances,
		RandaoMixes:                  randaoMixes,
		NextWithdrawalIndex:          nextWithdrawalIndex,
		NextWithdrawalValidatorIndex: nextWithdrawalValidatorIndex,
		Slashings:                    slashings,
		TotalSlashing:                totalSlashings,
	}).HashTreeRoot()
}
