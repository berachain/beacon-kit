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

package encoding

import (
	"cosmossdk.io/collections/codec"
)

type SSZMarshallable interface {
	MarshalSSZ() ([]byte, error)
	UnmarshalSSZ([]byte) error
	String() string
}

type SSZValueCodec[T SSZMarshallable] struct{}

// This is an assertion that Deposit implements the codec.ValueCodec interface.
var _ codec.ValueCodec[SSZMarshallable] = SSZValueCodec[SSZMarshallable]{}

func (SSZValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

func (SSZValueCodec[T]) Decode(b []byte) (T, error) {
	var v T
	if err := v.UnmarshalSSZ(b); err != nil {
		return v, err
	}
	return v, nil
}

func (SSZValueCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

func (SSZValueCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

func (SSZValueCodec[T]) Stringify(value T) string {
	return value.String()
}

func (SSZValueCodec[T]) ValueType() string {
	return "SSZMarshallable"
}
