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

package beacon

import (
	"github.com/berachain/beacon-kit/mod/core/state"
)

// Store is the interface for the beacon store.
func (s *Store) HashTreeRoot() ([32]byte, error) {
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

	slashings, err := s.GetSlashings()
	if err != nil {
		return [32]byte{}, err
	}

	totalSlashings, err := s.GetTotalSlashing()
	if err != nil {
		return [32]byte{}, err
	}

	return (&state.BeaconStateDeneb{
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
