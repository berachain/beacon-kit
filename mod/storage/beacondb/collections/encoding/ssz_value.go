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
	"reflect"

	"cosmossdk.io/collections/codec"
	fssz "github.com/ferranbt/fastssz"
)

// SSZMarshallable defines an interface for types that can be
// marshaled and unmarshaled using SSZ encoding,
// and also provides a string representation of the type.
type SSZMarshallable interface {
	fssz.Marshaler
	fssz.Unmarshaler
	String() string
}

// SSZValueCodec provides methods to encode and decode SSZ values.
type SSZValueCodec[T SSZMarshallable] struct{}

// Assert that SSZValueCodec implements codec.ValueCodec.
var _ codec.ValueCodec[SSZMarshallable] = SSZValueCodec[SSZMarshallable]{}

// Encode marshals the provided value into its SSZ encoding.
func (SSZValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (SSZValueCodec[T]) Decode(b []byte) (T, error) {
	var v T
	//nolint:errcheck // will error in unmarshal if there is a problem.
	v = reflect.New(reflect.TypeOf(v).Elem()).Interface().(T)
	if err := v.UnmarshalSSZ(b); err != nil {
		return v, err
	}
	return v, nil
}

// EncodeJSON is not implemented and will panic if called.
func (SSZValueCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (SSZValueCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (SSZValueCodec[T]) Stringify(value T) string {
	return value.String()
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZValueCodec[T]) ValueType() string {
	return "SSZMarshallable"
}

// SSZInterfaceCodec provides methods to encode and decode SSZ values.
//
// This type exists for codecs for interfaces, which require a factory function
// to create new instances of the underlying hard type since reflect cannot
// infer the type of an interface.
type SSZInterfaceCodec[T SSZMarshallable] struct {
	Factory func() T
}

// Assert that SSZInterfaceCodec implements codec.ValueCodec.
var _ codec.ValueCodec[SSZMarshallable] = SSZInterfaceCodec[SSZMarshallable]{}

// Encode marshals the provided value into its SSZ encoding.
func (SSZInterfaceCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (cdc SSZInterfaceCodec[T]) Decode(b []byte) (T, error) {
	v := cdc.Factory()
	if err := v.UnmarshalSSZ(b); err != nil {
		return v, err
	}

	return v, nil
}

// EncodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (SSZInterfaceCodec[T]) Stringify(value T) string {
	return value.String()
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZInterfaceCodec[T]) ValueType() string {
	return "SSZMarshallable"
}
