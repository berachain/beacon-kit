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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody,
	BeaconStateT state.BeaconState,
	BlobSidecarsT interface{ Len() int },
	ContextT Context,
] struct {
	cs              primitives.ChainSpec
	rp              RandaoProcessor[BeaconBlockT, BeaconStateT]
	signer          crypto.BLSSigner
	logger          log.Logger[any]
	executionEngine ExecutionEngine
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT],
	BeaconBlockBodyT BeaconBlockBody,
	BeaconStateT state.BeaconState,
	BlobSidecarsT interface{ Len() int },
	ContextT Context,
](
	cs primitives.ChainSpec,
	rp RandaoProcessor[
		BeaconBlockT, BeaconStateT,
	],
	executionEngine ExecutionEngine,
	signer crypto.BLSSigner,
	logger log.Logger[any],
) *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT, ContextT] {
	return &StateProcessor[
		BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, ContextT,
	]{
		cs:              cs,
		rp:              rp,
		executionEngine: executionEngine,
		signer:          signer,
		logger:          logger,
	}
}

// Transition is the main function for processing a state transition.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) Transition(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	blkSlot := blk.GetSlot()
	stateSlot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// We perform some initial logic to ensure the BeaconState is in the correct
	// state before we process the block.
	//
	//            +-------------------------------+
	//            |  Is state slot equal to the   |
	//            |  block slot minus one?        |
	//            +-------------------------------+
	//                           |
	//                           |
	//              +------------+------------+
	//              |                         |
	//           Yes, it is               No, it isn't
	//              |                         |
	//              |                         |
	//       Process the slot                 |
	//                                        |
	//                           +------------+
	//                           |
	//          Is state slot equal to the block slot?
	//                           |
	//              +------------+------------+
	//              |                         |
	//           Yes, it is               No, it isn't
	//              |                         |
	//     Skip slot processing          Return error:
	//                                   "out of sync"
	//
	// Unlike Ethereum, we error if the on disk state is greater than 1 slot
	// behind.
	// Due to CometBFT SSF nature, this SHOULD NEVER occur.
	//
	// TODO: We should probably not assume this to make our Transition
	// function more generalizable, since right now it makes an
	// assumption about the finalization properties of the cosnensus
	// engine.
	switch stateSlot {
	case blkSlot - 1:
		if err = sp.ProcessSlot(st); err != nil {
			return err
		}
	case blkSlot:
		// skip slot processing.
	default:
		return errors.Wrapf(
			ErrBeaconStateOutOfSync, "expected: %d, got: %d",
			stateSlot, blkSlot,
		)
	}

	// Process the block.
	if err = sp.ProcessBlock(ctx, st, blk); err != nil {
		sp.logger.Error(
			"failed to process block",
			"slot", blkSlot,
			"error", err,
		)
		return err
	}

	// We only want to persist state changes if we successfully
	// processed the block.
	st.Save()
	return nil
}

// ProcessSlot is run when a slot is missed.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) ProcessSlot(
	st BeaconStateT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Before we make any changes, we calculate the previous state root.
	prevStateRoot, err := st.HashTreeRoot()
	if err != nil {
		return err
	}

	// We update our state roots and block roots.
	if err = st.UpdateStateRootAtIndex(
		uint64(slot)%sp.cs.SlotsPerHistoricalRoot(),
		prevStateRoot,
	); err != nil {
		return err
	}

	// We get the latest block header, this will not have
	// a state root on it.
	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}

	// We set the "rawHeader" in the StateProcessor, but cannot fill in
	// the StateRoot until the following block.
	if (latestHeader.StateRoot == primitives.Root{}) {
		latestHeader.StateRoot = prevStateRoot
		if err = st.SetLatestBlockHeader(latestHeader); err != nil {
			return err
		}
	}

	// We update the block root.
	var prevBlockRoot primitives.Root
	prevBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	if err = st.UpdateBlockRootAtIndex(
		uint64(slot)%sp.cs.SlotsPerHistoricalRoot(), prevBlockRoot,
	); err != nil {
		return err
	}

	// Process the Epoch Boundary.
	if uint64(slot+1)%sp.cs.SlotsPerEpoch() == 0 {
		if err = sp.processEpoch(st); err != nil {
			return err
		}
		sp.logger.Info(
			"processed epoch transition 🔃",
			"old", uint64(slot)/sp.cs.SlotsPerEpoch(),
			"new", uint64(slot+1)/sp.cs.SlotsPerEpoch(),
		)
	}

	return st.SetSlot(slot + 1)
}

// ProcessBlock processes the block and ensures it matches the local state.
//
//nolint:funlen // todo fix.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) ProcessBlock(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	// process the freshly created header.
	if err := sp.processBlockHeader(st, blk); err != nil {
		return err
	}

	// process the execution payload.
	if err := sp.processExecutionPayload(
		ctx, st, blk,
	); err != nil {
		return err
	}

	// process the withdrawals.
	if err := sp.processWithdrawals(
		st, blk.GetBody(),
	); err != nil {
		return err
	}

	// phase0.ProcessProposerSlashings
	// phase0.ProcessAttesterSlashings

	// process the randao reveal.
	if err := sp.processRandaoReveal(st, blk); err != nil {
		return err
	}

	// phase0.ProcessEth1Vote ? forkchoice?

	// TODO: LOOK HERE
	//
	// process the deposits and ensure they match the local state.
	if err := sp.processOperations(st, blk); err != nil {
		return err
	}

	if ctx.GetValidateResult() {
		// Ensure the state root matches the block.
		//
		// TODO: We need to validate this in ProcessProposal as well.
		if stateRoot, err := st.HashTreeRoot(); err != nil {
			return err
		} else if blk.GetStateRoot() != stateRoot {
			return errors.Wrapf(
				ErrStateRootMismatch, "expected %s, got %s",
				primitives.Root(stateRoot), blk.GetStateRoot(),
			)
		}
	}
	return nil
}

// processEpoch processes the epoch and ensures it matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) processEpoch(
	st BeaconStateT,
) error {
	if err := sp.processRewardsAndPenalties(st); err != nil {
		return err
	} else if err = sp.processSlashingsReset(st); err != nil {
		return err
	}
	return sp.processRandaoMixesReset(st)
}

// processBlockHeader processes the header and ensures it matches the local
// state.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) processBlockHeader(
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	var (
		slot              math.Slot
		err               error
		latestBlockHeader *types.BeaconBlockHeader
		parentBlockRoot   primitives.Root
		bodyRoot          primitives.Root
		proposer          *types.Validator
	)

	// Ensure the block slot matches the state slot.
	if slot, err = st.GetSlot(); err != nil {
		return err
	} else if blk.GetSlot() != slot {
		return errors.Wrapf(ErrSlotMismatch,
			"expected: %d, got: %d",
			slot, blk.GetSlot(),
		)
	}

	// Verify the parent block root is correct.
	if latestBlockHeader, err = st.GetLatestBlockHeader(); err != nil {
		return err
	} else if blk.GetSlot() <= latestBlockHeader.GetSlot() {
		return errors.Wrapf(
			ErrBlockSlotTooLow, "expected: > %d, got: %d",
			latestBlockHeader.GetSlot(), blk.GetSlot(),
		)
	} else if parentBlockRoot, err = latestBlockHeader.HashTreeRoot(); err != nil {
		return err
	} else if parentBlockRoot != blk.GetParentBlockRoot() {
		return errors.Wrapf(ErrParentRootMismatch,
			"expected: %x, got: %x",
			parentBlockRoot, blk.GetParentBlockRoot(),
		)
	}

	// Ensure the block is within the acceptable range.
	// TODO: move this is in the wrong spot.
	deposits := blk.GetBody().GetDeposits()
	if uint64(len(deposits)) > sp.cs.MaxDepositsPerBlock() {
		return errors.Wrapf(ErrExceedsBlockDepositLimit,
			"expected: %d, got: %d",
			sp.cs.MaxDepositsPerBlock(), len(deposits),
		)
	}

	// Calculate the body root to place on the header.
	if bodyRoot, err = blk.GetBody().HashTreeRoot(); err != nil {
		return err
	} else if err = st.SetLatestBlockHeader(
		types.NewBeaconBlockHeader(
			blk.GetSlot(),
			blk.GetProposerIndex(),
			blk.GetParentBlockRoot(),
			// state_root is zeroed and overwritten
			// in the next `process_slot` call.
			[32]byte{},
			bodyRoot,
		),
	); err != nil {
		return err
	}

	// Check to make sure the proposer isn't slashed.
	if proposer, err = st.ValidatorByIndex(blk.GetProposerIndex()); err != nil {
		return err
	} else if proposer.Slashed {
		return errors.Wrapf(
			ErrSlashedProposer, "index: %d", blk.GetProposerIndex(),
		)
	}
	return nil
}

// getAttestationDeltas as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#get_attestation_deltas
//
//nolint:lll
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) getAttestationDeltas(
	st BeaconStateT,
) ([]math.Gwei, []math.Gwei, error) {
	// TODO: implement this function forreal
	validators, err := st.GetValidators()
	if err != nil {
		return nil, nil, err
	}
	placeholder := make([]math.Gwei, len(validators))
	return placeholder, placeholder, nil
}

// processRewardsAndPenalties as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#process_rewards_and_penalties
//
//nolint:lll
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, ContextT,
]) processRewardsAndPenalties(
	st BeaconStateT,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	if sp.cs.SlotToEpoch(slot) == math.U64(constants.GenesisEpoch) {
		return nil
	}

	rewards, penalties, err := sp.getAttestationDeltas(st)
	if err != nil {
		return err
	}

	validators, err := st.GetValidators()
	if err != nil {
		return err
	}

	if len(validators) != len(rewards) {
		return errors.Wrapf(
			ErrRewardsLengthMismatch, "expected: %d, got: %d",
			len(validators), len(rewards),
		)
	} else if len(validators) != len(penalties) {
		return errors.Wrapf(
			ErrPenaltiesLengthMismatch, "expected: %d, got: %d",
			len(validators), len(penalties),
		)
	}

	for i := range validators {
		// Increase the balance of the validator.
		if err = st.IncreaseBalance(
			math.ValidatorIndex(i),
			rewards[i],
		); err != nil {
			return err
		}

		// Decrease the balance of the validator.
		if err = st.DecreaseBalance(
			math.ValidatorIndex(i),
			penalties[i],
		); err != nil {
			return err
		}
	}

	return nil
}
