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

package ssz

import (
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/protolambda/ztyp/tree"

	// We need to import this package to use the VectorizedSha256 function.
	_ "github.com/minio/sha256-simd"
)

// Bytes is a byte slice representing a transaction.
type Bytes []byte

// HashTreeRoot returns the hash tree root of the transaction.
func (bz Bytes) HashTreeRoot() ([32]byte, error) {
	return tree.GetHashFn().ByteListHTR(bz, primitives.MaxBytesPerTxLength), nil
}

// TransactionsRoot computes the HTR for the Transactions' property of the ExecutionPayload
// The code was largely copy/pasted from the code generated to compute the HTR of the entire
// ExecutionPayload.
// TODO: create strong types and make put these functions on their receivers.
func TransactionsRoot(txs []Bytes) ([32]byte, error) {
	return MerkleizeVectorSSZAndMixinLength(txs, primitives.MaxTxsPerPayloadLength)
}
