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

package p2p

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
)

// NoopGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopBlockGossipHandler[ReqT encoding.ABCIRequest] struct {
	NoopGossipHandler[consensus.BeaconBlock, []byte]
	chainSpec common.ChainSpec
}

func NewNoopBlockGossipHandler[ReqT encoding.ABCIRequest](
	chainSpec common.ChainSpec,
) NoopBlockGossipHandler[ReqT] {
	return NoopBlockGossipHandler[ReqT]{
		NoopGossipHandler: NoopGossipHandler[consensus.BeaconBlock, []byte]{},
		chainSpec:         chainSpec,
	}
}

// Publish takes a BeaconBlock and returns the ssz marshalled data.
func (n NoopBlockGossipHandler[ReqT]) Publish(
	_ context.Context,
	data consensus.BeaconBlock,
) ([]byte, error) {
	return data.MarshalSSZ()
}

// Request takes an ABCI Request and returns a BeaconBlock.
func (n NoopBlockGossipHandler[ReqT]) Request(
	_ context.Context,
	req ReqT,
) (consensus.BeaconBlock, error) {
	return encoding.UnmarshalBeaconBlockFromABCIRequest(
		req,
		0,
		n.chainSpec.ActiveForkVersionForSlot(math.U64(req.GetHeight())),
	)
}
