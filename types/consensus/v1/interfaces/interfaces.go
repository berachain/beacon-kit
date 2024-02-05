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

package interfaces

import (
	"time"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
)

// ExecutionData is the interface for the execution data of a block.
type ExecutionData = interfaces.ExecutionData

// BeaconKitBlock is the interface for a beacon block.
type BeaconKitBlock interface {
	ReadOnlyBeaconKitBlock
	WriteOnlyBeaconKitBlock
}

// ReadOnlyBeaconKitBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconKitBlock interface {
	GetSlot() primitives.Slot
	// ProposerAddress() []byte
	IsNil() bool
	// Execution returns the execution data of the block.
	Execution() (interfaces.ExecutionData, error)

	// Marshal is the interface for marshalling a beacon block.
	Marshal() ([]byte, error)
	// Unmarshal is the interface for unmarshalling a beacon block.
	Unmarshal([]byte) error
}

// WriteOnlyBeaconKitBlock is the interface for a write-only beacon block.
type WriteOnlyBeaconKitBlock interface {
	AttachExecution(interfaces.ExecutionData) error
}

// ABCIRequest is the interface for an ABCI request.
type ABCIRequest interface {
	GetHeight() int64
	GetTime() time.Time
	GetTxs() [][]byte
}
