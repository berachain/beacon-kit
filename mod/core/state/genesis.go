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
	"github.com/berachain/beacon-kit/mod/core/state/deneb"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// WriteGenesisStateDeneb writes the genesis state to the state db.
// TODO: Implement the "from eth1 version" and move genesis.json
// to a custom type, not the full state.
//
//nolint:gocognit // splitting into more functions would be confusing.
func (s *StateDB) WriteGenesisStateDeneb(st *deneb.BeaconState) error {
	if err := s.SetGenesisValidatorsRoot(st.GenesisValidatorsRoot); err != nil {
		return err
	}

	if err := s.SetSlot(st.Slot); err != nil {
		return err
	}

	if err := s.SetFork(st.Fork); err != nil {
		return err
	}

	if err := s.SetLatestBlockHeader(st.LatestBlockHeader); err != nil {
		return err
	}

	for i, root := range st.BlockRoots {
		//#nosec:G701 // won't overflow in practice.
		if err := s.UpdateBlockRootAtIndex(uint64(i), root); err != nil {
			return err
		}
	}

	for i, root := range st.StateRoots {
		//#nosec:G701 // won't overflow in practice.
		if err := s.UpdateStateRootAtIndex(uint64(i), root); err != nil {
			return err
		}
	}

	if err := s.UpdateLatestExecutionPayload(
		st.LatestExecutionPayload,
	); err != nil {
		return err
	}

	if err := s.SetEth1Data(st.Eth1Data); err != nil {
		return err
	}

	if err := s.SetEth1DepositIndex(st.Eth1DepositIndex); err != nil {
		return err
	}

	for _, validator := range st.Validators {
		if err := s.AddValidator(validator); err != nil {
			return err
		}
	}

	for i, mix := range st.RandaoMixes {
		//#nosec:G701 // won't overflow in practice.
		if err := s.UpdateRandaoMixAtIndex(uint64(i), mix); err != nil {
			return err
		}
	}

	if err := s.SetNextWithdrawalIndex(st.NextWithdrawalIndex); err != nil {
		return err
	}

	if err := s.SetNextWithdrawalValidatorIndex(
		st.NextWithdrawalValidatorIndex,
	); err != nil {
		return err
	}

	totalSlashing := primitives.Gwei(0)
	for i, slashing := range st.Slashings {
		totalSlashing += primitives.Gwei(slashing)
		//#nosec:G701 // won't overflow in practice.
		if err := s.SetSlashingAtIndex(
			uint64(i), primitives.Gwei(slashing),
		); err != nil {
			return err
		}
	}

	return s.SetTotalSlashing(totalSlashing)
}
