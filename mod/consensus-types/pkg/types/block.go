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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// BeaconBlock is the interface for a beacon block.
type BeaconBlock struct {
	RawBeaconBlock[*BeaconBlockBody]
}

// Empty creates an empty beacon block.
func (w *BeaconBlock) Empty(forkVersion uint32) *BeaconBlock {
	switch forkVersion {
	case version.Deneb:
		return &BeaconBlock{
			RawBeaconBlock: (*BeaconBlockDeneb)(nil),
		}
	default:
		panic("fork version not supported")
	}
}

// NewWithVersion assembles a new beacon block from the given.
func (w *BeaconBlock) NewWithVersion(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	forkVersion uint32,
) (*BeaconBlock, error) {
	var (
		block RawBeaconBlock[*BeaconBlockBody]
		base  = BeaconBlockHeaderBase{
			Slot:            slot.Unwrap(),
			ProposerIndex:   proposerIndex.Unwrap(),
			ParentBlockRoot: parentBlockRoot,
			StateRoot:       bytes.B32{},
		}
	)

	switch forkVersion {
	case version.Deneb:
		block = &BeaconBlockDeneb{
			BeaconBlockHeaderBase: base,
			Body:                  &BeaconBlockBodyDeneb{},
		}
	default:
		return &BeaconBlock{}, ErrForkVersionNotSupported
	}

	return &BeaconBlock{
		RawBeaconBlock: block,
	}, nil
}

// UnmarshalSSZWithVersion creates a new beacon block from the given SSZ bytes.
func (w *BeaconBlock) UnmarshalSSZWithVersion(
	bz []byte,
	forkVersion uint32,
) (*BeaconBlock, error) {
	var block = new(BeaconBlock)
	switch forkVersion {
	case version.Deneb:
		block.RawBeaconBlock = &BeaconBlockDeneb{}
	default:
		return block, ErrForkVersionNotSupported
	}

	if err := block.UnmarshalSSZ(bz); err != nil {
		return block, err
	}
	return block, nil
}

// IsNil checks if the beacon block is nil.
func (w *BeaconBlock) IsNil() bool {
	return w == nil ||
		w.RawBeaconBlock == nil ||
		w.RawBeaconBlock.IsNil()
}

// BeaconBlockDeneb represents a block in the beacon chain during
// the Deneb fork.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path block.go -objs BeaconBlockDeneb -include ../../../primitives/pkg/common,../../../primitives/pkg/crypto,../../../primitives/pkg/math,..,./header.go,./withdrawal_credentials.go,../../../engine-primitives/pkg/engine-primitives/withdrawal.go,./deposit.go,./payload.go,./deposit.go,../../../primitives/pkg/eip4844,../../../primitives/pkg/bytes,./eth1data.go,../../../primitives/pkg/math,../../../primitives/pkg/common,./body.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output block.ssz.go
type BeaconBlockDeneb struct {
	// BeaconBlockHeaderBase is the base of the BeaconBlockDeneb.
	BeaconBlockHeaderBase
	// Body is the body of the BeaconBlockDeneb, containing the block's
	// operations.
	Body *BeaconBlockBodyDeneb
}

// Version identifies the version of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) Version() uint32 {
	return version.Deneb
}

// IsNil checks if the BeaconBlockDeneb instance is nil.
func (b *BeaconBlockDeneb) IsNil() bool {
	return b == nil
}

// SetStateRoot sets the state root of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) SetStateRoot(root common.Root) {
	b.StateRoot = root
}

// GetBody retrieves the body of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) GetBody() *BeaconBlockBody {
	return &BeaconBlockBody{RawBeaconBlockBody: b.Body}
}

// GetHeader builds a BeaconBlockHeader from the BeaconBlockDeneb.
func (b BeaconBlockDeneb) GetHeader() *BeaconBlockHeader {
	bodyRoot, err := b.GetBody().HashTreeRoot()
	if err != nil {
		return nil
	}

	return &BeaconBlockHeader{
		BeaconBlockHeaderBase: BeaconBlockHeaderBase{
			Slot:            b.Slot,
			ProposerIndex:   b.ProposerIndex,
			ParentBlockRoot: b.ParentBlockRoot,
			StateRoot:       b.StateRoot,
		},
		BodyRoot: bodyRoot,
	}
}
