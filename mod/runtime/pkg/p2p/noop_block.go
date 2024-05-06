package p2p

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ABCIRequest is the interface for an ABCI request.
type ABCIRequest interface {
	GetHeight() int64
	GetTime() time.Time
	GetTxs() [][]byte
}

// NoopGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopBlockGossipHandler[ReqT ABCIRequest] struct {
	NoopGossipHandler[consensus.BeaconBlock, []byte]
	chainSpec common.ChainSpec
}

func NewNoopBlockGossipHandler[ReqT ABCIRequest](chainSpec common.ChainSpec) NoopBlockGossipHandler[ReqT] {
	return NoopBlockGossipHandler[ReqT]{
		NoopGossipHandler: NoopGossipHandler[consensus.BeaconBlock, []byte]{},
		chainSpec:         chainSpec,
	}
}

func (n NoopBlockGossipHandler[ReqT]) Publish(ctx context.Context, data consensus.BeaconBlock) ([]byte, error) {
	return data.MarshalSSZ()
}

// Request takes an ABCI Request and returns a BeaconBlock.
func (n NoopBlockGossipHandler[ReqT]) Request(ctx context.Context, req ReqT, out consensus.BeaconBlock) error {
	txs := req.GetTxs()

	// Ensure there are transactions in the request and
	// that the request is valid.
	if lenTxs := uint(len(txs)); txs == nil || lenTxs == 0 {
		return ErrNoBeaconBlockInRequest
	}

	// Extract the beacon block from the ABCI request.
	blkBz := txs[0]
	if blkBz == nil {
		return ErrNilBeaconBlockInRequest
	}

	slot := math.Slot(req.GetHeight())

	var err error
	out, err = consensus.BeaconBlockFromSSZ(blkBz, n.chainSpec.ActiveForkVersionForSlot(slot))
	return err
}
