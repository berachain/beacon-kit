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

package serialization

import (
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
)

/* -------------------------------------------------------------------------- */
/*                                    Bool                                    */
/* -------------------------------------------------------------------------- */

// UnmarshalBool unmarshals a bool from SSZ format.
func UnmarshalBool[T ~bool](buf []byte) (T, error) {
	if len(buf) != constants.BoolSize {
		return false, fmt.Errorf("invalid bool length: %d", len(buf))
	}
	return T(buf[0] != 0), nil
}

// MarshalBool marshals a bool to SSZ format.
func MarshalBool[T ~bool](b T) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

/* -------------------------------------------------------------------------- */
/*                                    Byte                                    */
/* -------------------------------------------------------------------------- */

// UnmarshalByte unmarshals a byte from SSZ format.
func UnmarshalByte[T ~byte](buf []byte) (T, error) {
	if len(buf) != constants.ByteSize {
		return 0, fmt.Errorf("invalid byte length: %d", len(buf))
	}
	return T(buf[0]), nil
}

// MarshalByte marshals a byte to SSZ format.
func MarshalByte[T ~byte](b T) []byte {
	return []byte{byte(b)}
}

/* -------------------------------------------------------------------------- */
/*                                    U8                                      */
/* -------------------------------------------------------------------------- */

// UnmarshalU8 unmarshals a uint8 from SSZ format.
func UnmarshalU8[T ~uint8](buf []byte) (T, error) {
	b, err := UnmarshalByte[byte](buf)
	return T(b), err
}

// MarshalU8 marshals a uint8 to SSZ format.
func MarshalU8[T ~uint8](u T) []byte {
	return MarshalByte(byte(u))
}

/* -------------------------------------------------------------------------- */
/*                                    U16                                     */
/* -------------------------------------------------------------------------- */

// UnmarshalU16 unmarshals a uint16 from SSZ format.
func UnmarshalU16[T ~uint16](buf []byte) (T, error) {
	if len(buf) != constants.U16Size {
		return 0, fmt.Errorf("invalid uint16 length: %d", len(buf))
	}
	return T(binary.LittleEndian.Uint16(buf)), nil
}

// MarshalU16 marshals a uint16 to SSZ format.
func MarshalU16[T ~uint16](u T) []byte {
	buf := make([]byte, constants.U16Size)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf
}

/* -------------------------------------------------------------------------- */
/*                                    U32                                     */
/* -------------------------------------------------------------------------- */

// UnmarshalU32 unmarshals a uint32 from SSZ format.
func UnmarshalU32[T ~uint32](buf []byte) (T, error) {
	if len(buf) != constants.U32Size {
		return 0, fmt.Errorf("invalid uint32 length: %d", len(buf))
	}
	return T(binary.LittleEndian.Uint32(buf)), nil
}

// MarshalU32 marshals a uint32 to SSZ format.
func MarshalU32[T ~uint32](u T) []byte {
	buf := make([]byte, constants.U32Size)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf
}

/* -------------------------------------------------------------------------- */
/*                                    U64                                     */
/* -------------------------------------------------------------------------- */

// UnmarshalU64 unmarshals a uint64 from SSZ format.
func UnmarshalU64[T ~uint64](buf []byte) (T, error) {
	if len(buf) != constants.U64Size {
		return 0, fmt.Errorf("invalid uint64 length: %d", len(buf))
	}
	return T(binary.LittleEndian.Uint64(buf)), nil
}

// MarshalU64 marshals a uint64 to SSZ format.
func MarshalU64[T ~uint64](u T) []byte {
	buf := make([]byte, constants.U64Size)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf
}
