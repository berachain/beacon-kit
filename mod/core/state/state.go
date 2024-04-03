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

	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/storage/statedb"
)

type beaconState struct {
	*statedb.StateDB
	cfg params.BeaconChainConfig
}

func NewBeaconStateFromStore(sdb *statedb.StateDB, cfg params.BeaconChainConfig) BeaconState {
	return &beaconState{
		StateDB: sdb,
		cfg:     cfg,
	}
}

func (s *beaconState) Copy() BeaconState {
	return NewBeaconStateFromStore(s.StateDB.Copy(), s.cfg)
}

// Store is the interface for the beacon store.
func (s *beaconState) HashTreeRoot() ([32]byte, error) {
	slot, err := s.StateDB.GetSlot()
	if err != nil {
		return [32]byte{}, err
	}

	genesisValidatorsRoot, err := s.StateDB.GetGenesisValidatorsRoot()
	if err != nil {
		return [32]byte{}, err
	}

	latestBlockHeader, err := s.StateDB.GetLatestBlockHeader()
	if err != nil {
		return [32]byte{}, err
	}

	var blockRoot [32]byte
	blockRoots := make([][32]byte, s.cfg.SlotsPerHistoricalRoot)
	for i := uint64(0); i < s.cfg.SlotsPerHistoricalRoot; i++ {
		blockRoot, err = s.StateDB.GetBlockRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		blockRoots[i] = blockRoot
	}

	var stateRoot [32]byte
	stateRoots := make([][32]byte, s.cfg.SlotsPerHistoricalRoot)
	for i := uint64(0); i < s.cfg.SlotsPerHistoricalRoot; i++ {
		stateRoot, err = s.StateDB.StateRootAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		stateRoots[i] = stateRoot
	}

	eth1BlockHash, err := s.StateDB.GetEth1BlockHash()
	if err != nil {
		return [32]byte{}, err
	}

	eth1DepositIndex, err := s.StateDB.GetEth1DepositIndex()
	if err != nil {
		return [32]byte{}, err
	}

	validators, err := s.StateDB.GetValidators()
	if err != nil {
		return [32]byte{}, err
	}

	balances, err := s.StateDB.GetBalances()
	if err != nil {
		return [32]byte{}, err
	}

	var randaoMix [32]byte
	randaoMixes := make([][32]byte, s.cfg.EpochsPerHistoricalVector)
	for i := uint64(0); i < s.cfg.EpochsPerHistoricalVector; i++ {
		randaoMix, err = s.StateDB.GetRandaoMixAtIndex(i)
		if err != nil {
			return [32]byte{}, err
		}
		randaoMixes[i] = randaoMix
	}

	slashings, err := s.StateDB.GetSlashings()
	if err != nil {
		return [32]byte{}, err
	}

	totalSlashings, err := s.StateDB.GetTotalSlashing()
	if err != nil {
		return [32]byte{}, err
	}

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
		NextWithdrawalIndex:          0, // TODO
		NextWithdrawalValidatorIndex: 0, // TODO
		Slashings:                    slashings,
		TotalSlashing:                totalSlashings,
	}).HashTreeRoot()
}

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state consensus.
type BeaconState interface {
	Copy() BeaconState
	Save()
	Context() context.Context
	HashTreeRoot() ([32]byte, error)
	ReadOnlyBeaconState
	WriteOnlyBeaconState
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	ReadOnlyDeposits
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators
	ReadOnlyWithdrawals

	GetSlot() (primitives.Slot, error)
	GetCurrentEpoch() (primitives.Epoch, error)
	GetGenesisValidatorsRoot() (primitives.Root, error)
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	GetLatestBlockHeader() (*types.BeaconBlockHeader, error)
	GetTotalActiveBalances() (primitives.Gwei, error)
	GetValidators() ([]*types.Validator, error)
	GetEth1BlockHash() (primitives.ExecutionHash, error)
	GetTotalSlashing() (primitives.Gwei, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyDeposits
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators
	WriteOnlyWithdrawals
	SetSlot(primitives.Slot) error
	UpdateBlockRootAtIndex(uint64, primitives.Root) error
	SetLatestBlockHeader(*types.BeaconBlockHeader) error
	DecreaseBalance(primitives.ValidatorIndex, primitives.Gwei) error
	UpdateEth1BlockHash(primitives.ExecutionHash) error
	UpdateSlashingAtIndex(uint64, primitives.Gwei) error
}

// WriteOnlyStateRoots defines a struct which only has write access to state
// roots methods.
type WriteOnlyStateRoots interface {
	UpdateStateRootAtIndex(uint64, primitives.Root) error
}

// ReadOnlyStateRoots defines a struct which only has read access to state roots
// methods.
type ReadOnlyStateRoots interface {
	StateRootAtIndex(uint64) (primitives.Root, error)
}

// WriteOnlyRandaoMixes defines a struct which only has write access to randao
// mixes methods.
type WriteOnlyRandaoMixes interface {
	UpdateRandaoMixAtIndex(uint64, primitives.Bytes32) error
}

// ReadOnlyRandaoMixes defines a struct which only has read access to randao
// mixes methods.
type ReadOnlyRandaoMixes interface {
	GetRandaoMixAtIndex(uint64) (primitives.Bytes32, error)
}

// WriteOnlyValidators has write access to validator methods.
type WriteOnlyValidators interface {
	// Add methods here
	UpdateValidatorAtIndex(
		primitives.ValidatorIndex,
		*types.Validator,
	) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators interface {
	ValidatorIndexByPubkey(
		[]byte,
	) (primitives.ValidatorIndex, error)

	ValidatorByIndex(
		primitives.ValidatorIndex,
	) (*types.Validator, error)
}

// ReadWriteValidators has read and write access to validator methods.
type ReadWriteDeposits interface {
	ReadOnlyDeposits
	WriteOnlyDeposits
}

// ReadWriteDepositQueue has read and write access to deposit queue.
type WriteOnlyDeposits interface {
	SetEth1DepositIndex(uint64) error
	EnqueueDeposits([]*types.Deposit) error
	DequeueDeposits(uint64) ([]*types.Deposit, error)
}

// ReadOnlyDeposits has read access to deposit queue.
type ReadOnlyDeposits interface {
	GetEth1DepositIndex() (uint64, error)
	ExpectedDeposits(uint64) ([]*types.Deposit, error)
}

// ReadWriteWithdrawals has read and write access to withdrawal methods.
type ReadWriteWithdrawals interface {
	ReadOnlyWithdrawals
	WriteOnlyWithdrawals
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	ExpectedWithdrawals(uint64) ([]*primitives.Withdrawal, error)
}

// WriteOnlyWithdrawals only has write access to withdrawal methods.
type WriteOnlyWithdrawals interface {
	EnqueueWithdrawals([]*primitives.Withdrawal) error
	DequeueWithdrawals(uint64) ([]*primitives.Withdrawal, error)
}
