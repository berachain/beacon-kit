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

type SSZMarshaler interface {
	MarshalSSZ() ([]byte, error)
}

type SSZUnmarshaler interface {
	UnmarshalSSZ([]byte) error
}

type SSZRootable interface {
	HashTreeRoot() ([32]byte, error)
}

type SSZMarshalerUnmarshaler interface {
	SSZMarshaler
	SSZUnmarshaler
}

type SSZMarshalerUnmarshalerRootable interface {
	SSZMarshaler
	SSZUnmarshaler
	SSZRootable
}

type BaseSSZObject interface {
	// MarshalSSZTo marshals the object into the provided byte slice and returns
	// it along with any error.
	MarshalSSZTo([]byte) ([]byte, error)
	// MarshalSSZ marshals the object into a new byte slice and returns it along
	// with any error.
	MarshalSSZ() ([]byte, error)
	// UnmarshalSSZ unmarshals the object from the provided byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalSSZ([]byte) error
}

// JSONMarshallable is an interface that combines the json.Marshaler and
// json.Unmarshaler interfaces.
type JSONMarshallable interface {
	// MarshalJSON marshals the object into a JSON byte slice and returns it
	// along with any error.
	MarshalJSON() ([]byte, error)
	// UnmarshalJSON unmarshals the object from the provided JSON byte slice and
	// returns an error if the unmarshaling fails.
	UnmarshalJSON([]byte) error
}
