// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/minio/sha256-simd"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
	"github.com/prysmaticlabs/gohashtree"
	"golang.org/x/sync/errgroup"
)

const (
	// MinParallelizationSize is the minimum size of the input list that
	// should be hashed using the default method. If the input list is smaller
	// than this size, the overhead of parallelizing the hashing process is.
	//
	// TODO: This value is arbitrary and should be benchmarked to find the
	// optimal value.
	MinParallelizationSize = 5000
	// two is a constant to make the linter happy.
	two = 2
)

type HasherFn[RootT ~[32]byte] func([]RootT, []RootT) error

// Hasher can be re-used for constructing Merkle tree roots.
type Hasher[RootT ~[32]byte] struct {
	// buffer is a reusable buffer for hashing.
	buffer bytes.Buffer[RootT]
	// hasher is the hashing function to use.
	hasher HasherFn[RootT]
}

// NewHasher creates a new merkle Hasher.
func NewHasher[RootT ~[32]byte](
	buffer bytes.Buffer[RootT],
	hashFn HasherFn[RootT],
) *Hasher[RootT] {
	return &Hasher[RootT]{
		buffer: buffer,
		hasher: hashFn,
	}
}

// NewRootWithMaxLeaves constructs a Merkle tree root from a set of.
func (m *Hasher[RootT]) NewRootWithMaxLeaves(
	leaves []RootT,
	length math.U64,
) (RootT, error) {
	return m.NewRootWithDepth(leaves, length.NextPowerOfTwo().ILog2Ceil())
}

func (m *Hasher[RootT]) NewRootWithDepth(
	leaves []RootT,
	depth uint8,
) (RootT, error) {
	count := uint64(len(leaves))
	limit := uint64(1) << depth

	if count > limit {
		return zero.Hashes[depth], fmt.Errorf("merkleizing list that is too large, over limit")
	}

	if limit == 0 {
		return zero.Hashes[0], nil
	}

	if limit == 1 {
		if count == 1 {
			return leaves[0], nil
		}
		return zero.Hashes[0], nil
	}

	tmp := m.buffer.Get(int(depth + 1))
	var h RootT

	hh := NewHasherFunc[RootT](sha256.Sum256)

	var j uint8
	merge := func(i uint64) error {
		for j = 0; ; j++ {
			if i&(uint64(1)<<j) == 0 {
				if i == count && j < depth {
					h = hh.Combi(h, zero.Hashes[j])
				} else {
					break
				}
			} else {
				h = hh.Combi(tmp[j], h)
			}
		}
		tmp[j] = h
		return nil
	}

	for i := uint64(0); i < count; i++ {
		h = leaves[i]
		if err := merge(i); err != nil {
			return zero.Hashes[depth], err
		}
	}

	if (uint64(1) << depth) != count {
		h = zero.Hashes[0]
		if err := merge(count); err != nil {
			return zero.Hashes[depth], err
		}
	}

	for j := uint8(depth); j < uint8(depth); j++ {
		tmp[j+1] = hh.Combi(tmp[j], zero.Hashes[j])
	}

	return tmp[depth], nil
}

// BuildParentTreeRoots calls BuildParentTreeRootsWithNRoutines with the
// number of routines set to runtime.GOMAXPROCS(0)-1.
func BuildParentTreeRoots[RootT ~[32]byte](
	outputList, inputList []RootT,
) error {
	err := BuildParentTreeRootsWithNRoutines(
		//#nosec:G103 // on purpose.
		*(*[][32]byte)(unsafe.Pointer(&outputList)),
		//#nosec:G103 // on purpose.
		*(*[][32]byte)(unsafe.Pointer(&inputList)),
		runtime.GOMAXPROCS(0)-1,
	)

	// Convert out back to []RootT using unsafe pointer cas
	return err
}

// BuildParentTreeRootsWithNRoutines optimizes hashing of a list of roots
// using CPU-specific vector instructions and parallel processing. This
// method adapts to the host machine's hardware for potential performance
// gains over sequential hashing.
//
// TODO: We do not use generics here due to the gohashtree library not
// supporting generics.
func BuildParentTreeRootsWithNRoutines(
	outputList, inputList [][32]byte, n int,
) error {
	// Validate input list length.
	inputLength := len(inputList)
	if inputLength%2 != 0 {
		return ErrOddLengthTreeRoots
	}

	// Build output variables
	outputLength := inputLength / two

	// If the input list is small, hash it using the default method since
	// the overhead of parallelizing the hashing process is not worth it.
	if inputLength < MinParallelizationSize {
		//#nosec:G103 // used of unsafe calls should be audited.
		return gohashtree.Hash(outputList, inputList)
	}

	// Otherwise parallelize the hashing process for large inputs.
	// Take the max(n, 1) to prevent division by 0.
	groupSize := inputLength / (two * max(n, 1))
	twiceGroupSize := two * groupSize
	eg := new(errgroup.Group)

	// if n is 0 the parallelization is disabled and the whole inputList is
	// hashed in the main goroutine at the end of this function.
	for j := range n + 1 {
		eg.Go(func() error {
			// inputList:  [-------------------2*groupSize-------------------]
			//              ^                  ^                    ^        ^
			//              |                  |                    |        |
			// j*2*groupSize   (j+1)*2*groupSize    (j+2)*2*groupSize   End
			//
			// outputList: [---------groupSize---------]
			//              ^                         ^
			//              |                         |
			//             j*groupSize         (j+1)*groupSize
			//
			// Each goroutine processes a segment of inputList that is twice as
			// large as the segment it fills in outputList. This is because the
			// hash
			// operation reduces the
			// size of the input by half.
			// Define the segment of the inputList each goroutine will process.
			segmentStart := j * twiceGroupSize
			segmentEnd := min((j+1)*twiceGroupSize, inputLength)

			return gohashtree.Hash(
				outputList[j*groupSize:min((j+1)*groupSize, outputLength)],
				inputList[segmentStart:segmentEnd],
			)
		})
	}

	// Wait for all goroutines to complete.
	return eg.Wait()
}
