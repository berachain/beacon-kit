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

//nolint:mnd // lots of magic numbers here.
package ssz

import (
	"encoding/binary"
	"fmt"
)

/* -------------------------------------------------------------------------- */
/*                                    Bool                                    */
/* -------------------------------------------------------------------------- */

type Bool bool

// SizeSSZ returns the size of the bool in bytes.
func (b Bool) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the bool into SSZ format.
func (b Bool) MarshalSSZ() ([]byte, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// NewFromSSZ creates a new Bool from SSZ format.
func (Bool) NewFromSSZ(buf []byte) (Bool, error) {
	if len(buf) != 1 {
		return false, fmt.Errorf(
			"invalid buffer length: expected 1, got %d",
			len(buf),
		)
	}
	return Bool(buf[0] != 0), nil
}

// HashTreeRoot returns the hash tree root of the bool.
func (b Bool) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	if b {
		buf[0] = 1
	}
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the bool.
func (Bool) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U8                                     */
/* -------------------------------------------------------------------------- */

type U8 uint8

// SizeSSZ returns the size of the uint8 in bytes.
func (u U8) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the uint8 into SSZ format.
func (u U8) MarshalSSZ() ([]byte, error) {
	return []byte{byte(u)}, nil
}

// NewFromSSZ creates a new U8 from SSZ format.
func (U8) NewFromSSZ(buf []byte) (U8, error) {
	if len(buf) != 1 {
		return 0, fmt.Errorf(
			"invalid buffer length: expected 1, got %d",
			len(buf),
		)
	}

	//#nosec:G701 // the check above protects against overflow.
	return U8(buf[0]), nil
}

// HashTreeRoot returns the hash tree root of the uint8.
func (u U8) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	buf[0] = byte(u)
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the uint8.
func (U8) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U16                                     */
/* -------------------------------------------------------------------------- */

type U16 uint16

// SizeSSZ returns the size of the uint16 in bytes.
func (u U16) SizeSSZ() int {
	return 2
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u U16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// NewFromSSZ creates a new U16 from SSZ format.
func (U16) NewFromSSZ(buf []byte) (U16, error) {
	if len(buf) != 2 {
		return 0, fmt.Errorf(
			"invalid buffer length: expected 2, got %d",
			len(buf),
		)
	}
	return U16(binary.LittleEndian.Uint16(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint16.
func (u U16) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint16(buf[:2], uint16(u))
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the uint16.
func (U16) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U32                                    */
/* -------------------------------------------------------------------------- */

type U32 uint32

// SizeSSZ returns the size of the uint32 in bytes.
func (u U32) SizeSSZ() int {
	return 4
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u U32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// NewFromSSZ creates a new U32 from SSZ format.
func (U32) NewFromSSZ(buf []byte) (U32, error) {
	if len(buf) != 4 {
		return 0, fmt.Errorf(
			"invalid buffer length: expected 4, got %d",
			len(buf),
		)
	}
	return U32(binary.LittleEndian.Uint32(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint32.
func (u U32) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint32(buf[:4], uint32(u))
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the uint32.
func (U32) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U64                                    */
/* -------------------------------------------------------------------------- */

type U64 uint64

// SizeSSZ returns the size of the uint64 in bytes.
func (u U64) SizeSSZ() int {
	return 8
}

// MarshalSSZ marshals the uint64 into SSZ format.
func (u U64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// NewFromSSZ creates a new U64 from SSZ format.
func (U64) NewFromSSZ(buf []byte) (U64, error) {
	if len(buf) != 8 {
		return 0, fmt.Errorf(
			"invalid buffer length: expected 8, got %d",
			len(buf),
		)
	}
	return U64(binary.LittleEndian.Uint64(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint64.
func (u U64) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	binary.LittleEndian.PutUint64(buf[:8], uint64(u))
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the uint64.
func (U64) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                    Byte                                    */
/* -------------------------------------------------------------------------- */

type Byte byte

// SizeSSZ returns the size of the byte slice in bytes.
func (b Byte) SizeSSZ() int {
	return 1
}

// MarshalSSZ marshals the byte into SSZ format.
func (b Byte) MarshalSSZ() ([]byte, error) {
	return []byte{byte(b)}, nil
}

// NewFromSSZ creates a new Byte from SSZ format.
func (Byte) NewFromSSZ(buf []byte) (Byte, error) {
	if len(buf) != 1 {
		return 0, fmt.Errorf(
			"invalid buffer length: expected 1, got %d",
			len(buf),
		)
	}
	return Byte(buf[0]), nil
}

// HashTreeRoot returns the hash tree root of the byte.
func (b Byte) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, 32)
	buf[0] = byte(b)
	return [32]byte(buf), nil
}

// ChunkCount returns the number of chunks required to store the byte.
func (Byte) ChunkCount() uint64 {
	return 1
}
