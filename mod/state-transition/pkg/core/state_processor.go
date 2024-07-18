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
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor[
	BeaconBlockT BeaconBlock[
		DepositT, BeaconBlockBodyT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT,
		BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, WithdrawalT,
	],
	ContextT Context,
	DepositT Deposit[ForkDataT, WithdrawalCredentialsT],
	Eth1DataT interface {
		New(common.Root, math.U64, gethprimitives.ExecutionHash) Eth1DataT
		GetDepositCount() math.U64
	},
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkT interface {
		New(common.Version, common.Version, math.Epoch) ForkT
	},
	ForkDataT ForkData[ForkDataT],
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT ~[32]byte,
] struct {
	// cs is the chain specification for the beacon chain.
	cs common.ChainSpec
	// signer is the BLS signer used for cryptographic operations.
	signer crypto.BLSSigner
	// executionEngine is the engine responsible for executing transactions.
	executionEngine ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	]
	// txsMerkleizer is the merkleizer used for calculating transaction roots.
	txsMerkleizer *merkle.Merkleizer[[32]byte, common.Root]
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor[
	BeaconBlockT BeaconBlock[
		DepositT, BeaconBlockBodyT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT,
		DepositT, ExecutionPayloadT,
		ExecutionPayloadHeaderT,
		WithdrawalT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		KVStoreT, ValidatorT, WithdrawalT,
	],
	ContextT Context,
	DepositT Deposit[ForkDataT, WithdrawalCredentialsT],
	Eth1DataT interface {
		New(common.Root, math.U64, gethprimitives.ExecutionHash) Eth1DataT
		GetDepositCount() math.U64
	},
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkT interface {
		New(common.Version, common.Version, math.Epoch) ForkT
	},
	ForkDataT ForkData[ForkDataT],
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT ~[32]byte,
](
	cs common.ChainSpec,
	executionEngine ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	signer crypto.BLSSigner,
) *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, ContextT, DepositT, Eth1DataT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
	WithdrawalT, WithdrawalCredentialsT,
] {
	return &StateProcessor[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, ContextT, DepositT, Eth1DataT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
		WithdrawalT, WithdrawalCredentialsT,
	]{
		cs:              cs,
		executionEngine: executionEngine,
		signer:          signer,
		txsMerkleizer:   merkle.NewMerkleizer[[32]byte, common.Root](),
	}
}

// Transition is the main function for processing a state transition.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT, _, _, _, _, _, _, _, _, _, _,
]) Transition(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
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

	// We only want to persist state changes if we successfully
	// processed the block.
	st.Save()
	return validatorUpdates, nil
}

func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) ProcessSlots(
	st BeaconStateT, slot math.U64,
) (transition.ValidatorUpdates, error) {
	var (
		validatorUpdates      transition.ValidatorUpdates
		epochValidatorUpdates transition.ValidatorUpdates
	)

	stateSlot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	// Iterate until we are "caught up".
	for ; stateSlot < slot; stateSlot++ {
		// Process the slot
		if err = sp.processSlot(st); err != nil {
			return nil, err
		}

		// Process the Epoch Boundary.
		if uint64(stateSlot+1)%sp.cs.SlotsPerEpoch() == 0 {
			if epochValidatorUpdates, err =
				sp.processEpoch(st); err != nil {
				return nil, err
			}
			validatorUpdates = append(
				validatorUpdates,
				epochValidatorUpdates...,
			)
		}

		// We update on the state because we need to
		// update the state for calls within processSlot/Epoch().
		if err = st.SetSlot(stateSlot + 1); err != nil {
			return nil, err
		}
	}

	return validatorUpdates, nil
}

// processSlot is run when a slot is missed.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processSlot(
	st BeaconStateT,
) error {
	stateSlot, err := st.GetSlot()
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
		uint64(stateSlot)%sp.cs.SlotsPerHistoricalRoot(), prevStateRoot,
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
	var prevBlockRoot common.Root
	prevBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	if err = st.UpdateBlockRootAtIndex(
		uint64(stateSlot)%sp.cs.SlotsPerHistoricalRoot(), prevBlockRoot,
	); err != nil {
		return err
	}

	return nil
}

// ProcessBlock processes the block, it optionally verifies the
// state root.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT, _, _, _, _, _, _, _, _, _, _,
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

	// TODO:
	//
	// phase0.ProcessProposerSlashings
	// phase0.ProcessAttesterSlashings

	// process the randao reveal.
	if err := sp.processRandaoReveal(
		st, blk, ctx.GetSkipValidateRandao(),
	); err != nil {
		return err
	}

	// TODO:
	//
	// phase0.ProcessEth1Vote

	// process the deposits and ensure they match the local state.
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
	if stateRoot, err := st.HashTreeRoot(); err != nil {
		return err
	} else if blk.GetStateRoot() != stateRoot {
		return errors.Wrapf(
			ErrStateRootMismatch, "expected %s, got %s",
			common.Root(stateRoot), blk.GetStateRoot(),
		)
	}

	return nil
}

// processEpoch processes the epoch and ensures it matches the local state.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
]) processEpoch(
	st BeaconStateT,
) (transition.ValidatorUpdates, error) {
	if err := sp.processRewardsAndPenalties(st); err != nil {
		return nil, err
	} else if err = sp.processSlashingsReset(st); err != nil {
		return nil, err
	} else if err = sp.processRandaoMixesReset(st); err != nil {
		return nil, err
	}
	return sp.processSyncCommitteeUpdates(st)
}

// processBlockHeader processes the header and ensures it matches the local
// state.
func (sp *StateProcessor[
	BeaconBlockT, _, BeaconBlockHeaderT, BeaconStateT,
	_, _, _, _, _, _, _, _, ValidatorT, _, _,
]) processBlockHeader(
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	var (
		slot              math.Slot
		err               error
		latestBlockHeader BeaconBlockHeaderT
		parentBlockRoot   common.Root
		bodyRoot          common.Root
		proposer          ValidatorT
	)

	// Ensure the block slot matches the state slot.
	if slot, err = st.GetSlot(); err != nil {
		return err
	} else if blk.GetSlot() != slot {
		return errors.Wrapf(
			ErrSlotMismatch,
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
			"expected: %s, got: %s",
			parentBlockRoot.String(), blk.GetParentBlockRoot().String(),
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
	var lbh BeaconBlockHeaderT
	if bodyRoot, err = blk.GetBody().HashTreeRoot(); err != nil {
		return err
	} else if err = st.SetLatestBlockHeader(
		lbh.New(
			blk.GetSlot(),
			blk.GetProposerIndex(),
			blk.GetParentBlockRoot(),
			// state_root is zeroed and overwritten
			// in the next `process_slot` call.
			common.Root{},
			bodyRoot,
		),
	); err != nil {
		return err
	}

	// Check to make sure the proposer isn't slashed.
	if proposer, err = st.ValidatorByIndex(blk.GetProposerIndex()); err != nil {
		return err
	} else if proposer.IsSlashed() {
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
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
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
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _,
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
