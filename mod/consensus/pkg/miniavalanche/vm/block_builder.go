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

package vm

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/snow/engine/common"
	"go.uber.org/zap"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// how often VM should ping consensus to try and build blocks.
// ProposerVM is active, hence it will eventually decide when BuildBlock is really called
const pingInterval = time.Second

type blockBuilder struct {
	shutdown chan struct{}
	toEngine chan<- common.Message
	vm       *VM
}

func newBlockBuilder(toEngine chan<- common.Message, vm *VM) *blockBuilder {
	bb := &blockBuilder{
		shutdown: make(chan struct{}),
		toEngine: toEngine,
		vm:       vm,
	}
	return bb
}

func (bb *blockBuilder) Start() {
	go bb.listen()
}

func (bb *blockBuilder) Shutdown() {
	close(bb.shutdown)
}

func (bb *blockBuilder) BuildBlock(ctx context.Context) (*StatefulBlock, error) {
	// STEP 1: retrieve parent block data
	parentBlkID := bb.vm.preferredBlkID
	parentBlk, err := bb.vm.getBlock(parentBlkID)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving preferred block, ID: %v: %w", bb.vm.preferredBlkID, err)
	}

	// STEP 2: generate block content
	childBlkHeight := parentBlk.Height() + 1
	childChainTime := parentBlk.Timestamp()
	slotData := &miniavalanche.SlotDataT{
		Slot: math.Slot(childBlkHeight),
	}
	blkBytes, blobBytes, err := bb.vm.middleware.BuildBlock(ctx, slotData)
	if err != nil {
		return nil, fmt.Errorf("failed building block or blob: %w", err)
	}

	// STEP 3: finally build the block
	b, err := block.NewStatelessBlock(
		parentBlkID,
		childBlkHeight,
		childChainTime,
		block.BlockContent{
			BeaconBlockByte: blkBytes,
			BlobsBytes:      blobBytes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed building stateless block: %w", err)
	}

	bb.vm.chainCtx.Log.Info(
		"built block",
		zap.Stringer("blkID", b.ID()),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	return &StatefulBlock{
		StatelessBlock: b,
		vm:             bb.vm,
	}, nil
}

func (bb *blockBuilder) listen() {
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case <-ticker.C:
			ticker.Stop()

			// ping engine in case we are ready to build a block
			bb.toEngine <- common.PendingTxs
		case <-bb.shutdown:
			ticker.Stop()
			return
		}
	}
}
