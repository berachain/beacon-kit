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

// SideCars is a slice of blob side cars to be included in the block.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen -path ./sidecars.go -objs BlobSidecars -include ./sidecar.go,../../primitives/kzg,../../primitives,$GETH_PKG_INCLUDE/common -output sidecars.ssz.go
type BlobSidecars struct {
	// Sidecars is a slice of blob side cars to be included in the block.
	Sidecars []*BlobSidecar `ssz-max:"6"`
}

// ValidateBlockRoots checks to make sure that
// all blobs in the sidecar are from the same block.
func (bs *BlobSidecars) ValidateBlockRoots() error {
	// We only need to check if there is more than
	// a single blob in the sidecar.
	if sc := bs.Sidecars; len(sc) > 1 {
		firstHtr, err := sc[0].BeaconBlockHeader.HashTreeRoot()
		if err != nil {
			return err
		}

		var nextHtr [32]byte
		for i := 1; i < len(sc); i++ {
			nextHtr, err = sc[i].BeaconBlockHeader.HashTreeRoot()
			if err != nil {
				return err
			}
			if firstHtr != nextHtr {
				return ErrSidecarContainsDifferingBlockRoots
			}
		}
	}
	return nil
}
