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

import (
	"fmt"

	fssz "github.com/prysmaticlabs/fastssz"
)

const (
	epochSize = 8
)

var _ fssz.HashRoot = (Epoch)(0)
var _ fssz.Marshaler = (*Epoch)(nil)
var _ fssz.Unmarshaler = (*Epoch)(nil)

// Epoch represents a single epoch.
// We don't really use epochs in BeaconKit.
// But we keep them around for compatibility with the Ethereum 2.0 spec.
type Epoch uint64

// HashTreeRoot --.
func (e Epoch) HashTreeRoot() ([32]byte, error) {
	return fssz.HashWithDefaultHasher(e)
}

// HashTreeRootWith --.
func (e Epoch) HashTreeRootWith(hh *fssz.Hasher) error {
	hh.PutUint64(uint64(e))
	return nil
}

// UnmarshalSSZ --.
func (e *Epoch) UnmarshalSSZ(buf []byte) error {
	if len(buf) != e.SizeSSZ() {
		return fmt.Errorf(
			"expected buffer of length %d received %d",
			e.SizeSSZ(),
			len(buf),
		)
	}
	*e = Epoch(fssz.UnmarshallUint64(buf))
	return nil
}

// MarshalSSZTo --.
func (e *Epoch) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalled, err := e.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, marshalled...), nil
}

// MarshalSSZ --.
func (e *Epoch) MarshalSSZ() ([]byte, error) {
	marshalled := fssz.MarshalUint64([]byte{}, uint64(*e))
	return marshalled, nil
}

// SizeSSZ --.
func (e *Epoch) SizeSSZ() int {
	return epochSize
}
