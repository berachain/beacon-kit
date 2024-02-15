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

package sha256

import (
	"runtime"

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

// BuildParentTreeRoots calls BuildParentTreeRootsWithNRoutines with the number of
// routines set to runtime.GOMAXPROCS(0)-1.
func BuildParentTreeRoots(inputList [][32]byte) ([][32]byte, error) {
	return BuildParentTreeRootsWithNRoutines(inputList, runtime.GOMAXPROCS(0)-1)
}

// HashTreeRoot takes a list of roots and hashes them using CPU
// specific vector instructions. Depending on host machine's specific
// hardware configuration, using this routine can lead to a significant
// performance improvement compared to the default method of hashing
// lists.
func BuildParentTreeRootsWithNRoutines(inputList [][32]byte, n int) ([][32]byte, error) {
	if len(inputList)%2 != 0 {
		return nil, ErrOddLengthTreeRoots
	}
	outputList := make([][32]byte, len(inputList)/two)
	// If the input list is small, hash it using the default method since
	// the overhead of parallelizing the hashing process is not worth it.
	if len(inputList) < MinParallelizationSize {
		return outputList, gohashtree.Hash(outputList, inputList)
	}

	// Otherwise parallelize the hashing process for large inputs.

	groupSize := len(inputList) / (two * (n + 1))
	eg := new(errgroup.Group)

	// if n is 0 the parallelization is disabled and the whole inputList is hashed in the main
	// goroutine at the end of this function.
	for j := 0; j < n; j++ {
		// capture loop variable
		cj := j

		// inputList:  [---------------------2*groupSize---------------------]
		//              ^                    ^                    ^          ^
		//              |                    |                    |          |
		//             j*2*groupSize   (j+1)*2*groupSize    (j+2)*2*groupSize  End
		//
		// outputList: [---------groupSize---------]
		//              ^                         ^
		//              |                         |
		//             j*groupSize         (j+1)*groupSize
		//
		// Each goroutine processes a segment of inputList that is twice as large as the
		// segment it fills in outputList. This is because the hash operation reduces the
		// size of the input by half.
		eg.Go(func() error {
			return gohashtree.Hash(
				outputList[cj*groupSize:], inputList[cj*two*groupSize:(cj+1)*two*groupSize],
			)
		})
	}

	// The last segment of inputList is processed here because the division of the inputList
	// among the goroutines might leave a remainder segment that is not exactly divisible by
	// the number of goroutines spawned. This remainder segment is processed in the main goroutine
	// to ensure all parts of the inputList are hashed.
	remainderStartIndex := n * two * groupSize
	if remainderStartIndex < len(inputList) { // Check if there's a remainder segment to process.
		err := gohashtree.Hash(outputList[n*groupSize:], inputList[remainderStartIndex:])
		if err != nil {
			return nil, err
		}
	}

	// Wait for all goroutines to finish processing their segments.
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return outputList, nil
}
