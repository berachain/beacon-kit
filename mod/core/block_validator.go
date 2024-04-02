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

	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
)

// BlockValidator is responsible for validating incoming
// BeaconBlocks.
type BlockValidator struct {
	cfg *params.BeaconChainConfig
}

// NewBlockValidator creates a new block validator.
func NewBlockValidator(cfg *params.BeaconChainConfig) *BlockValidator {
	return &BlockValidator{
		cfg: cfg,
	}
}

// ValidateBlock validates the incoming block.
func (bv *BlockValidator) ValidateBlock(
	st state.BeaconState,
	blk types.ReadOnlyBeaconBlock,
) error {
	// Get the block body.
	body := blk.GetBody()
	if body == nil || body.IsNil() {
		return types.ErrNilBlkBody
	}

	// Get the current slot.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Ensure the block slot matches the state slot.
	if blk.GetSlot() != slot {
		return fmt.Errorf(
			"slot does not match, expected: %d, got: %d",
			slot,
			blk.GetSlot(),
		)
	}

	// Get the latest block header.
	latestBlockHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}

	// Ensure the block is within the acceptable range.
	if blk.GetSlot() <= latestBlockHeader.Slot {
		return fmt.Errorf(
			"block slot is too low, expected: > %d, got: %d",
			latestBlockHeader.Slot,
			blk.GetSlot(),
		)
	}

	// Ensure the block is within the acceptable range.
	// TODO: move this is in the wrong spot.
	if deposits := body.GetDeposits(); uint64(
		len(deposits),
	) > bv.cfg.MaxDepositsPerBlock {
		return fmt.Errorf(
			"too many deposits, expected: %d, got: %d",
			bv.cfg.MaxDepositsPerBlock, len(deposits),
		)
	}

	// Ensure the parent root matches the latest block header.
	parentBlockRoot, err := latestBlockHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	// Ensure the parent root matches the latest block header.
	if parentBlockRoot != blk.GetParentBlockRoot() {
		return fmt.Errorf(
			"parent root does not match, expected: %x, got: %x",
			parentBlockRoot,
			blk.GetParentBlockRoot(),
		)
	}
	return nil
}
