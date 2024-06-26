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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"encoding/binary"
	"fmt"
)

type Basic[BasicT any] interface {
	// As per the spec, a basic type is one of the following:
	~bool | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	// TODO add uint256, uint128

	// We also need the support the following functions.
	NewFromSSZ([]byte) (BasicT, error)
	MarshalSSZ() ([]byte, error)
	SizeSSZ() int
}

type SSZBool bool

// SizeSSZ returns the size of the bool in bytes.
func (b SSZBool) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the bool into SSZ format.
func (b SSZBool) MarshalSSZ() ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// NewFromSSZ creates a new SSZBool from SSZ format.
func (SSZBool) NewFromSSZ(buf []byte) (SSZBool, error) {
	if len(buf) != 1 {
		return false, fmt.Errorf("invalid buffer length: expected 1, got %d", len(buf))
	}
	return SSZBool(buf[0] != 0), nil
}

// -----------------------------

// -----------------------------

type SSZUInt8 uint8

// SizeSSZ returns the size of the uint8 in bytes.
func (u SSZUInt8) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the uint8 into SSZ format.
func (u SSZUInt8) MarshalSSZ() ([]byte, error) {
	return []byte{byte(u)}, nil
}

// NewFromSSZ creates a new SSZUInt8 from SSZ format.
func (SSZUInt8) NewFromSSZ(buf []byte) (SSZUInt8, error) {
	if len(buf) != 1 {
		return 0, fmt.Errorf("invalid buffer length: expected 1, got %d", len(buf))
	}
	return SSZUInt8(buf[0]), nil
}

// -----------------------------

type SSZUInt16 uint16

// SizeSSZ returns the size of the uint16 in bytes.
func (u SSZUInt16) SizeSSZ() int {
	return 2
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u SSZUInt16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// NewFromSSZ creates a new SSZUInt16 from SSZ format.
func (SSZUInt16) NewFromSSZ(buf []byte) (SSZUInt16, error) {
	if len(buf) != 2 {
		return 0, fmt.Errorf("invalid buffer length: expected 2, got %d", len(buf))
	}
	return SSZUInt16(binary.LittleEndian.Uint16(buf)), nil
}

// -----------------------------

type SSZUInt32 uint32

// SizeSSZ returns the size of the uint32 in bytes.
func (u SSZUInt32) SizeSSZ() int {
	return 4
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u SSZUInt32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// NewFromSSZ creates a new SSZUInt32 from SSZ format.
func (SSZUInt32) NewFromSSZ(buf []byte) (SSZUInt32, error) {
	if len(buf) != 4 {
		return 0, fmt.Errorf("invalid buffer length: expected 4, got %d", len(buf))
	}
	return SSZUInt32(binary.LittleEndian.Uint32(buf)), nil
}

// -----------------------------

type SSZUInt64 uint64

// SizeSSZ returns the size of the uint64 in bytes.
func (u SSZUInt64) SizeSSZ() int {
	return 8
}

// MarshalSSZ marshals the uint64 into SSZ format.
func (u SSZUInt64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// NewFromSSZ creates a new SSZUInt64 from SSZ format.
func (SSZUInt64) NewFromSSZ(buf []byte) (SSZUInt64, error) {
	if len(buf) != 8 {
		return 0, fmt.Errorf("invalid buffer length: expected 8, got %d", len(buf))
	}
	return SSZUInt64(binary.LittleEndian.Uint64(buf)), nil
}

// -----------------------------

type SSZByte byte

// SizeSSZ returns the size of the byte slice in bytes.
func (b SSZByte) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the byte into SSZ format.
func (b SSZByte) MarshalSSZ() ([]byte, error) {
	return []byte{byte(b)}, nil
}

// NewFromSSZ creates a new SSZByte from SSZ format.
func (SSZByte) NewFromSSZ(buf []byte) (SSZByte, error) {
	if len(buf) != 1 {
		return 0, fmt.Errorf("invalid buffer length: expected 1, got %d", len(buf))
	}
	return SSZByte(buf[0]), nil
}
