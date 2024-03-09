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
	"errors"

	"github.com/berachain/beacon-kit/config/version"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs BeaconBlockDeneb,BeaconBlockBodyDeneb,Deposit -include ../../../primitives,../../../engine/types,$GOPATH/pkg/mod/github.com/ethereum/go-ethereum@$GETH_GO_GENERATE_VERSION/common -output generated.ssz.go
type BeaconBlockDeneb struct {
	Slot            primitives.Slot
	ParentBlockRoot [32]byte `ssz-size:"32"`
	Body            *BeaconBlockBodyDeneb
	PayloadValue    [32]byte `ssz-size:"32"`
}

// IsEmpty returns true if the block is nil or the body is nil.
func (b *BeaconBlockDeneb) IsNil() bool {
	return b == nil
}

func (b *BeaconBlockDeneb) GetBody() BeaconBlockBody {
	return b.Body
}

// Version returns the version of the block.
func (b *BeaconBlockDeneb) Version() int {
	return version.Deneb
}

func (b *BeaconBlockDeneb) GetSlot() primitives.Slot {
	return b.Slot
}

func (b *BeaconBlockDeneb) GetParentBlockRoot() [32]byte {
	return b.ParentBlockRoot
}

type BeaconBlockBodyDeneb struct {
	RandaoReveal       [96]byte   `ssz-size:"96"`
	Graffiti           [32]byte   `ssz-size:"32"`
	Deposits           []*Deposit `                ssz-max:"16"`
	ExecutionPayload   *enginetypes.ExecutableDataDeneb
	BlobKzgCommitments [][48]byte `ssz-size:"?,48" ssz-max:"16"`
}

func (b *BeaconBlockBodyDeneb) GetRandaoReveal() []byte {
	return b.RandaoReveal[:]
}

//
//nolint:lll
func (b *BeaconBlockBodyDeneb) GetExecutionPayload() enginetypes.ExecutionPayload {
	return b.ExecutionPayload
}

func (b *BeaconBlockBodyDeneb) AttachExecution(
	executionData enginetypes.ExecutionPayload,
) error {
	var ok bool
	b.ExecutionPayload, ok = executionData.(*enginetypes.ExecutableDataDeneb)
	if !ok {
		return errors.New("invalid execution data type")
	}
	return nil
}

func (b *BeaconBlockBodyDeneb) GetBlobKzgCommitments() [][48]byte {
	return b.BlobKzgCommitments
}
