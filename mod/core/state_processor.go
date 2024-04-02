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
	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor struct {
	cfg    *params.BeaconChainConfig
	bp     BlobsProcessor
	rp     RandaoProcessor
	logger log.Logger
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor(
	cfg *params.BeaconChainConfig,
	bp BlobsProcessor,
	rp RandaoProcessor,
	logger log.Logger,
) *StateProcessor {
	return &StateProcessor{
		cfg:    cfg,
		bp:     bp,
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
		uint64(slot)%sp.cfg.SlotsPerHistoricalRoot,
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
		uint64(slot)%sp.cfg.SlotsPerHistoricalRoot, prevBlockRoot,
	); err != nil {
		return err
	}

	// Process the Epoch Boundary.
	if uint64(slot+1)%sp.cfg.SlotsPerEpoch == 0 {
		if err = sp.processEpoch(st); err != nil {
			return err
		}
		sp.logger.Info(
			"processed epoch transition ⏰ ",
			"old", uint64(slot)/sp.cfg.SlotsPerEpoch,
			"new", uint64(slot+1)/sp.cfg.SlotsPerEpoch,
		)
	}

	return st.SetSlot(slot + 1)
}

// ProcessBlobs processes the blobs and ensures they match the local state.
func (sp *StateProcessor) ProcessBlobs(
	avs state.AvailabilityStore,
	blk types.BeaconBlock,
	sidecars *datypes.BlobSidecars,
) error {
	return sp.bp.ProcessBlobs(avs, blk, sidecars)
}

// ProcessBlock processes the block and ensures it matches the local state.
func (sp *StateProcessor) ProcessBlock(
	st state.BeaconState,
	blk types.BeaconBlock,
) error {
	header, err := types.NewBeaconBlockHeader(blk)
	if err != nil {
		return err
	}

	// process the freshly created header.
	if err = sp.processHeader(st, header); err != nil {
		return err
	}

	// process the withdrawals.
	body := blk.GetBody()
	if err = sp.processWithdrawals(
		st, body.GetExecutionPayload().GetWithdrawals(),
	); err != nil {
		return err
	}

	// phase0.ProcessProposerSlashings
	// phase0.ProcessAttesterSlashings

	// process the randao reveal.
	if err = sp.processRandaoReveal(st, blk); err != nil {
		return err
	}

	// phase0.ProcessEth1Vote ? forkchoice?

	// process the deposits and ensure they match the local state.
	if err = sp.processOperations(st, body); err != nil {
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
	header *types.BeaconBlockHeader,
) error {
	// Store as the new latest block
	headerRaw := &types.BeaconBlockHeader{
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
	// if len(body.GetDeposits()) == min(0, len(body.GetDeposits())) {
	return sp.processDeposits(st, body.GetDeposits())
}

// ProcessDeposits processes the deposits and ensures they match the
// local state.
func (sp *StateProcessor) processDeposits(
	st state.BeaconState,
	deposits []*types.Deposit,
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
	dep *types.Deposit,
) {
	idx, err := st.ValidatorIndexByPubkey(dep.Pubkey[:])
	if err != nil {
		_ = 0
		// # Verify the deposit signature (proof of possession) which is not
		// checked by the deposit contract
		// deposit_message = DepositMessage(
		//     pubkey=pubkey,
		//     withdrawal_credentials=withdrawal_credentials,
		//     amount=amount,
		// )
		// domain = compute_domain(DOMAIN_DEPOSIT)  # Fork-agnostic domain since
		// deposits are valid across forks
		// signing_root = compute_signing_root(deposit_message, domain)
		// if bls.Verify(pubkey, signing_root, signature):
		// add_validator_to_registry(state, pubkey, withdrawal_credentials,
		// amount)
	} else {
		var val *types.Validator
		val, err = st.ValidatorByIndex(idx)
		if err != nil {
			return
		}

		// TODO: Modify balance here and then effective balance once per epoch.
		val.EffectiveBalance = min(val.EffectiveBalance+dep.Amount,
			primitives.Gwei(sp.cfg.MaxEffectiveBalance))
		if err = st.UpdateValidatorAtIndex(idx, val); err != nil {
			return
		}
	}
}

// processWithdrawals processes the withdrawals and ensures they match the
// local state.
func (sp *StateProcessor) processWithdrawals(
	st state.BeaconState,
	withdrawals []*primitives.Withdrawal,
) error {
	// Dequeue and verify the withdrawals.
	localWithdrawals, err := st.DequeueWithdrawals(uint64(len(withdrawals)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, wd := range withdrawals {
		if wd == nil {
			return types.ErrNilWithdrawal
		}
		if wd.Index != localWithdrawals[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localWithdrawals[i].Index, wd.Index)
		}

		var val *types.Validator
		val, err = st.ValidatorByIndex(wd.Validator)
		if err != nil {
			continue
		}

		// TODO: Modify balance here and then effective balance once per epoch.
		val.EffectiveBalance -= min(
			val.EffectiveBalance, wd.Amount,
		)
		if err = st.UpdateValidatorAtIndex(wd.Validator, val); err != nil {
			return err
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
	epoch, err := st.GetCurrentEpoch(sp.cfg.SlotsPerEpoch)
	if err != nil {
		return err
	}

	index := (uint64(epoch) + 1) % sp.cfg.EpochsPerSlashingsVector
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
	slotsPerEpoch := sp.cfg.SlotsPerEpoch
	totalBalance, err := st.GetTotalActiveBalances(slotsPerEpoch)
	if err != nil {
		return err
	}

	totalSlashings, err := st.GetTotalSlashing()
	if err != nil {
		return err
	}
	proportionalSlashingMultiplier := sp.cfg.ProportionalSlashingMultiplier
	adjustedTotalSlashingBalance := min(
		uint64(totalSlashings)*proportionalSlashingMultiplier,
		uint64(totalBalance),
	)
	vals, err := st.GetValidators()
	if err != nil {
		return err
	}

	// Get the current epoch
	epoch, err := st.GetCurrentEpoch(slotsPerEpoch)
	if err != nil {
		return err
	}

	// Iterate through the validators.
	for _, val := range vals {
		// Checks if the validator is slashable.
		//nolint:gomnd // this is in the spec
		slashableEpoch := (uint64(epoch) + sp.cfg.EpochsPerSlashingsVector) / 2
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
	increment := sp.cfg.EffectiveBalanceIncrement
	balDivIncrement := uint64(val.EffectiveBalance) / increment
	penaltyNumerator := balDivIncrement * adjustedTotalSlashingBalance
	penalty := penaltyNumerator / totalBalance * increment

	// Get the val index and decrease the balance of the validator.
	idx, err := st.ValidatorIndexByPubkey(val.Pubkey[:])
	if err != nil {
		return err
	}

	return st.DecreaseBalance(idx, primitives.Gwei(penalty))
}
