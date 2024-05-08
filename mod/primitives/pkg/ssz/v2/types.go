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

import "reflect"

type (
	// Kind represents different types of SSZ'able data.
	Kind uint8

	// Type represents a SSZ type.
	Type interface {
		Kind() Kind
	}
)
type GenericSSZIFace interface {
	Type
	SizeSSZ() int
	HashTreeRoot() ([32]byte, error)
	UnmarshalSSZ(buf []byte) error
	MarshalSSZ() ([]byte, error)
}

type SSZWrapper struct {
	GenericSSZIFace
	data       any
	serializer Serializer
}

func Wrap[T any](t T) SSZWrapper {
	return SSZWrapper{
		data:       t,
		serializer: NewSerializer(),
	}
}

func (s *SSZWrapper) SizeSSZ() int {
	res, err := s.serializer.GetSize(reflect.ValueOf(s.data), reflect.TypeOf(s.data))
	if err != nil {
		return 0
	}
	return res
}

func (s *SSZWrapper) MarshalSSZ() ([]byte, error) {
	res, err := s.serializer.MarshalSSZ(reflect.ValueOf(s.data).Interface())
	if err != nil {
		return nil, err
	}
	return res, nil
}
