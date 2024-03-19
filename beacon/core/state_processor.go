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
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cockroachdb/errors"
	"github.com/sourcegraph/conc/iter"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor struct {
	cfg *config.Beacon
	rp  RandaoProcessor
	vsu ValsetUpdater
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor(
	cfg *config.Beacon,
	rp RandaoProcessor,
	vsu ValsetUpdater,
) *StateProcessor {
	return &StateProcessor{
		cfg: cfg,
		rp:  rp,
		vsu: vsu,
	}
}

// ProcessSlot is run when a slot is missed.
func (sp *StateProcessor) ProcessSlot(
	st state.BeaconState,
) error {
	// Before we make any changes, we calculate the previous state root.
	prevStateRoot, err := st.HashTreeRoot()
	if err != nil {
		return err
	}

	// We update our state roots and block roots. Note: we use
	// st.GetSlot() even though technically this was the state root from
	// end of the previous slot.
	if err = st.UpdateStateRootAtIndex(
		st.GetSlot()%sp.cfg.Limits.SlotsPerHistoricalRoot,
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
	if (latestHeader.StateRoot == primitives.HashRoot{}) {
		latestHeader.StateRoot = prevStateRoot
		if err = st.SetLatestBlockHeader(latestHeader); err != nil {
			return err
		}
	}

	// We update the block root.
	var prevBlockRoot primitives.HashRoot
	prevBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	if err = st.UpdateBlockRootAtIndex(
		st.GetSlot()%sp.cfg.Limits.SlotsPerHistoricalRoot, prevBlockRoot,
	); err != nil {
		return err
	}
	return nil
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
	if err = sp.processDeposits(st, body.GetDeposits()); err != nil {
		return err
	}

	// process the redirects and ensure they match the local state.
	if err = sp.processRedirects(st, body.GetRedirects()); err != nil {
		return err
	}

	// ProcessVoluntaryExits

	return nil
}

// ProcessBlob processes a blob.
func (sp *StateProcessor) ProcessBlobs(
	avs state.AvailabilityStore,
	blk types.BeaconBlock,
	sidecars *types.BlobSidecars,
) error {
	// Verify the KZG inclusion proofs.
	bodyRoot, err := blk.GetBody().HashTreeRoot()
	if err != nil {
		return err
	}

	// Ensure the blobs are available.
	if err = errors.Join(iter.Map(
		sidecars.Sidecars,
		func(sidecar **types.BlobSidecar) error {
			if *sidecar == nil {
				return ErrAttemptedToVerifyNilSidecar
			}
			// Store the blobs under a single height.
			return types.VerifyKZGInclusionProof(
				bodyRoot[:], *sidecar, (*sidecar).Index,
			)
		},
	)...); err != nil {
		return err
	}

	return avs.Persist(blk.GetSlot(), sidecars.Sidecars...)
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
		if dep == nil {
			return types.ErrNilDeposit
		}
		if dep.Index != localDeposits[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localDeposits[i].Index, dep.Index)
		}

		// TODO: These changes are not encapsulated in the state root of
		// the beacon store. @po-bera needs for EIP-4788.
		if err = sp.vsu.IncreaseConsensusPower(
			st.Context(),
			dep.Credentials,
			[bls12381.PubKeyLength]byte(dep.Pubkey),
			dep.Amount,
			dep.Signature,
			dep.Index,
		); err != nil {
			return err
		}
	}
	return nil
}

// processRedirects processes the redirects and ensures they match the
// local state.
func (sp *StateProcessor) processRedirects(
	st state.BeaconState,
	redirects []*types.Redirect,
) error {
	// Dequeue and verify the logs.
	localRedirects, err := st.DequeueRedirects(uint64(len(redirects)))
	if err != nil {
		return err
	}

	// Ensure the redirects match the local state.
	for i, red := range redirects {
		if red == nil {
			return types.ErrNilRedirect
		}
		if red.Index != localRedirects[i].Index {
			return fmt.Errorf(
				"redirect index does not match, expected: %d, got: %d",
				localRedirects[i].Index, red.Index)
		}

		// TODO: These changes are not encapsulated in the state root of
		// the beacon store. @po-bera needs for EIP-4788.
		if err = sp.vsu.RedirectConsensusPower(
			st.Context(),
			red.Credentials,
			[bls12381.PubKeyLength]byte(red.Pubkey),
			[bls12381.PubKeyLength]byte(red.NewPubkey),
			red.Amount,
			red.Signature,
			red.Index,
		); err != nil {
			return err
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
	for i, wd := range withdrawals {
		if wd == nil {
			return types.ErrNilWithdrawal
		}
		if wd.Index != localWithdrawals[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localWithdrawals[i].Index, wd.Index)
		}
		// TODO: These changes are not encapsulated in the state root of
		// the beacon store. @po-bera needs for EIP-4788.
		var pk []byte
		pk, err = st.ValidatorPubKeyByIndex(
			wd.Validator,
		)
		if err != nil {
			return err
		}
		if err = sp.vsu.DecreaseConsensusPower(
			st.Context(),
			wd.Address,
			[bls12381.PubKeyLength]byte(pk),
			wd.Amount,
		); err != nil {
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
