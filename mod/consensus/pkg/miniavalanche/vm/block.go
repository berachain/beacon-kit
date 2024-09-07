package vm

import (
	"context"

	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"go.uber.org/zap"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
)

var _ snowman.Block = (*StatefulBlock)(nil)

type StatefulBlock struct {
	*block.StatelessBlock
	vm *VM
}

func (b *StatefulBlock) Verify(ctx context.Context) error {
	if err := b.vm.middleware.VerifyBlock(ctx, b.StatelessBlock); err != nil {
		return err
	}

	b.vm.verifiedBlocks[b.ID()] = b
	return nil
}

func (b *StatefulBlock) Accept(ctx context.Context) error {
	delete(b.vm.verifiedBlocks, b.ID())

	b.vm.state.SetLastAccepted(b.ID())
	b.vm.state.AddStatelessBlock(b.StatelessBlock)
	b.vm.preferredBlkID = b.ID()

	if err := b.vm.state.Commit(); err != nil {
		return err
	}

	b.vm.chainCtx.Log.Info(
		"accepted block",
		zap.Stringer("blkID", b.ID()),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	// TODO: handle dynamic validator set
	// At this stage of hooking stuff up, we consider a static validators set
	// where validators (and especially the mapping validator -> NodeID) is
	// setup in Genesis
	_, err := b.vm.middleware.AcceptBlock(ctx, b.StatelessBlock)
	return err
}

func (b *StatefulBlock) Reject(context.Context) error {
	delete(b.vm.verifiedBlocks, b.ID())
	return nil
}
