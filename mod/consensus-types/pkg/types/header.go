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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockHeaderBase represents the base of a beacon block header.
type BeaconBlockHeaderBase struct {
	// Slot represents the position of the block in the chain.
	// TODO: Put back to math.Slot after fastssz fixes.
	Slot uint64
	// ProposerIndex is the index of the validator who proposed the block.
	// TODO: Put back to math.ProposerIndex after fastssz fixes.
	ProposerIndex uint64
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetSlot() math.Slot {
	return math.Slot(b.Slot)
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetProposerIndex() math.ValidatorIndex {
	return math.ValidatorIndex(b.ProposerIndex)
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlockDeneb.
func (b *BeaconBlockHeaderBase) GetStateRoot() common.Root {
	return b.StateRoot
}

// BeaconBlockHeader is the header of a beacon block.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path header.go -objs BeaconBlockHeader -include ../../../primitives/pkg/common,../../../primitives/pkg/math,../../../primitives/pkg/bytes,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output header.ssz.go
type BeaconBlockHeader struct {
	// BeaconBlockHeaderBase is the base of the block.
	BeaconBlockHeaderBase
	// 	// BodyRoot is the root of the block body.
	BodyRoot common.Root `json:"bodyRoot"`
}

// NewBeaconBlockHeader creates a new BeaconBlockHeader.
func NewBeaconBlockHeader(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	stateRoot common.Root,
	bodyRoot common.Root,
) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		BeaconBlockHeaderBase: BeaconBlockHeaderBase{
			Slot:            uint64(slot),
			ProposerIndex:   uint64(proposerIndex),
			ParentBlockRoot: parentBlockRoot,
			StateRoot:       stateRoot,
		},
		BodyRoot: bodyRoot,
	}
}
