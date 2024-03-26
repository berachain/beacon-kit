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
	"github.com/berachain/beacon-kit/config/params"
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
	if blk == nil || blk.IsNil() {
		return types.ErrNilBlk
	}

	body := blk.GetBody()
	if body == nil || body.IsNil() {
		return types.ErrNilBlkBody
	}

	// Ensure the block slot matches the state slot.
	if blk.GetSlot() != st.GetSlot() {
		return fmt.Errorf(
			"slot does not match, expected: %d, got: %d",
			st.GetSlot(),
			blk.GetSlot(),
		)
	}

	if deposits := body.GetDeposits(); uint64(
		len(deposits),
	) > bv.cfg.MaxDepositsPerBlock {
		return fmt.Errorf(
			"too many deposits, expected: %d, got: %d",
			bv.cfg.MaxDepositsPerBlock, len(deposits),
		)
	}

	// Ensure the parent block root matches what we have locally.
	parentBlockRoot, err := st.GetBlockRootAtIndex(
		(blk.GetSlot() - 1) % bv.cfg.SlotsPerHistoricalRoot)
	if err != nil {
		return err
	}

	if parentBlockRoot != blk.GetParentBlockRoot() {
		return fmt.Errorf(
			"parent root does not match, expected: %x, got: %x",
			parentBlockRoot,
			blk.GetParentBlockRoot(),
		)
	}
	return nil
}
