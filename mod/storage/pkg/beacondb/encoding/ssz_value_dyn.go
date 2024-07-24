// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package encoding

import (
	"reflect"

	"cosmossdk.io/collections/codec"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/davecgh/go-spew/spew"
)

// SSZValueCodecDyn provides methods to encode and decode SSZ values.
type SSZValueCodecDyn[T constraints.SSZMarshallableDynamic] struct{}

// Assert that SSZValueCodecDyn implements codec.ValueCodec.
//
//nolint:lll // annoying formatter.
var _ codec.ValueCodec[constraints.SSZMarshallableDynamic] = SSZValueCodecDyn[constraints.SSZMarshallableDynamic]{}

// Encode marshals the provided value into its SSZ encoding.
func (SSZValueCodecDyn[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (SSZValueCodecDyn[T]) Decode(b []byte) (T, error) {
	var v T
	//nolint:errcheck // will error in unmarshal if there is a problem.
	v = reflect.New(reflect.TypeOf(v).Elem()).Interface().(T)
	if err := v.UnmarshalSSZ(b); err != nil {
		return v, err
	}
	return v, nil
}

// EncodeJSON is not implemented and will panic if called.
func (SSZValueCodecDyn[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (SSZValueCodecDyn[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (SSZValueCodecDyn[T]) Stringify(value T) string {
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZValueCodecDyn[T]) ValueType() string {
	return "SSZMarshallableDynamic"
}

// SSZInterfaceCodecDyn provides methods to encode and decode SSZ values.
//
// This type exists for codecs for interfaces, which require a factory function
// to create new instances of the underlying hard type since reflect cannot
// infer the type of an interface.
type SSZInterfaceCodecDyn[T interface {
	constraints.SSZMarshallableDynamic
	NewFromSSZ([]byte, uint32) (T, error)
	Version() uint32
}] struct {
	latestVersion uint32
}

// SetActiveForkVersion sets the fork version for the codec.
func (cdc *SSZInterfaceCodecDyn[T]) SetActiveForkVersion(version uint32) {
	cdc.latestVersion = version
}

// Encode marshals the provided value into its SSZ encoding.
func (cdc *SSZInterfaceCodecDyn[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (cdc SSZInterfaceCodecDyn[T]) Decode(b []byte) (T, error) {
	var t T
	return t.NewFromSSZ(b, cdc.latestVersion)
}

// EncodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodecDyn[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodecDyn[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (SSZInterfaceCodecDyn[T]) Stringify(value T) string {
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZInterfaceCodecDyn[T]) ValueType() string {
	return "SSZMarshallableDynamic"
}
