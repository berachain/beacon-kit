// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/davecgh/go-spew/spew"
)

// SSZValueCodec provides methods to encode and decode SSZ values.
type SSZValueCodec[T constraints.SSZMarshallable] struct {
	NewEmptyF func() T // constructor
}

// Encode marshals the provided value into its SSZ encoding.
func (SSZValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (sc SSZValueCodec[T]) Decode(bz []byte) (T, error) {
	dest := sc.NewEmptyF()
	return dest, constraints.SSZUnmarshal(bz, dest)
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
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (SSZValueCodec[T]) ValueType() string {
	return "SSZMarshallable"
}

// SSZVersionedValueCodec provides methods to encode and decode SSZ values for a specific version.
type SSZVersionedValueCodec[T constraints.SSZMarshallable] struct {
	NewEmptyF     func(common.Version) T // constructor
	latestVersion common.Version
}

// SetActiveForkVersion sets the fork version for the codec.
func (cdc *SSZVersionedValueCodec[T]) SetActiveForkVersion(version common.Version) {
	cdc.latestVersion = version
}

// Encode marshals the provided value into its SSZ encoding.
func (cdc *SSZVersionedValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (cdc *SSZVersionedValueCodec[T]) Decode(b []byte) (T, error) {
	dest := cdc.NewEmptyF(cdc.latestVersion)
	return dest, constraints.SSZUnmarshal(b, dest)
}

// EncodeJSON is not implemented and will panic if called.
func (cdc *SSZVersionedValueCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (cdc *SSZVersionedValueCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (cdc *SSZVersionedValueCodec[T]) Stringify(value T) string {
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (cdc *SSZVersionedValueCodec[T]) ValueType() string {
	return "SSZVersionedMarshallable"
}
