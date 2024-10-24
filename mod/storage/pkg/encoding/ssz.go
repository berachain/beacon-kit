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
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/davecgh/go-spew/spew"
)

var errNotImplemented = errors.New("not implemented")

// SSZValueCodec provides methods to encode and decode SSZ values.
type SSZValueCodec[T interface {
	constraints.SSZMarshallable
	constraints.Empty[T]
}] struct{}

// Encode marshals the provided value into its SSZ encoding.
func (SSZValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (SSZValueCodec[T]) Decode(bz []byte) (T, error) {
	var v T
	v = (v).Empty()
	return v, v.UnmarshalSSZ(bz)
}

// EncodeJSON is not implemented and will panic if called.
func (SSZValueCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic(fmt.Errorf("EncodeJson %w", errNotImplemented))
}

// DecodeJSON is not implemented and will panic if called.
func (SSZValueCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic(fmt.Errorf("DecodeJSON %w", errNotImplemented))
}

// Stringify returns the string representation of the provided value.
func (SSZValueCodec[T]) Stringify(value T) string {
	return spew.Sdump(value)
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
type SSZInterfaceCodec[T interface {
	constraints.SSZMarshallable
	constraints.Versionable
	NewFromSSZ([]byte, uint32) (T, error)
}] struct {
	latestVersion uint32
}

// SetActiveForkVersion sets the fork version for the codec.
func (cdc *SSZInterfaceCodec[T]) SetActiveForkVersion(version uint32) {
	cdc.latestVersion = version
}

// Encode marshals the provided value into its SSZ encoding.
func (cdc *SSZInterfaceCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (cdc SSZInterfaceCodec[T]) Decode(b []byte) (T, error) {
	var t T
	return t.NewFromSSZ(b, cdc.latestVersion)
}

// EncodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic(fmt.Errorf("EncodeJSON %w", errNotImplemented))
}

// DecodeJSON is not implemented and will panic if called.
func (SSZInterfaceCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic(fmt.Errorf("DecodeJSON %w", errNotImplemented))
}

// Stringify returns the string representation of the provided value.
func (SSZInterfaceCodec[T]) Stringify(value T) string {
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZInterfaceCodec[T]) ValueType() string {
	return "SSZMarshallable"
}
