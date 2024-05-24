package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/mod/beacon/dispatcher"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// TODO: unhood this name
const EventType = "new_block"

type Event struct {
	Ctx  context.Context
	Slot math.U64
}

func NewEvent(
	ctx context.Context, slot math.U64,
) dispatcher.Event {
	return &Event{
		Ctx:  ctx,
		Slot: slot,
	}
}

func (e *Event) Type() dispatcher.EventType {
	return EventType
}
