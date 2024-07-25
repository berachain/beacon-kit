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

package constraints

// SSZMarshaler is an interface for objects that can be
// marshaled to SSZ format.
type SSZMarshaler interface {
	// MarshalSSZ marshals the object into SSZ format.
	MarshalSSZ() ([]byte, error)
}

// SSZUnmarshaler is an interface for objects that can be
// unmarshaled from SSZ format.
type SSZUnmarshaler interface {
	// UnmarshalSSZ unmarshals the object from SSZ format.
	UnmarshalSSZ([]byte) error
}

// SSZRootable is an interface for objects that can compute
// their hash tree root.
type SSZRootable interface {
	// HashTreeRoot computes the hash tree root of the object.
	HashTreeRoot() ([32]byte, error)
}

// SSZMarshalerUnmarshaler is an interface that combines
// SSZMarshaler and SSZUnmarshaler.
type SSZMarshalerUnmarshaler interface {
	SSZMarshaler
	SSZUnmarshaler
}

// SSZMarshalerUnmarshalerRootable is an interface that combines
// SSZMarshaler, SSZUnmarshaler, and SSZRootable.
type SSZMarshalerUnmarshalerRootable interface {
	SSZMarshaler
	SSZUnmarshaler
	SSZRootable
}

// JSONMarshallable is an interface that combines the json.Marshaler and
// json.Unmarshaler interfaces.
type JSONMarshallable interface {
	// MarshalJSON marshals the object into a JSON byte slice and returns it
	// along with any error.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals the object from the provided JSON byte slice
	// and returns an error if the unmarshaling fails.
	UnmarshalJSON([]byte) error
}
