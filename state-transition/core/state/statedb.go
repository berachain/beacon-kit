// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/storage/beacondb"
)

// StateDB is the underlying struct behind the BeaconState interface.
//
//nolint:revive // todo fix somehow
type StateDB struct {
	beacondb.KVStore

	cs ChainSpec
}

// NewBeaconStateFromDB creates a new beacon state from an underlying state db.
func NewBeaconStateFromDB(bdb *beacondb.KVStore, cs ChainSpec) *StateDB {
	return &StateDB{
		KVStore: *bdb,
		cs:      cs,
	}
}

// Copy returns a copy of the beacon state.
func (s *StateDB) Copy(ctx context.Context) *StateDB {
	return NewBeaconStateFromDB(s.KVStore.Copy(ctx), s.cs)
}

// IncreaseBalance increases the balance of a validator.
func (s *StateDB) IncreaseBalance(idx math.ValidatorIndex, delta math.Gwei) error {
	balance, err := s.GetBalance(idx)
	if err != nil {
		return err
	}
	return s.SetBalance(idx, balance+delta)
}

// DecreaseBalance decreases the balance of a validator.
func (s *StateDB) DecreaseBalance(idx math.ValidatorIndex, delta math.Gwei) error {
	balance, err := s.GetBalance(idx)
	if err != nil {
		return err
	}
	return s.SetBalance(idx, balance-min(balance, delta))
}

// ExpectedWithdrawals as defined in the Ethereum 2.0 Specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/electra/beacon-chain.md#new-get_expected_withdrawals
//
// NOTE: This function is modified from the spec to allow a fixed withdrawal
// (as the first withdrawal) used for EVM inflation.
//
//nolint:gocognit,funlen // TODO(pectra).
func (s *StateDB) ExpectedWithdrawals(timestamp math.U64) (engineprimitives.Withdrawals, uint64, error) {
	var (
		validator         *ctypes.Validator
		balance           math.Gwei
		withdrawalAddress common.ExecutionAddress
	)

	processedPartialWithdrawals := uint64(0)

	slot, err := s.GetSlot()
	if err != nil {
		return nil, 0, err
	}
	epoch := s.cs.SlotToEpoch(slot)
	maxWithdrawals := s.cs.MaxWithdrawalsPerPayload()
	withdrawals := make([]*engineprimitives.Withdrawal, 0, maxWithdrawals)

	// The first withdrawal is fixed to be the EVM inflation withdrawal.
	withdrawals = append(withdrawals, s.EVMInflationWithdrawal(timestamp))

	withdrawalIndex, err := s.GetNextWithdrawalIndex()
	if err != nil {
		return nil, 0, err
	}

	validatorIndex, err := s.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return nil, 0, err
	}

	totalValidators, err := s.GetTotalValidators()
	if err != nil {
		return nil, 0, err
	}

	// [New in Electra:EIP7251] Consume pending partial withdrawals
	timestampFork := s.cs.ActiveForkVersionForTimestamp(timestamp)
	if version.EqualsOrIsAfter(timestampFork, version.Electra()) {
		withdrawals, withdrawalIndex, processedPartialWithdrawals, err =
			s.consumePendingPartialWithdrawals(timestampFork, epoch, withdrawals, withdrawalIndex)
		if err != nil {
			return nil, 0, err
		}
	}

	bound := min(totalValidators, s.cs.MaxValidatorsPerWithdrawalsSweep())

	// Iterate through indices to find the next validators to withdraw.
	for range bound {
		validator, err = s.ValidatorByIndex(validatorIndex)
		if err != nil {
			return nil, 0, err
		}

		balance, err = s.GetBalance(validatorIndex)
		if err != nil {
			return nil, 0, err
		}

		// Set the amount of the withdrawal depending on the balance of the validator.
		if validator.IsFullyWithdrawable(balance, epoch) {
			withdrawalAddress, err = validator.GetWithdrawalCredentials().ToExecutionAddress()
			if err != nil {
				return nil, 0, err
			}

			withdrawals = append(withdrawals, engineprimitives.NewWithdrawal(
				math.U64(withdrawalIndex),
				validatorIndex,
				withdrawalAddress,
				balance,
			))

			// Increment the withdrawal index to process the next withdrawal.
			withdrawalIndex++
		} else if validator.IsPartiallyWithdrawable(
			balance, math.Gwei(s.cs.MaxEffectiveBalance()),
		) {
			withdrawalAddress, err = validator.GetWithdrawalCredentials().ToExecutionAddress()
			if err != nil {
				return nil, 0, err
			}

			withdrawals = append(withdrawals, engineprimitives.NewWithdrawal(
				math.U64(withdrawalIndex),
				validatorIndex,
				withdrawalAddress,
				balance-math.Gwei(s.cs.MaxEffectiveBalance()),
			))

			// Increment the withdrawal index to process the next withdrawal.
			withdrawalIndex++
		}

		// Cap the number of withdrawals to the maximum allowed per payload.
		if uint64(len(withdrawals)) == maxWithdrawals {
			break
		}

		// Increment the validator index to process the next validator.
		validatorIndex = (validatorIndex + 1) % math.ValidatorIndex(totalValidators)
	}

	return withdrawals, processedPartialWithdrawals, nil
}

func (s *StateDB) consumePendingPartialWithdrawals(
	timestampFork common.Version,
	epoch math.Epoch,
	withdrawals engineprimitives.Withdrawals,
	withdrawalIndex uint64,
) (
	[]*engineprimitives.Withdrawal,
	uint64, // withdrawalIndex
	uint64, // processedPartialWithdrawals
	error,
) {
	// By this point, if we're post-Electra, the fork version on the BeaconState will have been set as part of `PrepareStateForFork`.
	// We sanity check this by enforcing the state fork and timestamp version are equal.
	// We cannot do this check before Electra as Deneb1 did not have a beacon state fork version set.
	stateFork, forkErr := s.GetFork()
	if forkErr != nil {
		return nil, 0, 0, forkErr
	}
	if !version.Equals(stateFork.CurrentVersion, timestampFork) {
		return nil, 0, 0, errors.New("state fork version does not match timestamp fork in expected withdrawals")
	}
	ppWithdrawals, getErr := s.GetPendingPartialWithdrawals()
	if getErr != nil {
		return nil, 0, 0, getErr
	}
	processedPartialWithdrawals := uint64(0)

	for _, withdrawal := range ppWithdrawals {
		if withdrawal.WithdrawableEpoch > epoch || len(withdrawals) == int(constants.MaxPendingPartialsPerWithdrawalsSweep) {
			// If the first withdrawal in the queue is not withdrawable, then all subsequent withdrawals will also be in later
			// epochs and hence are not withdrawable, so we can break early.
			break
		}
		validator, err := s.ValidatorByIndex(withdrawal.ValidatorIndex)
		if err != nil {
			return nil, 0, 0, err
		}
		minActivationBalance := math.Gwei(s.cs.MinActivationBalance())
		hasSufficientEffectiveBalance := validator.GetEffectiveBalance() >= minActivationBalance
		balance, err := s.GetBalance(withdrawal.ValidatorIndex)
		if err != nil {
			return nil, 0, 0, err
		}
		hasExcessBalance := balance > minActivationBalance
		if validator.ExitEpoch == math.Epoch(constants.FarFutureEpoch) && hasSufficientEffectiveBalance && hasExcessBalance {
			// A validator can only partial withdraw an amount such that:
			// 1. never withdraw more than what the validator asked for.
			// 2. never withdraw so much that the validator’s remaining balance would drop below MIN_ACTIVATION_BALANCE
			withdrawableBalance := min(balance-minActivationBalance, withdrawal.Amount)

			withdrawalAddress, addrErr := validator.WithdrawalCredentials.ToExecutionAddress()
			if addrErr != nil {
				return nil, 0, 0, addrErr
			}
			withdrawals = append(
				withdrawals,
				engineprimitives.NewWithdrawal(math.U64(withdrawalIndex), withdrawal.ValidatorIndex, withdrawalAddress, withdrawableBalance),
			)
			// Increment the withdrawal index to process the next withdrawal.
			withdrawalIndex++
		}
		// Even if a withdrawal was not created, e.g. the validator did not have sufficient balance, we will consider
		// this withdrawal processed (spec defined) and hence increment the processedPartialWithdrawals count.
		processedPartialWithdrawals++
	}
	return withdrawals, withdrawalIndex, processedPartialWithdrawals, nil
}

// EVMInflationWithdrawal returns the withdrawal used for EVM balance inflation.
//
// NOTE: The withdrawal index and validator index are both set to max(uint64) as
// they are not used during processing.
func (s *StateDB) EVMInflationWithdrawal(timestamp math.U64) *engineprimitives.Withdrawal {
	return engineprimitives.NewWithdrawal(
		EVMInflationWithdrawalIndex,
		EVMInflationWithdrawalValidatorIndex,
		s.cs.EVMInflationAddress(timestamp),
		math.Gwei(s.cs.EVMInflationPerBlock(timestamp)),
	)
}

// GetMarshallable is the interface for the beacon store.
//
//nolint:funlen,gocognit // todo fix somehow
func (s *StateDB) GetMarshallable() (*ctypes.BeaconState, error) {
	slot, err := s.GetSlot()
	if err != nil {
		return nil, err
	}

	fork, err := s.GetFork()
	if err != nil {
		return nil, err
	}
	genesisValidatorsRoot, err := s.GetGenesisValidatorsRoot()
	if err != nil {
		return nil, err
	}

	latestBlockHeader, err := s.GetLatestBlockHeader()
	if err != nil {
		return nil, err
	}

	blockRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		blockRoots[i], err = s.GetBlockRootAtIndex(i)
		if err != nil {
			return nil, err
		}
	}

	stateRoots := make([]common.Root, s.cs.SlotsPerHistoricalRoot())
	for i := range s.cs.SlotsPerHistoricalRoot() {
		stateRoots[i], err = s.StateRootAtIndex(i)
		if err != nil {
			return nil, err
		}
	}

	latestExecutionPayloadHeader, err := s.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}

	eth1Data, err := s.GetEth1Data()
	if err != nil {
		return nil, err
	}

	eth1DepositIndex, err := s.GetEth1DepositIndex()
	if err != nil {
		return nil, err
	}

	validators, err := s.GetValidators()
	if err != nil {
		return nil, err
	}

	balances, err := s.GetBalances()
	if err != nil {
		return nil, err
	}

	randaoMixes := make([]common.Bytes32, s.cs.EpochsPerHistoricalVector())
	for i := range s.cs.EpochsPerHistoricalVector() {
		randaoMixes[i], err = s.GetRandaoMixAtIndex(i)
		if err != nil {
			return nil, err
		}
	}

	nextWithdrawalIndex, err := s.GetNextWithdrawalIndex()
	if err != nil {
		return nil, err
	}

	nextWithdrawalValidatorIndex, err := s.GetNextWithdrawalValidatorIndex()
	if err != nil {
		return nil, err
	}

	slashings, err := s.GetSlashings()
	if err != nil {
		return nil, err
	}

	totalSlashings, err := s.GetTotalSlashing()
	if err != nil {
		return nil, err
	}

	beaconState := ctypes.NewEmptyBeaconStateWithVersion(fork.CurrentVersion)
	beaconState.Slot = slot
	beaconState.GenesisValidatorsRoot = genesisValidatorsRoot
	beaconState.Fork = fork
	beaconState.LatestBlockHeader = latestBlockHeader
	beaconState.BlockRoots = blockRoots
	beaconState.StateRoots = stateRoots
	beaconState.LatestExecutionPayloadHeader = latestExecutionPayloadHeader
	beaconState.Eth1Data = eth1Data
	beaconState.Eth1DepositIndex = eth1DepositIndex
	beaconState.Validators = validators
	beaconState.Balances = balances
	beaconState.RandaoMixes = randaoMixes
	beaconState.NextWithdrawalIndex = nextWithdrawalIndex
	beaconState.NextWithdrawalValidatorIndex = nextWithdrawalValidatorIndex
	beaconState.Slashings = slashings
	beaconState.TotalSlashing = totalSlashings

	if version.EqualsOrIsAfter(beaconState.GetForkVersion(), version.Electra()) {
		pendingPartialWithdrawals, getErr := s.GetPendingPartialWithdrawals()
		if getErr != nil {
			return nil, getErr
		}
		beaconState.PendingPartialWithdrawals = pendingPartialWithdrawals
	}

	return beaconState, nil
}

// HashTreeRoot is the interface for the beacon store.
func (s *StateDB) HashTreeRoot() common.Root {
	st, err := s.GetMarshallable()
	if err != nil {
		panic(err)
	}
	return st.HashTreeRoot()
}
