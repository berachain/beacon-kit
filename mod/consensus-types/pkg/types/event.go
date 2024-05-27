package types

import "context"

type BlockEvent struct {
	ctx   context.Context
	block BeaconBlock
}

func NewBlockEvent(ctx context.Context, block BeaconBlock) BlockEvent {
	return BlockEvent{
		ctx:   ctx,
		block: block,
	}
}

func (e BlockEvent) Context() context.Context {
	return e.ctx
}

func (e BlockEvent) Block() BeaconBlock {
	return e.block
}
