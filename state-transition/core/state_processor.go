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

package core

import (
	"bytes"

	"github.com/berachain/beacon-kit/chain-spec/chain"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core/state"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor[
	ContextT Context,
	KVStoreT any,
] struct {
	// logger is used for logging information and errors.
	logger log.Logger
	// cs is the chain specification for the beacon chain.
	cs chain.ChainSpec
	// signer is the BLS signer used for cryptographic operations.
	signer crypto.BLSSigner
	// fGetAddressFromPubKey verifies that a validator public key
	// matches with the proposer address passed by the consensus
	// Injected via ctor to simplify testing.
	fGetAddressFromPubKey func(crypto.BLSPubkey) ([]byte, error)
	// executionEngine is the engine responsible for executing transactions.
	executionEngine ExecutionEngine
	// ds allows checking payload deposits against the deposit contract
	ds DepositStore
	// metrics is the metrics for the service.
	metrics *stateProcessorMetrics
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor[
	ContextT Context,
	KVStoreT any,
](
	logger log.Logger,
	cs chain.ChainSpec,
	executionEngine ExecutionEngine,
	ds DepositStore,
	signer crypto.BLSSigner,
	fGetAddressFromPubKey func(crypto.BLSPubkey) ([]byte, error),
	telemetrySink TelemetrySink,
) *StateProcessor[
	ContextT,
	KVStoreT,
] {
	return &StateProcessor[
		ContextT,
		KVStoreT,
	]{
		logger:                logger,
		cs:                    cs,
		executionEngine:       executionEngine,
		signer:                signer,
		fGetAddressFromPubKey: fGetAddressFromPubKey,
		ds:                    ds,
		metrics:               newStateProcessorMetrics(telemetrySink),
	}
}

// Transition is the main function for processing a state transition.
func (sp *StateProcessor[
	ContextT, _,
]) Transition(
	ctx ContextT,
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
) (transition.ValidatorUpdates, error) {
	if blk.IsNil() {
		return nil, nil
	}

	// Process the slots.
	validatorUpdates, err := sp.ProcessSlots(st, blk.GetSlot())
	if err != nil {
		return nil, err
	}

	// Process the block.
	if err = sp.ProcessBlock(ctx, st, blk); err != nil {
		return nil, err
	}

	return validatorUpdates, nil
}

func (sp *StateProcessor[
	_, _,
]) ProcessSlots(
	st *state.StateDB, slot math.Slot,
) (transition.ValidatorUpdates, error) {
	var res transition.ValidatorUpdates

	stateSlot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	// Iterate until we are "caught up".
	for ; stateSlot < slot; stateSlot++ {
		if err = sp.processSlot(st); err != nil {
			return nil, err
		}

		// Process the Epoch Boundary.
		boundary := (stateSlot.Unwrap()+1)%sp.cs.SlotsPerEpoch() == 0
		if boundary {
			var epochUpdates transition.ValidatorUpdates
			if epochUpdates, err = sp.processEpoch(st); err != nil {
				return nil, err
			}
			res = append(res, epochUpdates...)
		}

		// We update on the state because we need to
		// update the state for calls within processSlot/Epoch().
		if err = st.SetSlot(stateSlot + 1); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// processSlot is run when a slot is missed.
func (sp *StateProcessor[
	_, _,
]) processSlot(
	st *state.StateDB,
) error {
	stateSlot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Before we make any changes, we calculate the previous state root.
	prevStateRoot := st.HashTreeRoot()
	if err = st.UpdateStateRootAtIndex(
		stateSlot.Unwrap()%sp.cs.SlotsPerHistoricalRoot(), prevStateRoot,
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
	if (latestHeader.GetStateRoot() == common.Root{}) {
		latestHeader.SetStateRoot(prevStateRoot)
		if err = st.SetLatestBlockHeader(latestHeader); err != nil {
			return err
		}
	}

	// We update the block root.
	return st.UpdateBlockRootAtIndex(
		stateSlot.Unwrap()%sp.cs.SlotsPerHistoricalRoot(),
		latestHeader.HashTreeRoot(),
	)
}

// ProcessBlock processes the block, it optionally verifies the
// state root.
func (sp *StateProcessor[
	ContextT, _,
]) ProcessBlock(
	ctx ContextT,
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	if err := sp.processBlockHeader(ctx, st, blk); err != nil {
		return err
	}

	if err := sp.processExecutionPayload(ctx, st, blk); err != nil {
		return err
	}

	if err := sp.processWithdrawals(st, blk); err != nil {
		return err
	}

	if err := sp.processRandaoReveal(ctx, st, blk); err != nil {
		return err
	}

	if err := sp.processOperations(st, blk); err != nil {
		return err
	}

	// If we are skipping validate, we can skip calculating the state
	// root to save compute.
	if ctx.GetSkipValidateResult() {
		return nil
	}

	// Ensure the calculated state root matches the state root on
	// the block.
	stateRoot := st.HashTreeRoot()
	if blk.GetStateRoot() != stateRoot {
		return errors.Wrapf(
			ErrStateRootMismatch, "expected %s, got %s",
			stateRoot, blk.GetStateRoot(),
		)
	}

	return nil
}

// processEpoch processes the epoch and ensures it matches the local state.
func (sp *StateProcessor[
	_, _,
]) processEpoch(
	st *state.StateDB,
) (transition.ValidatorUpdates, error) {
	slot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	// track validators set before updating it, to be able to
	// inform consensus of the validators set changes
	currentEpoch := sp.cs.SlotToEpoch(slot)
	currentActiveVals, err := getActiveVals(sp.cs, st, currentEpoch)
	if err != nil {
		return nil, err
	}

	if err = sp.processRewardsAndPenalties(st); err != nil {
		return nil, err
	}
	if err = sp.processRegistryUpdates(st); err != nil {
		return nil, err
	}
	if err = sp.processEffectiveBalanceUpdates(st, slot); err != nil {
		return nil, err
	}
	if err = sp.processSlashingsReset(st); err != nil {
		return nil, err
	}
	if err = sp.processRandaoMixesReset(st); err != nil {
		return nil, err
	}

	// only after we have fully updated validators, we enforce
	// a cap on the validators set
	if err = sp.processValidatorSetCap(st); err != nil {
		return nil, err
	}

	// finally compute diffs in validator set to duly update consensus
	nextEpoch := currentEpoch + 1
	nextActiveVals, err := getActiveVals(sp.cs, st, nextEpoch)
	if err != nil {
		return nil, err
	}

	return validatorSetsDiffs(currentActiveVals, nextActiveVals), nil
}

// processBlockHeader processes the header and ensures it matches the local
// state.
func (sp *StateProcessor[
	ContextT, _,
]) processBlockHeader(
	ctx ContextT,
	st *state.StateDB,
	blk *ctypes.BeaconBlock,
) error {
	// Ensure the block slot matches the state slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	if blk.GetSlot() != slot {
		return errors.Wrapf(
			ErrSlotMismatch, "expected: %d, got: %d",
			slot, blk.GetSlot(),
		)
	}

	// Verify that the block is newer than latest block header
	latestBlockHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}
	if blk.GetSlot() <= latestBlockHeader.GetSlot() {
		return errors.Wrapf(
			ErrBlockSlotTooLow, "expected: > %d, got: %d",
			latestBlockHeader.GetSlot(), blk.GetSlot(),
		)
	}

	// Verify that proposer matches with what consensus declares as proposer
	proposer, err := st.ValidatorByIndex(blk.GetProposerIndex())
	if err != nil {
		return err
	}
	stateProposerAddress, err := sp.fGetAddressFromPubKey(proposer.GetPubkey())
	if err != nil {
		return err
	}
	if !bytes.Equal(stateProposerAddress, ctx.GetProposerAddress()) {
		return errors.Wrapf(
			ErrProposerMismatch, "store key: %s, consensus key: %s",
			stateProposerAddress, ctx.GetProposerAddress(),
		)
	}

	// Verify that the parent matches
	parentBlockRoot := latestBlockHeader.HashTreeRoot()
	if parentBlockRoot != blk.GetParentBlockRoot() {
		return errors.Wrapf(
			ErrParentRootMismatch, "expected: %s, got: %s",
			parentBlockRoot.String(), blk.GetParentBlockRoot().String(),
		)
	}

	// Verify proposer is not slashed
	if proposer.IsSlashed() {
		return errors.Wrapf(
			ErrSlashedProposer, "index: %d",
			blk.GetProposerIndex(),
		)
	}

	// Cache current block as the new latest block
	bodyRoot := blk.GetBody().HashTreeRoot()
	var lbh *ctypes.BeaconBlockHeader
	lbh = lbh.New(
		blk.GetSlot(),
		blk.GetProposerIndex(),
		blk.GetParentBlockRoot(),
		// state_root is zeroed and overwritten
		// in the next `process_slot` call.
		common.Root{},
		bodyRoot,
	)
	return st.SetLatestBlockHeader(lbh)
}

// processEffectiveBalanceUpdates as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#effective-balances-updates
func (sp *StateProcessor[
	_, _,
]) processEffectiveBalanceUpdates(
	st *state.StateDB,
	slot math.Slot,
) error {
	// Update effective balances with hysteresis
	validators, err := st.GetValidators()
	if err != nil {
		return err
	}

	var (
		hysteresisIncrement = sp.cs.EffectiveBalanceIncrement() / sp.cs.HysteresisQuotient()
		downwardThreshold   = math.Gwei(
			hysteresisIncrement * sp.cs.HysteresisDownwardMultiplier(),
		)
		upwardThreshold = math.Gwei(
			hysteresisIncrement * sp.cs.HysteresisUpwardMultiplier(),
		)

		idx     math.U64
		balance math.Gwei
	)

	for _, val := range validators {
		idx, err = st.ValidatorIndexByPubkey(val.GetPubkey())
		if err != nil {
			return err
		}

		balance, err = st.GetBalance(idx)
		if err != nil {
			return err
		}

		if balance+downwardThreshold < val.GetEffectiveBalance() ||
			val.GetEffectiveBalance()+upwardThreshold < balance {
			updatedBalance := ctypes.ComputeEffectiveBalance(
				balance,
				math.U64(sp.cs.EffectiveBalanceIncrement()),
				math.U64(sp.cs.MaxEffectiveBalance(
					state.IsPostFork3(sp.cs.DepositEth1ChainID(), slot),
				)),
			)
			val.SetEffectiveBalance(updatedBalance)
			if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
				return err
			}
		}
	}
	return nil
}
