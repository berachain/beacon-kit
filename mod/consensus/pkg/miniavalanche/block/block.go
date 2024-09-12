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

package block

import (
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/hashing"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/encoding"
)

type StatelessBlock struct {
	ParentID   ids.ID
	BlkHeight  uint64
	ChainTime  time.Time
	BlkContent Content

	// in memory attributes, unexported since they won't be serialized
	bytes []byte
	blkID ids.ID
}

type Content struct {
	GenesisContent  []byte
	BeaconBlockByte []byte
	BlobsBytes      []byte
}

func NewStatelessBlock(
	parentID ids.ID,
	blkHeight uint64,
	chainTime time.Time,
	blkContent Content,
) (*StatelessBlock, error) {
	b := &StatelessBlock{
		ParentID:   parentID,
		BlkHeight:  blkHeight,
		ChainTime:  chainTime,
		BlkContent: blkContent,
	}

	// setup block bytes and id
	bytes, err := encoding.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed encoding block: %w", err)
	}
	b.initBlkID(bytes)
	return b, nil
}

func ParseStatelessBlock(blockBytes []byte) (*StatelessBlock, error) {
	blk := &StatelessBlock{}
	if err := encoding.Decode(blockBytes, &blk); err != nil {
		return nil, fmt.Errorf("unable to parse block: %w", err)
	}
	blk.initBlkID(blockBytes)
	return blk, nil
}

func (b *StatelessBlock) ID() ids.ID {
	return b.blkID
}

func (b *StatelessBlock) Parent() ids.ID {
	return b.ParentID
}

func (b *StatelessBlock) Bytes() []byte {
	return b.bytes
}

func (b *StatelessBlock) Height() uint64 {
	return b.BlkHeight
}

func (b *StatelessBlock) Timestamp() time.Time {
	return b.ChainTime
}

func (b *StatelessBlock) initBlkID(blkBytes []byte) {
	b.bytes = make([]byte, len(blkBytes))
	copy(b.bytes, blkBytes)
	b.blkID = hashing.ComputeHash256Array(b.bytes)
}
