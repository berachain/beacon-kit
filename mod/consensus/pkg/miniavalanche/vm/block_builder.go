package vm

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/snow/engine/common"
	"go.uber.org/zap"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
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
	var parentStatelessBlk *block.StatelessBlock
	parentBlk, found := bb.vm.verifiedBlocks[parentBlkID]
	if found {
		parentStatelessBlk = parentBlk.StatelessBlock
	} else {
		var err error
		parentStatelessBlk, err = bb.vm.state.GetBlock(parentBlkID)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving preferred block, ID: %v: %w", bb.vm.preferredBlkID, err)
		}
	}
	parentBlkHeight := parentStatelessBlk.Height()

	// STEP 2: generate block content
	childBlkHeight := parentBlkHeight + 1
	childChainTime := parentStatelessBlk.Timestamp()

	slotData := &miniavalanche.SlotDataT{} // TODO: FILL IT UP
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
