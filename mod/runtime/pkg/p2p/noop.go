package p2p

import (
	"context"

	ssz "github.com/ferranbt/fastssz"
)

// NoopGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopGossipHandler[DataT interface {
	ssz.Marshaler
	ssz.Unmarshaler
}, BytesT ~[]byte] struct{}

// Publish creates a new NoopGossipHandler.
func (n NoopGossipHandler[DataT, BytesT]) Publish(ctx context.Context, data DataT) (BytesT, error) {
	return data.MarshalSSZ()
}

// Request simply returns the reference it receives.
func (n NoopGossipHandler[DataT, BytesT]) Request(ctx context.Context, ref BytesT, out DataT) error {
	return out.UnmarshalSSZ(ref)
}
