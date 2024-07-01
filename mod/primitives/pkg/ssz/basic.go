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


package ssz

import (
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/types"
	"github.com/holiman/uint256"
)

/* -------------------------------------------------------------------------- */
/*                                Type Definitions                            */
/* -------------------------------------------------------------------------- */

// Ensure types implement types.SSZType.
var (
	_ types.SSZType[Bool]  = (*Bool)(nil)
	_ types.SSZType[U8]    = (*U8)(nil)
	_ types.SSZType[U16]   = (*U16)(nil)
	_ types.SSZType[U32]   = (*U32)(nil)
	_ types.SSZType[U64]   = (*U64)(nil)
	_ types.SSZType[*U256] = (*U256)(nil)
	_ types.SSZType[Byte]  = (*Byte)(nil)
)

type (
	Bool bool
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	U256 uint256.Int
	Byte byte
)

/* -------------------------------------------------------------------------- */
/*                                    Bool                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the bool in bytes.
func (Bool) SizeSSZ() int {
	return constants.BoolSize
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
	if len(buf) != constants.BoolSize {
		return false, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.BoolSize,
			len(buf),
		)
	}
	return Bool(buf[0] != 0), nil
}

// HashTreeRoot returns the hash tree root of the bool.
func (b Bool) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	if b {
		buf[0] = 1
	}
	return [constants.BytesPerChunk]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (Bool) IsFixed() bool {
	return true
}

// Type returns the type of the bool.
func (Bool) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the bool.
func (Bool) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U8                                     */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint8 in bytes.
func (U8) SizeSSZ() int {
	return constants.U8Size
}

// MarshalSSZ marshals the uint8 into SSZ format.
func (u U8) MarshalSSZ() ([]byte, error) {
	return []byte{byte(u)}, nil
}

// NewFromSSZ creates a new U8 from SSZ format.
func (U8) NewFromSSZ(buf []byte) (U8, error) {
	if len(buf) != constants.U8Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U8Size,
			len(buf),
		)
	}

	//#nosec:G701 // the check above protects against overflow.
	return U8(buf[0]), nil
}

// HashTreeRoot returns the hash tree root of the uint8.
func (u U8) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	buf[0] = byte(u)
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U8) IsFixed() bool {
	return true
}

// Type returns the type of the U8.
func (U8) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint8.
func (U8) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U16                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint16 in bytes.
func (U16) SizeSSZ() int {
	return constants.U16Size
}

// MarshalSSZ marshals the uint16 into SSZ format.
func (u U16) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U16Size)
	binary.LittleEndian.PutUint16(buf, uint16(u))
	return buf, nil
}

// NewFromSSZ creates a new U16 from SSZ format.
func (U16) NewFromSSZ(buf []byte) (U16, error) {
	if len(buf) != constants.U16Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U16Size,
			len(buf),
		)
	}
	return U16(binary.LittleEndian.Uint16(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint16.
func (u U16) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint16(buf[:constants.U16Size], uint16(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U16) IsFixed() bool {
	return true
}

// Type returns the type of the U16.
func (U16) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint16.
func (U16) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U32                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint32 in bytes.
func (U32) SizeSSZ() int {
	return constants.U32Size
}

// MarshalSSZ marshals the uint32 into SSZ format.
func (u U32) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U32Size)
	binary.LittleEndian.PutUint32(buf, uint32(u))
	return buf, nil
}

// NewFromSSZ creates a new U32 from SSZ format.
func (U32) NewFromSSZ(buf []byte) (U32, error) {
	if len(buf) != constants.U32Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U32Size,
			len(buf),
		)
	}
	return U32(binary.LittleEndian.Uint32(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint32.
func (u U32) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint32(buf[:constants.U32Size], uint32(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U32) IsFixed() bool {
	return true
}

// Type returns the type of the U32.
func (U32) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint32.
func (U32) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                     U64                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the uint64 in bytes.
func (U64) SizeSSZ() int {
	return constants.U64Size
}

// MarshalSSZ marshals the uint64 into SSZ format.
func (u U64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, constants.U64Size)
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// NewFromSSZ creates a new U64 from SSZ format.
func (U64) NewFromSSZ(buf []byte) (U64, error) {
	if len(buf) != constants.U64Size {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U64Size,
			len(buf),
		)
	}
	return U64(binary.LittleEndian.Uint64(buf)), nil
}

// HashTreeRoot returns the hash tree root of the uint64.
func (u U64) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	binary.LittleEndian.PutUint64(buf[:constants.U64Size], uint64(u))
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (U64) IsFixed() bool {
	return true
}

// Type returns the type of the U64.
func (U64) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the uint64.
func (U64) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                    U256                                    */
/* -------------------------------------------------------------------------- */

func NewU256FromUint64(v uint64) *U256 {
	return (*U256)(uint256.NewInt(0).SetUint64(v))
}

// SizeSSZ returns the size of the U256 in bytes.
func (U256) SizeSSZ() int {
	return constants.U256Size
}

// MarshalSSZ marshals the U256 into SSZ format.
func (u *U256) MarshalSSZ() ([]byte, error) {
	return (*uint256.Int)(u).MarshalSSZ()
}

// NewFromSSZ creates a new U256 from SSZ format.
func (U256) NewFromSSZ(buf []byte) (*U256, error) {
	if len(buf) != constants.U256Size {
		return nil, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.U256Size,
			len(buf),
		)
	}
	u := new(uint256.Int)
	return (*U256)(u), u.UnmarshalSSZ(buf)
}

// HashTreeRoot returns the hash tree root of the U256.
func (u *U256) HashTreeRoot() ([32]byte, error) {
	return (*uint256.Int)(u).HashTreeRoot()
}

// IsFixed returns true if the U256 is fixed size.
func (*U256) IsFixed() bool {
	return true
}

// Type returns the type of the U256.
func (*U256) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the U256.
func (*U256) ChunkCount() uint64 {
	return 1
}

/* -------------------------------------------------------------------------- */
/*                                    Byte                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the byte slice in bytes.
func (Byte) SizeSSZ() int {
	return constants.ByteSize
}

// MarshalSSZ marshals the byte into SSZ format.
func (b Byte) MarshalSSZ() ([]byte, error) {
	return []byte{byte(b)}, nil
}

// NewFromSSZ creates a new Byte from SSZ format.
func (Byte) NewFromSSZ(buf []byte) (Byte, error) {
	if len(buf) != constants.ByteSize {
		return 0, fmt.Errorf(
			"invalid buffer length: expected %d, got %d",
			constants.ByteSize,
			len(buf),
		)
	}
	return Byte(buf[0]), nil
}

// HashTreeRoot returns the hash tree root of the byte.
func (b Byte) HashTreeRoot() ([32]byte, error) {
	buf := make([]byte, constants.BytesPerChunk)
	buf[0] = byte(b)
	return [32]byte(buf), nil
}

// IsFixed returns true if the bool is fixed size.
func (Byte) IsFixed() bool {
	return true
}

// Type returns the type of the Byte.
func (Byte) Type() types.Type {
	return types.Basic
}

// ChunkCount returns the number of chunks required to store the byte.
func (Byte) ChunkCount() uint64 {
	return 1
}
