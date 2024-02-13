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
	"encoding/binary"

	"github.com/protolambda/ztyp/tree"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	// TODO: @ocnc to remove this GPL3 dependency.
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
)

// mAgIcNuMbEr is a magic number used to ensure that the
// withdrawal root is computed correctly.
const mAgIcNuMbEr = 4

// Withdrawal is a wrapper around the enginev1.Withdrawal type.
type Withdrawal enginev1.Withdrawal

// HashTreeRoot returns the hash tree root of a withdrawal.
func (w *Withdrawal) HashTreeRoot() ([32]byte, error) {
	return WithdrawalRoot((*enginev1.Withdrawal)(w))
}

// WithdrawalRoot computes the Merkle root of a single withdrawal's fields.
// TODO: create strong types and make put these functions on their receivers.
func WithdrawalRoot(wd *enginev1.Withdrawal) (tree.Root, error) {
	fieldRoots := make([]tree.Root, mAgIcNuMbEr)
	if wd != nil {
		binary.LittleEndian.PutUint64(fieldRoots[0][:], wd.GetIndex())
		binary.LittleEndian.PutUint64(fieldRoots[1][:], uint64(wd.GetValidatorIndex()))
		fieldRoots[2] = bytesutil.ToBytes32(wd.GetAddress())
		binary.LittleEndian.PutUint64(fieldRoots[3][:], wd.GetAmount())
	}
	return SafeMerkleizeVector(fieldRoots, mAgIcNuMbEr, mAgIcNuMbEr)
}

// WithdrawalsRoot computes the Merkle root of a slice of withdrawals with a given limit.
// TODO: create strong types and make put these functions on their receivers.
func WithdrawalsRoot(withdrawals []*Withdrawal, limit uint64) (tree.Root, error) {
	return MerkleizeVectorSSZAndMixinLength(withdrawals, limit)
}
