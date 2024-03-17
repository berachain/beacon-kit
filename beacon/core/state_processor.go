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

package core

import (
	"fmt"

	"github.com/berachain/beacon-kit/beacon/core/state"
	"github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/config"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor struct {
	cfg *config.Beacon
	rp  RandaoProcessor
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor(
	cfg *config.Beacon,
	rp RandaoProcessor,
) *StateProcessor {
	return &StateProcessor{
		cfg: cfg,
		rp:  rp,
	}
}

// ProcessSlot processes the slot and ensures it matches the local state.
func (sp *StateProcessor) ProcessSlot(
	_ state.BeaconState,
	_ uint64,
) error {
	return nil
}

// ProcessBlock processes the block and ensures it matches the local state.
func (sp *StateProcessor) ProcessBlock(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	// Ensure Body is non nil.
	body := blk.GetBody()
	if body.IsNil() {
		return types.ErrNilBlkBody
	}

	// process the eth1 vote.
	payload := body.GetExecutionPayload()
	if payload.IsNil() {
		return types.ErrNilPayload
	}

	// common.ProcessHeader

	// process the withdrawals.
	if err := sp.processWithdrawals(st, payload.GetWithdrawals()); err != nil {
		return err
	}

	// phase0.ProcessProposerSlashings
	// phase0.ProcessAttesterSlashings

	// process the randao reveal.
	if err := sp.processRandaoReveal(st, blk); err != nil {
		return err
	}

	// phase0.ProcessEth1Vote ? forkchoice?

	// process the deposits and ensure they match the local state.
	if err := sp.processDeposits(st, body.GetDeposits()); err != nil {
		return err
	}

	// ProcessVoluntaryExits

	return nil
}

// ProcessBlob processes a blob.
func (sp *StateProcessor) ProcessBlob(_ state.BeaconState) error {
	// TODO: 4844.
	return nil
}

// ProcessDeposits processes the deposits and ensures they match the
// local state.
func (sp *StateProcessor) processDeposits(
	st state.BeaconState,
	deposits []*types.Deposit,
) error {
	// Dequeue and verify the logs.
	localDeposits, err := st.ExpectedDeposits(uint64(len(deposits)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, dep := range deposits {
		if dep == nil {
			return types.ErrNilDeposit
		}
		if dep.Index != localDeposits[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localDeposits[i].Index, dep.Index)
		}
	}
	return nil
}

// processWithdrawals processes the withdrawals and ensures they match the
// local state.
func (sp *StateProcessor) processWithdrawals(
	st state.BeaconState,
	withdrawals []*enginetypes.Withdrawal,
) error {
	// Dequeue and verify the withdrawals.
	localWithdrawals, err := st.DequeueWithdrawals(uint64(len(withdrawals)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, dep := range withdrawals {
		if dep == nil {
			return types.ErrNilWithdrawal
		}
		if dep.Index != localWithdrawals[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localWithdrawals[i].Index, dep.Index)
		}
	}
	return nil
}

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor) processRandaoReveal(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	// Ensure the proposer index is valid.
	pubkey, err := st.ValidatorPubKeyByIndex(blk.GetProposerIndex())
	if err != nil {
		return err
	}

	// Verify the RANDAO Reveal.
	reveal := blk.GetBody().GetRandaoReveal()
	if err = sp.rp.VerifyReveal(
		st,
		[bls12381.PubKeyLength]byte(pubkey),
		reveal,
	); err != nil {
		return err
	}

	// Mixin the reveal.
	return sp.rp.MixinNewReveal(st, reveal)
}
