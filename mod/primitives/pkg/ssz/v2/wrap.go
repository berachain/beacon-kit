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

import (
	"reflect"

	fastssz "github.com/ferranbt/fastssz"
)

// wrap implements a wrapper class which can wrap ssz-able go structs and add
// fastssz iface compliant functions.
type GenericSSZIFace interface {
	Type
	fastssz.Marshaler
	// Todo
	// fastssz.Unmarshaler
	// fastssz.HashRoot
}

type Wrapper struct {
	GenericSSZIFace
	data       any
	serializer Serializer
}

func Wrap[T any](t T) Wrapper {
	return Wrapper{
		data:       t,
		serializer: NewSerializer(),
	}
}

func (s *Wrapper) SizeSSZ() int {
	res, err := s.serializer.GetSize(
		reflect.ValueOf(s.data),
		reflect.TypeOf(s.data),
	)
	if err != nil {
		return 0
	}
	return res
}

func (s *Wrapper) MarshalSSZ() ([]byte, error) {
	res, err := s.serializer.MarshalSSZ(reflect.ValueOf(s.data).Interface())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Wrapper) MarshalSSZTo(buf []byte) ([]byte, error) {
	_, err := s.serializer.Marshal(
		reflect.ValueOf(s.data),
		reflect.TypeOf(s.data),
		buf,
		uint64(len(buf)),
	)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *Wrapper) HashTreeRoot() ([32]byte, error) {
	return fastssz.HashWithDefaultHasher(s)
}

// GetTree ssz hashes the AttestationData object.
func (s *Wrapper) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(s)
}

func (s *Wrapper) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// todo - add each field to hh

	hh.Merkleize(indx)
	return nil
}
