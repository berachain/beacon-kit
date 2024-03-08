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

package types

import (
	"github.com/berachain/beacon-kit/config/version"
	"github.com/berachain/beacon-kit/primitives"
)

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs BeaconBlockDeneb,BeaconBlockBodyDeneb,Deposit -include ../../../primitives,../../../engine/types,$GOPATH/pkg/mod/github.com/ethereum/go-ethereum@$GETH_GO_GENERATE_VERSION/common -output generated.ssz.go
type BeaconBlockDeneb struct {
	// Slot represents the position of the block in the chain.
	Slot primitives.Slot
	// ParentBlockRoot is the hash of the parent block.
	ParentBlockRoot [32]byte `ssz-size:"32"`
	// Body is the body of the BeaconBlockDeneb, containing the block's
	// operations.
	Body *BeaconBlockBodyDeneb
	// PayloadValue is a value used in the block's payload.
	PayloadValue [32]byte `ssz-size:"32"`
}

// Version identifies the version of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) Version() int {
	return version.Deneb
}

// IsNil checks if the BeaconBlockDeneb instance is nil.
func (b *BeaconBlockDeneb) IsNil() bool {
	return b == nil
}

// GetBody retrieves the body of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) GetBody() BeaconBlockBody {
	return b.Body
}

// GetSlot retrieves the slot of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) GetSlot() primitives.Slot {
	return b.Slot
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockDeneb.
func (b *BeaconBlockDeneb) GetParentBlockRoot() [32]byte {
	return b.ParentBlockRoot
}
