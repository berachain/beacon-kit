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

package p2p

import (
	"context"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
)

// NoopBlockGossipHandler is a gossip handler that simply returns the
// ssz marshalled data as a "reference" to the object it receives.
type NoopBlockGossipHandler[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	ReqT encoding.ABCIRequest,
] struct {
	NoopGossipHandlerDyn[BeaconBlockT, []byte]
	chainSpec common.ChainSpec
}

func NewNoopBlockGossipHandler[
	BeaconBlockT BeaconBlock[BeaconBlockT],
	ReqT encoding.ABCIRequest,
](
	chainSpec common.ChainSpec,
) NoopBlockGossipHandler[BeaconBlockT, ReqT] {
	return NoopBlockGossipHandler[BeaconBlockT, ReqT]{
		NoopGossipHandlerDyn: NoopGossipHandlerDyn[BeaconBlockT, []byte]{},
		chainSpec:            chainSpec,
	}
}

// Publish takes a BeaconBlock and returns the ssz marshalled data.
func (n NoopBlockGossipHandler[BeaconBlockT, ReqT]) Publish(
	_ context.Context,
	data BeaconBlockT,
) ([]byte, error) {
	fmt.Println("NOOP BLOCK PUBLISH")
	x, err := data.MarshalSSZ()
	fmt.Println("NOOP BLOCK PUBLISH", len(x), err)
	return x, err
}

// Request takes an ABCI Request and returns a BeaconBlock.
func (n NoopBlockGossipHandler[BeaconBlockT, ReqT]) Request(
	_ context.Context,
	req ReqT,
) (BeaconBlockT, error) {
	b, err := encoding.UnmarshalBeaconBlockFromABCIRequest[BeaconBlockT](
		req,
		0,
		n.chainSpec.ActiveForkVersionForSlot(math.U64(req.GetHeight())),
	)
	return b, err
}
