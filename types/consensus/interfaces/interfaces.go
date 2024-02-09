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

	"github.com/itsdevbear/bolaris/math"
	enginev1 "github.com/itsdevbear/bolaris/third_party/prysm/proto/engine/v1"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	ssz "github.com/prysmaticlabs/fastssz"
	"google.golang.org/protobuf/proto"
)

// ExecutionData is the interface for the execution data of a block.
type ExecutionData interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	IsBlinded() bool
	Proto() proto.Message
	ParentHash() []byte
	FeeRecipient() []byte
	StateRoot() []byte
	ReceiptsRoot() []byte
	LogsBloom() []byte
	PrevRandao() []byte
	BlockNumber() uint64
	GasLimit() uint64
	GasUsed() uint64
	Timestamp() uint64
	ExtraData() []byte
	BaseFeePerGas() []byte
	BlobGasUsed() (uint64, error)
	ExcessBlobGas() (uint64, error)
	BlockHash() []byte
	Transactions() ([][]byte, error)
	TransactionsRoot() ([]byte, error)
	Withdrawals() ([]*enginev1.Withdrawal, error)
	WithdrawalsRoot() ([]byte, error)
	PbCapella() (*enginev1.ExecutionPayloadCapella, error)
	PbBellatrix() (*enginev1.ExecutionPayload, error)
	PbDeneb() (*enginev1.ExecutionPayloadDeneb, error)
	ValueInWei() (math.Wei, error)
	ValueInGwei() (uint64, error)
}

// BeaconKitBlock is the interface for a beacon block.
type BeaconKitBlock interface {
	ReadOnlyBeaconKitBlock
	WriteOnlyBeaconKitBlock
}

type BeaconBlockBody interface {
	ReadOnlyBeaconKitBlockBody
}

type ReadOnlyBeaconKitBlockBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
}

// ReadOnlyBeaconKitBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconKitBlock interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	GetSlot() primitives.Slot
	// ProposerAddress() []byte
	IsNil() bool
	// Execution returns the execution data of the block.
	Execution() (ExecutionData, error)
	Version() int
}

// WriteOnlyBeaconKitBlock is the interface for a write-only beacon block.
type WriteOnlyBeaconKitBlock interface {
	AttachExecution(ExecutionData) error
}

// ABCIRequest is the interface for an ABCI request.
type ABCIRequest interface {
	GetHeight() int64
	GetTime() time.Time
	GetTxs() [][]byte
}
