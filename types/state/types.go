// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package state

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

// BeaconStateProvider provides access to the current beacon state.
type BeaconStateProvider interface {
	// BeaconState returns the current beacon state based on the supplied context.
	BeaconState(ctx context.Context) BeaconState
}

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state interfaces.
type BeaconState interface {
	ReadOnlyBeaconState
	WriteOnlyBeaconState
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	ReadOnlyForkChoice
	ReadOnlyGenesis
	Slot() primitives.Slot
	Time() uint64
	Version() int
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyForkChoice
	WriteOnlyGenesis
}

// Write Only Fork Choice.
type WriteOnlyForkChoice interface {
	// Write Only Fork Choice.
	SetSafeEth1BlockHash(safeBlockHash common.Hash)
	SetFinalizedEth1BlockHash(finalizedBlockHash common.Hash)
	SetLastValidHead(lastValidHead common.Hash)
}

type ReadOnlyForkChoice interface {
	GetLastValidHead() common.Hash
	GetSafeEth1BlockHash() common.Hash
	GetFinalizedEth1BlockHash() common.Hash
}

type ReadOnlyGenesis interface {
	GenesisEth1Hash() common.Hash
}

type WriteOnlyGenesis interface {
	SetGenesisEth1Hash(genesisEth1Hash common.Hash)
}
