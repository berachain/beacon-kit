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

package primitives

import "github.com/berachain/beacon-kit/mod/primitives/math"

// BeaconBlockHeader is the header of a beacon block.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path header.go -objs BeaconBlockHeader -include ./primitives.go,./math,./bytes.go,./math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output header.ssz.go
type BeaconBlockHeader struct {
	// Slot is the slot number of the block.
	Slot math.Slot `json:"slot"`
	// ProposerIndex is the index of the proposer of the block.
	ProposerIndex math.ValidatorIndex `json:"proposerIndex"`
	// ParentRoot is the root of the parent block.
	ParentRoot Root `json:"parentRoot"    ssz-size:"32"`
	// StateRoot is the root of the beacon state after executing
	// the block. Will be 0x00...00 prior to execution.
	StateRoot Root `json:"stateRoot"     ssz-size:"32"`
	// 	// BodyRoot is the root of the block body.
	BodyRoot Root `json:"bodyRoot"      ssz-size:"32"`
}

// // HashTreeRoot ssz hashes the BeaconBlockHeader object
// func (b *BeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
// 	x, err := ssz.MerkleizeContainer[U64](b)
// 	if err != nil {
// 		return [32]byte{}, err
// 	}
// 	y, _ := fssz.HashWithDefaultHasher(b)
// 	if x != y {
// 		fmt.Println("HashTreeRoot mismatch", Root(x), Root(y))
// 		panic("HashTreeRoot mismatch")
// 	}
// 	return x, nil
// }
