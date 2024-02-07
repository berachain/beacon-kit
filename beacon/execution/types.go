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

package execution

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/state"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

// BeaconStateProvider is an interface that wraps the basic BeaconState method.
type BeaconStateProvider interface {
	// BeaconState provides access to the underlying beacon state.
	BeaconState(ctx context.Context) state.BeaconState
}

// FCUConfig is a struct that holds the configuration for a fork choice update.
type FCUConfig struct {
	// HeadEth1Hash is the hash of the head eth1 block we are updating the
	// execution client's head to be.
	HeadEth1Hash common.Hash

	// ProposingSlot is the slot that the execution client should propose a block
	// for if Attributes neither nil nor empty.
	ProposingSlot primitives.Slot

	// Attributes is a list of payload attributes to include in a forkchoice update
	// to the execution client. It is used to signal to the execution client that
	// it should build a payload.
	Attributes payloadattribute.Attributer
}
