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

package ssz

import "github.com/berachain/beacon-kit/lib/ssz/common"

const (
	two = 2
)

func MerkleizeBasic(value common.SSZObject) ([32]byte, error) {
	trie, err := BuildTreeFromBasic(value)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}

// MerkleizeComposite hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func MerkleizeComposite(value common.Composite) ([32]byte, error) {
	trie, err := BuildTreeFromComposite(value)
	if err != nil {
		return [32]byte{}, err
	}
	return trie.HashTreeRoot(), nil
}
