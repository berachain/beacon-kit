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
	"fmt"

	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/davecgh/go-spew/spew"
)

// FastSSZValueCodec provides methods to encode and decode SSZ values.
type FastSSZValueCodec[T constraints.FastSSZMarshallable] struct {
	NewEmptyF func() T // constructor
}

// Encode marshals the provided value into its SSZ encoding.
func (FastSSZValueCodec[T]) Encode(value T) ([]byte, error) {
	return value.MarshalSSZ()
}

// Decode unmarshals the provided bytes into a value of type T.
func (sc FastSSZValueCodec[T]) Decode(bz []byte) (T, error) {
	dest := sc.NewEmptyF()

	if err := dest.UnmarshalSSZ(bz); err != nil {
		return dest, fmt.Errorf("failed decoding %T: %w", dest, err)
	}

	// Note: ValidateAfterDecodingSSZ may change v even if it returns error
	// (depending on the specific implementations)
	if err := dest.ValidateAfterDecodingSSZ(); err != nil {
		return dest, fmt.Errorf("failed validating %T: %w", dest, err)
	}

	return dest, nil
}

// EncodeJSON is not implemented and will panic if called.
func (FastSSZValueCodec[T]) EncodeJSON(_ T) ([]byte, error) {
	panic("not implemented")
}

// DecodeJSON is not implemented and will panic if called.
func (FastSSZValueCodec[T]) DecodeJSON(_ []byte) (T, error) {
	panic("not implemented")
}

// Stringify returns the string representation of the provided value.
func (FastSSZValueCodec[T]) Stringify(value T) string {
	return spew.Sdump(value)
}

// ValueType returns the name of the interface that this codec is intended for.
func (FastSSZValueCodec[T]) ValueType() string {
	return "SSZMarshallable"
}
