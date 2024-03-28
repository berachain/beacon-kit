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
	"encoding/binary"

	"github.com/berachain/beacon-kit/primitives"
)

type (
	// SSZUint64 is a special uint64 that implements
	// signature.SSZObject interface.
	SSZUInt64 uint64
)

// HashTreeRoot return the merklized epoch,
// represented as bytes in little endian,
// padded on the right side with zeroed bytes
// to a total of 32 bytes.
func (e SSZUInt64) HashTreeRoot() (primitives.Root, error) {
	bz := make([]byte, primitives.RootLength)
	binary.LittleEndian.PutUint64(bz, uint64(e))
	return primitives.Root(bz), nil
}
