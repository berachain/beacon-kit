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

package abci

import (
	"context"

	"github.com/berachain/beacon-kit/mod/core/state"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type BuilderService interface {
	RequestBestBlock(
		context.Context,
		state.BeaconState,
		math.Slot,
	) (consensus.BeaconBlock, *datypes.BlobSidecars, error)
}

type BlockchainService interface {
	ProcessSlot(state.BeaconState) error
	BeaconState(context.Context) state.BeaconState
	ProcessBeaconBlock(
		context.Context,
		state.BeaconState,
		consensus.ReadOnlyBeaconBlock,
		*datypes.BlobSidecars,
	) error
	PostBlockProcess(
		context.Context,
		state.BeaconState,
		consensus.ReadOnlyBeaconBlock,
	) error
	ChainSpec() primitives.ChainSpec
	VerifyPayloadOnBlk(context.Context, consensus.ReadOnlyBeaconBlock) error
}
