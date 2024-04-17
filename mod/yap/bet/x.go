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

package bet

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/ssz"
)

// maxDepProvider is an interface for types that provide the maximum number of
// deposits per block.
type maxDepProvider interface {
	MaxDepositsPerBlock() uint64
}

// Deposits is a list of deposits.
type Deposits[T maxDepProvider] []*primitives.Deposit

// HashTreeRoot returns the root of the hash tree of the deposits.
func (d Deposits[T]) HashTreeRoot(chainSpec ...T) primitives.Root {
	if chainSpec == nil {
		return primitives.Root{}
	}
	root, err := ssz.MerkleizeList(d, chainSpec[0].MaxDepositsPerBlock())
	fmt.Println("DEPOSITS ROOT", err, root)
	return root
}

type BasicType[T any] uint64

// Encode as little endian and then compute SHA256 hash.
func (b BasicType[T]) HashTreeRoot(...T) primitives.Root {
	var bytes [32]byte
	binary.LittleEndian.PutUint64(bytes[:], uint64(b))
	hash := sha256.Sum256(bytes[:])
	return primitives.Root(hash)
}

type YapCave[T primitives.ChainSpec] struct {
	Item2 Deposits[T]
}

//go:generate go run github.com/ferranbt/fastssz/sszgen --path . --objs YapCave2 -include ../../primitives
type YapCave2 struct {
	Item2 []*primitives.Deposit `ssz-max:"16"`
}

func (y *YapCave[T]) HashTreeRoot(chainSpec ...T) primitives.Root {
	vec := []ssz.Hashable2[T]{y.Item2}
	return ssz.MerkleizeVector2(vec, uint64(len(vec)), chainSpec...)
}
