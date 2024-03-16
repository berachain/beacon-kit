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

package htr

import (
	"runtime"
	"sync"

	"github.com/prysmaticlabs/gohashtree"
)

const minSliceSizeToParallelize = 5000

func hashParallel(
	inputList [][32]byte,
	outputList [][32]byte,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	err := gohashtree.Hash(outputList, inputList)
	if err != nil {
		panic(err)
	}
}

// VectorizedSha256 takes a list of roots and hashes them using CPU
// specific vector instructions. Depending on host machine's specific
// hardware configuration, using this routine can lead to a significant
// performance improvement compared to the default method of hashing
// lists.
func VectorizedSha256(inputList [][32]byte) [][32]byte {
	outputList := make([][32]byte, len(inputList)/2)
	if len(inputList) < minSliceSizeToParallelize {
		err := gohashtree.Hash(outputList, inputList)
		if err != nil {
			panic(err)
		}
		return outputList
	}
	n := runtime.GOMAXPROCS(0) - 1
	wg := sync.WaitGroup{}
	wg.Add(n)
	groupSize := len(inputList) / (2 * (n + 1))
	for j := 0; j < n; j++ {
		go hashParallel(
			inputList[j*2*groupSize:(j+1)*2*groupSize],
			outputList[j*groupSize:],
			&wg,
		)
	}
	err := gohashtree.Hash(outputList[n*groupSize:], inputList[n*2*groupSize:])
	if err != nil {
		panic(err)
	}
	wg.Wait()
	return outputList
}
