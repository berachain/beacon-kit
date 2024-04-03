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
		NextWithdrawalIndex:          0, // TODO
		NextWithdrawalValidatorIndex: 0, // TODO
		Slashings:                    slashings,
		TotalSlashing:                totalSlashings,
	}).HashTreeRoot()
}
