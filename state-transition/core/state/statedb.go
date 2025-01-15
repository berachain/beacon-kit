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
	"context"

	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/storage/beacondb"
)

// StateDB is the underlying struct behind the BeaconState interface.
//
//nolint:revive // todo fix somehow
type StateDB struct {
	beacondb.KVStore

	cs chain.ChainSpec
}

// NewBeaconStateFromDB creates a new beacon state from an underlying state db.
func (s *StateDB) NewFromDB(
	bdb *beacondb.KVStore,
	cs chain.ChainSpec,
) *StateDB {
	return &StateDB{
		KVStore: *bdb,
		cs:      cs,
	}
}

// Copy returns a copy of the beacon state.
func (s *StateDB) Copy(ctx context.Context) *StateDB {
	return s.NewFromDB(s.KVStore.Copy(ctx), s.cs)
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB) IncreaseBalance(
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
func (s *StateDB) DecreaseBalance(
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
func (s *StateDB) UpdateSlashingAtIndex(
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
// NOTE: This function is modified from the spec to allow a fixed withdrawal
// (as the first withdrawal) used for EVM inflation.
//
//nolint:funlen,gocognit // TODO: Simplify when dropping special cases.
func (s *StateDB) ExpectedWithdrawals() (engineprimitives.Withdrawals, error) {
	var (
		validator         *ctypes.Validator
		balance           math.Gwei
		withdrawalAddress common.ExecutionAddress
		withdrawals       = make([]*engineprimitives.Withdrawal, 0)
		withdrawal        *engineprimitives.Withdrawal
	)

	slot, err := s.GetSlot()
	if err != nil {
		return nil, err
	}

	// Handle special cases wherever it's necessary
	switch {
	case s.cs.DepositEth1ChainID() == spec.BartioChainID:
		// nothing special to do

	case s.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot == math.U64(spec.BoonetFork1Height):
		// Slot used to emergency mint EVM tokens on Boonet.
		withdrawals = append(withdrawals, withdrawal.New(
			0, // NOT USED
			0, // NOT USED
			common.NewExecutionAddressFromHex(EVMMintingAddress),
			math.Gwei(EVMMintingAmount),
		))
		return withdrawals, nil

	case s.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height):
		// Boonet inherited the Bartio behaviour pre BoonetFork2Height
		// nothing specific to do

	default:
		// The first withdrawal is fixed to be the EVM inflation withdrawal.
		withdrawals = append(withdrawals, s.EVMInflationWithdrawal())
	}

	epoch := math.Epoch(slot.Unwrap() / s.cs.SlotsPerEpoch())

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

	bound := min(
		totalValidators, s.cs.MaxValidatorsPerWithdrawalsSweep(
			IsPostFork2(s.cs.DepositEth1ChainID(), slot),
		),
	)

	// Iterate through indices to find the next validators to withdraw.
	for range bound {
		validator, err = s.ValidatorByIndex(validatorIndex)
		if err != nil {
			return nil, err
		}

		balance, err = s.GetBalance(validatorIndex)
		if err != nil {
			return nil, err
		}

		// Set the amount of the withdrawal depending on the balance of the
		// validator.
		//nolint:gocritic,nestif // ok.
		if validator.IsFullyWithdrawable(balance, epoch) {
			withdrawalAddress, err = validator.
				GetWithdrawalCredentials().ToExecutionAddress()
			if err != nil {
				return nil, err
			}

			withdrawals = append(withdrawals, withdrawal.New(
				math.U64(withdrawalIndex),
				validatorIndex,
				withdrawalAddress,
				balance,
			))

			// Increment the withdrawal index to process the next withdrawal.
			withdrawalIndex++
		} else if validator.IsPartiallyWithdrawable(
			balance, math.Gwei(s.cs.MaxEffectiveBalance(
				IsPostFork3(s.cs.DepositEth1ChainID(), slot),
			)),
		) {
			withdrawalAddress, err = validator.
				GetWithdrawalCredentials().ToExecutionAddress()
			if err != nil {
				return nil, err
			}

			withdrawals = append(withdrawals, withdrawal.New(
				math.U64(withdrawalIndex),
				validatorIndex,
				withdrawalAddress,
				balance-math.Gwei(s.cs.MaxEffectiveBalance(
					IsPostFork3(s.cs.DepositEth1ChainID(), slot),
				)),
			))

			// Increment the withdrawal index to process the next withdrawal.
			withdrawalIndex++
		} else if s.cs.DepositEth1ChainID() == spec.BartioChainID {
			// Backward compatibility with Bartio
			// TODO: Drop this when we drop other Bartio special cases.

			withdrawalAddress, err = validator.
				GetWithdrawalCredentials().ToExecutionAddress()
			if err != nil {
				return nil, err
			}

			withdrawal = withdrawal.New(
				math.U64(withdrawalIndex),
				validatorIndex,
				withdrawalAddress,
				0,
			)

			withdrawals = append(withdrawals, withdrawal)
			withdrawalIndex++
		}

		// Cap the number of withdrawals to the maximum allowed per payload.
		if uint64(len(withdrawals)) == s.cs.MaxWithdrawalsPerPayload() {
			break
		}

		// Increment the validator index to process the next validator.
		validatorIndex = (validatorIndex + 1) % math.ValidatorIndex(
			totalValidators,
		)
	}

	return withdrawals, nil
}

// EVMInflationWithdrawal returns the withdrawal used for EVM balance inflation.
//
// NOTE: The withdrawal index and validator index are both set to 0 as they are
// not used during processing.
func (s *StateDB) EVMInflationWithdrawal() *engineprimitives.Withdrawal {
	var withdrawal *engineprimitives.Withdrawal
	return withdrawal.New(
		EVMInflationWithdrawalIndex,
		EVMInflationWithdrawalValidatorIndex,
		s.cs.EVMInflationAddress(),
		math.Gwei(s.cs.EVMInflationPerBlock()),
	)
}

// GetMarshallable is the interface for the beacon store.
//
//nolint:funlen,gocognit // todo fix somehow
func (s *StateDB) GetMarshallable() (*ctypes.BeaconState, error) {
	var empty *ctypes.BeaconState

	slot, err := s.GetSlot()
	if err != nil {
		return empty, err
	}

	fork, err := s.GetFork()
	if err != nil {
		return empty, err
	}

	genesisValidatorsRoot, err := s.GetGenesisValidatorsRoot()
	if err != nil {
		return empty, err
	}

	latestBlockHeader, err := s.GetLatestBlockHeader()
	if err != nil {
		return empty, err
	}

	blockRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		blockRoots[i], err = s.GetBlockRootAtIndex(i)
		if err != nil {
			return empty, err
		}
	}

	stateRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		stateRoots[i], err = s.StateRootAtIndex(i)
		if err != nil {
			return empty, err
		}
	}

	latestExecutionPayloadHeader, err := s.GetLatestExecutionPayloadHeader()
	if err != nil {
		return empty, err
	}

	eth1Data, err := s.GetEth1Data()
	if err != nil {
		return empty, err
	}

	eth1DepositIndex, err := s.GetEth1DepositIndex()
	if err != nil {
		return empty, err
	}

	validators, err := s.GetValidators()
	if err != nil {
		return empty, err
	}

	balances, err := s.GetBalances()
	if err != nil {
		return empty, err
	}

	randaoMixes := make([]common.Bytes32, s.cs.EpochsPerHistoricalVector())
	for i := range s.cs.EpochsPerHistoricalVector() {
		randaoMixes[i], err = s.GetRandaoMixAtIndex(i)
		if err != nil {
			return empty, err
		}
	}

	nextWithdrawalIndex, err := s.GetNextWithdrawalIndex()
	if err != nil {
		return empty, err
	}

	nextWithdrawalValidatorIndex, err := s.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return empty, err
	}

	slashings, err := s.GetSlashings()
	if err != nil {
		return empty, err
	}

	totalSlashings, err := s.GetTotalSlashing()
	if err != nil {
		return empty, err
	}

	var bs *ctypes.BeaconState
	return bs.New(
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
}

// HashTreeRoot is the interface for the beacon store.
func (s *StateDB) HashTreeRoot() common.Root {
	st, err := s.GetMarshallable()
	if err != nil {
		panic(err)
	}
	return st.HashTreeRoot()
}
