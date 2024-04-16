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

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/version"
	"github.com/davecgh/go-spew/spew"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor struct {
	cs     primitives.ChainSpec
	bv     BlobVerifier
	rp     RandaoProcessor
	logger log.Logger
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor(
	cs primitives.ChainSpec,
	bv BlobVerifier,
	rp RandaoProcessor,
	logger log.Logger,
) *StateProcessor {
	return &StateProcessor{
		cs:     cs,
		bv:     bv,
		rp:     rp,
		logger: logger.With("module", "state-processor"),
	}
}

// StateTransition is the main function for processing a state transition.
func (sp *StateProcessor) Transition(
	st state.BeaconState,
	blk types.ReadOnlyBeaconBlock,
	/*validateSignature bool, */
	validateResult bool,
) error {
	// Process the slot.
	if err := sp.ProcessSlot(st); err != nil {
		return err
	}

	// Process the block.
	if err := sp.ProcessBlock(st, blk); err != nil {
		return err
	}

	if validateResult {
		stateRoot, err := st.HashTreeRoot()
		if err != nil {
			return err
		}

		if stateRoot != blk.GetStateRoot() {
			return ErrStateRootMismatch
		}
	}

	return nil
}

// ProcessSlot is run when a slot is missed.
func (sp *StateProcessor) ProcessSlot(
	st state.BeaconState,
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
			"processed epoch transition ⏰ ",
			"old", uint64(slot)/sp.cs.SlotsPerEpoch(),
			"new", uint64(slot+1)/sp.cs.SlotsPerEpoch(),
		)
	}

	return st.SetSlot(slot + 1)
}

// ProcessBlobs processes the blobs and ensures they match the local state.
func (sp *StateProcessor) ProcessBlobs(
	st state.BeaconState,
	avs AvailabilityStore,
	sidecars *datypes.BlobSidecars,
) error {
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// If there are no blobs to verify, return early.
	numBlobs := len(sidecars.Sidecars)
	if numBlobs == 0 {
		sp.logger.Info(
			"no blobs to verify, skipping verifier 🧢",
			"slot",
			slot,
		)
		return nil
	}

	// Otherwise, we run the verification checks on the blobs.
	if err = sp.bv.VerifyBlobs(
		sidecars,
		types.KZGOffset(sp.cs.MaxBlobCommitmentsPerBlock()),
	); err != nil {
		return err
	}

	sp.logger.Info(
		"successfully verified all blob sidecars 💦",
		"num_blobs",
		numBlobs,
		"slot",
		slot,
	)

	// Lastly, we store the blobs in the availability store.
	return avs.Persist(slot, sidecars)
}

// ProcessBlock processes the block and ensures it matches the local state.
func (sp *StateProcessor) ProcessBlock(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	// process the freshly created header.
	if err := sp.processHeader(st, blk); err != nil {
		return err
	}

	// process the withdrawals.
	body := blk.GetBody()
	if err := sp.processWithdrawals(
		st, body.GetExecutionPayload(),
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

	// process the deposits and ensure they match the local state.
	if err := sp.processOperations(st, body); err != nil {
		return err
	}

	// ProcessVoluntaryExits

	return nil
}

// processEpoch processes the epoch and ensures it matches the local state.
func (sp *StateProcessor) processEpoch(st state.BeaconState) error {
	var err error
	if err = sp.processSlashingsReset(st); err != nil {
		return err
	}
	if err = sp.processRandaoMixesReset(st); err != nil {
		return err
	}
	return nil
}

// processHeader processes the header and ensures it matches the local state.
func (sp *StateProcessor) processHeader(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	// TODO: this function is really confusing, can probably just
	// be removed and the logic put in the ProcessBlock function.
	header := blk.GetHeader()
	if header == nil {
		return types.ErrNilBlockHeader
	}

	// Store as the new latest block
	headerRaw := &primitives.BeaconBlockHeader{
		Slot:          header.Slot,
		ProposerIndex: header.ProposerIndex,
		ParentRoot:    header.ParentRoot,
		// state_root is zeroed and overwritten in the next `process_slot` call.
		// with BlockHeaderState.UpdateStateRoot(), once the post state is
		// available.
		StateRoot: [32]byte{},
		BodyRoot:  header.BodyRoot,
	}
	return st.SetLatestBlockHeader(headerRaw)
}

// processOperations processes the operations and ensures they match the
// local state.
func (sp *StateProcessor) processOperations(
	st state.BeaconState,
	body types.BeaconBlockBody,
) error {
	return sp.processDeposits(st, body.GetDeposits())
}

// ProcessDeposits processes the deposits and ensures they match the
// local state.
func (sp *StateProcessor) processDeposits(
	st state.BeaconState,
	deposits []*primitives.Deposit,
) error {
	// Dequeue and verify the logs.
	localDeposits, err := st.DequeueDeposits(uint64(len(deposits)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, dep := range deposits {
		if dep.Index != localDeposits[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localDeposits[i].Index, dep.Index)
		}

		var depIdx uint64
		depIdx, err = st.GetEth1DepositIndex()
		if err != nil {
			return err
		}

		// TODO: this is bad but safe.
		if dep.Index != depIdx {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				depIdx, dep.Index)
		}

		// TODO: this is a shitty spot for this.
		// TODO: deprecate using this.
		if err = st.SetEth1DepositIndex(depIdx + 1); err != nil {
			return err
		}
		sp.processDeposit(st, dep)
	}
	return nil
}

// processDeposit processes the deposit and ensures it matches the local state.
func (sp *StateProcessor) processDeposit(
	st state.BeaconState,
	dep *primitives.Deposit,
) {
	idx, err := st.ValidatorIndexByPubkey(dep.Pubkey)
	// If the validator already exists, we update the balance.
	if err == nil {
		var val *types.Validator
		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return
		}

		// TODO: Modify balance here and then effective balance once per epoch.
		val.EffectiveBalance = min(val.EffectiveBalance+dep.Amount,
			primitives.Gwei(sp.cs.MaxEffectiveBalance()))
		if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
			return
		}
		// Exiting early because we only check signature on creation
		return
	}
	// If the validator does not exist, we add the validator.
	// Add the validator to the registry.
	if err = sp.createValidator(st, dep); err != nil {
		sp.logger.Error("failed to create validator", "error", err)
	}
}

// createValidator creates a validator if the deposit is valid.
func (sp *StateProcessor) createValidator(
	st state.BeaconState,
	dep *primitives.Deposit,
) error {
	var (
		genesisValidatorsRoot primitives.Root
		epoch                 primitives.Epoch
		err                   error
	)

	// Get the genesis validators root to be used to find fork data later.
	genesisValidatorsRoot, err = st.GetGenesisValidatorsRoot()
	if err != nil {
		return err
	}

	// Get the current epoch.
	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}
	epoch = sp.cs.SlotToEpoch(slot)

	// Get the fork data for the current epoch.
	fd := primitives.NewForkData(
		version.FromUint32[primitives.Version](
			sp.cs.ActiveForkVersionForEpoch(epoch),
		), genesisValidatorsRoot,
	)

	depositMessage := primitives.DepositMessage{
		Pubkey:      dep.Pubkey,
		Credentials: dep.Credentials,
		Amount:      dep.Amount,
	}
	if err = depositMessage.VerifyCreateValidator(
		fd, dep.Signature, blst.VerifySignaturePubkeyBytes, sp.cs.DomainTypeDeposit(),
	); err != nil {
		return err
	}

	// Add the validator to the registry.
	return sp.addValidatorToRegistry(st, dep)
}

// addValidatorToRegistry adds a validator to the registry.
func (sp *StateProcessor) addValidatorToRegistry(
	st state.BeaconState,
	dep *primitives.Deposit,
) error {
	val := types.NewValidatorFromDeposit(
		dep.Pubkey,
		dep.Credentials,
		dep.Amount,
		primitives.Gwei(sp.cs.EffectiveBalanceIncrement()),
		primitives.Gwei(sp.cs.MaxEffectiveBalance()),
	)
	if err := st.AddValidator(val); err != nil {
		return err
	}

	idx, err := st.ValidatorIndexByPubkey(val.Pubkey)
	if err != nil {
		return err
	}
	return st.IncreaseBalance(idx, dep.Amount)
}

// processWithdrawals as per the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#new-process_withdrawals
//
//nolint:lll
func (sp *StateProcessor) processWithdrawals(
	st state.BeaconState,
	payload engineprimitives.ExecutionPayload,
) error {
	// Dequeue and verify the logs.
	var nextValidatorIndex primitives.ValidatorIndex
	payloadWithdrawals := payload.GetWithdrawals()
	expectedWithdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		return err
	}

	// Ensure the withdrawals have the same length
	if len(expectedWithdrawals) != len(payloadWithdrawals) {
		return fmt.Errorf(
			"withdrawals do not match expected length %d, got %d",
			len(expectedWithdrawals), len(payloadWithdrawals),
		)
	}

	// Compare and process each withdrawal.
	for i, wd := range expectedWithdrawals {
		// Ensure the withdrawals match the local state.
		if !wd.Equals(payloadWithdrawals[i]) {
			return fmt.Errorf(
				"withdrawals do not match expected %s, got %s",
				spew.Sdump(wd), spew.Sdump(payloadWithdrawals[i]),
			)
		}

		// Then we process the withdrawal.
		if err = st.DecreaseBalance(wd.Validator, wd.Amount); err != nil {
			return err
		}
	}

	// Update the next withdrawal index if this block contained withdrawals
	numWithdrawals := len(expectedWithdrawals)
	if numWithdrawals != 0 {
		// Next sweep starts after the latest withdrawal's validator index
		if err = st.SetNextWithdrawalIndex(
			(expectedWithdrawals[len(expectedWithdrawals)-1].Index + 1).Unwrap(),
		); err != nil {
			return err
		}
	}

	totalValidators, err := st.GetTotalValidators()
	if err != nil {
		return err
	}

	// Update the next validator index to start the next withdrawal sweep
	//#nosec:G701 // won't overflow in practice.
	if numWithdrawals == int(sp.cs.MaxWithdrawalsPerPayload()) {
		// Next sweep starts after the latest withdrawal's validator index
		nextValidatorIndex =
			(expectedWithdrawals[len(expectedWithdrawals)-1].Index + 1) %
				primitives.U64(totalValidators)
	} else {
		// Advance sweep by the max length of the sweep if there was not
		// a full set of withdrawals
		nextValidatorIndex, err = st.GetNextWithdrawalValidatorIndex()
		if err != nil {
			return err
		}
		nextValidatorIndex += primitives.ValidatorIndex(
			sp.cs.MaxValidatorsPerWithdrawalsSweep())
		nextValidatorIndex %= primitives.ValidatorIndex(totalValidators)
	}

	return st.SetNextWithdrawalValidatorIndex(nextValidatorIndex)
}

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor) processRandaoReveal(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	return sp.rp.ProcessRandao(st, blk)
}

// processRandaoMixesReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#randao-mixes-updates
//
//nolint:lll
func (sp *StateProcessor) processRandaoMixesReset(
	st state.BeaconState,
) error {
	return sp.rp.ProcessRandaoMixesReset(st)
}

// processSlashingsReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#slashings-balances-updates
//
//nolint:lll
func (sp *StateProcessor) processSlashingsReset(
	st state.BeaconState,
) error {
	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	index := (uint64(sp.cs.SlotToEpoch(slot)) + 1) % sp.cs.EpochsPerSlashingsVector()
	return st.UpdateSlashingAtIndex(index, 0)
}

// processProposerSlashing as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#proposer-slashings
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor) processProposerSlashing(
	_ state.BeaconState,
	// ps types.ProposerSlashing,
) error {
	return nil
}

// processAttesterSlashing as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#attester-slashings
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor) processAttesterSlashing(
	_ state.BeaconState,
	// as types.AttesterSlashing,
) error {
	return nil
}

// processSlashings as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#slashings
//
// processSlashings processes the slashings and ensures they match the local
// state.
//
//nolint:lll,unused // will be used later
func (sp *StateProcessor) processSlashings(
	st state.BeaconState,
) error {
	totalBalance, err := st.GetTotalActiveBalances(sp.cs.SlotsPerEpoch())
	if err != nil {
		return err
	}

	totalSlashings, err := st.GetTotalSlashing()
	if err != nil {
		return err
	}
	proportionalSlashingMultiplier := sp.cs.ProportionalSlashingMultiplier
	adjustedTotalSlashingBalance := min(
		uint64(totalSlashings)*proportionalSlashingMultiplier(),
		uint64(totalBalance),
	)
	vals, err := st.GetValidators()
	if err != nil {
		return err
	}

	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Iterate through the validators.
	for _, val := range vals {
		// Checks if the validator is slashable.
		//nolint:gomnd // this is in the spec
		slashableEpoch := (uint64(sp.cs.SlotToEpoch(slot)) + sp.cs.EpochsPerSlashingsVector()) / 2
		// If the validator is slashable, and slashed
		if val.Slashed && (slashableEpoch == uint64(val.WithdrawableEpoch)) {
			if err = sp.processSlash(
				st,
				val,
				adjustedTotalSlashingBalance,
				uint64(totalBalance),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

// processSlash handles the logic for slashing a validator.
//
//nolint:unused // will be used later
func (sp *StateProcessor) processSlash(
	st state.BeaconState,
	val *types.Validator,
	adjustedTotalSlashingBalance uint64,
	totalBalance uint64,
) error {
	// Calculate the penalty.
	increment := sp.cs.EffectiveBalanceIncrement()
	balDivIncrement := uint64(val.GetEffectiveBalance()) / increment
	penaltyNumerator := balDivIncrement * adjustedTotalSlashingBalance
	penalty := penaltyNumerator / totalBalance * increment

	// Get the val index and decrease the balance of the validator.
	idx, err := st.ValidatorIndexByPubkey(val.Pubkey)
	if err != nil {
		return err
	}

	return st.DecreaseBalance(idx, primitives.Gwei(penalty))
}
