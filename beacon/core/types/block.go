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
	"fmt"

	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	"github.com/itsdevbear/bolaris/config/version"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	"github.com/itsdevbear/bolaris/primitives"
)

//go:generate go run github.com/prysmaticlabs/fastssz/sszgen -path . -objs BeaconBlockDeneb,BeaconBlockBodyDeneb -include ../../../primitives,../../../engine/types,./v1,$HOME/go/pkg/mod/github.com/ethereum/go-ethereum@v1.13.14/common -output generated.ssz.go
type BeaconBlockDeneb struct {
	Slot            primitives.Slot
	ParentBlockRoot [32]byte `ssz-size:"32"`
	Body            *BeaconBlockBodyDeneb
	PayloadValue    [32]byte `ssz-size:"32"`
}

func (b *BeaconBlockDeneb) ExecutionPayload() (enginetypes.ExecutionPayload, error) {
	return b.Body.ExecutionPayload, nil
}

func (b *BeaconBlockDeneb) AttachExecution(
	executionData enginetypes.ExecutionPayload,
) error {
	fmt.Println("AttachExecution")
	fmt.Println(executionData.GetBlockHash())
	fmt.Println(executionData)
	b.Body.ExecutionPayload = executionData.(*enginetypes.ExecutableDataDeneb)
	return nil
}

func (b *BeaconBlockDeneb) Version() int {
	return version.Deneb
}

func (b *BeaconBlockDeneb) IsNil() bool {
	return b == nil
}

func (b *BeaconBlockDeneb) GetSlot() primitives.SSZUint64 {
	return primitives.SSZUint64(b.Slot)
}

func (b *BeaconBlockDeneb) GetRandaoReveal() []byte {
	return b.Body.RandaoReveal[:]
}

func (b *BeaconBlockDeneb) GetGraffiti() []byte {
	return b.Body.Graffiti[:]
}

func (b *BeaconBlockDeneb) GetDeposits() []*beacontypesv1.Deposit {
	return b.Body.Deposits
}

func (b *BeaconBlockDeneb) GetParentBlockRoot() []byte {
	return b.ParentBlockRoot[:]
}

func (b *BeaconBlockDeneb) GetBlobKzgCommitments() [][]byte {
	make := make([][]byte, len(b.Body.BlobKzgCommitments))
	for i, v := range b.Body.BlobKzgCommitments {
		make[i] = v[:]
	}
	return make
}

type BeaconBlockBodyDeneb struct {
	RandaoReveal       [96]byte                 `ssz-size:"96"`
	Graffiti           [32]byte                 `ssz-size:"32"`
	Deposits           []*beacontypesv1.Deposit `                ssz-max:"16"`
	ExecutionPayload   *enginetypes.ExecutableDataDeneb
	BlobKzgCommitments [][48]byte `ssz-size:"?,48" ssz-max:"16"`
}
