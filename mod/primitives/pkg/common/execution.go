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

package common

import (
	"encoding"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"golang.org/x/crypto/sha3"
)

var (
	_ encoding.TextMarshaler   = (*ExecutionHash)(nil)
	_ encoding.TextUnmarshaler = (*ExecutionHash)(nil)
	_ json.Marshaler           = (*ExecutionHash)(nil)
	_ json.Unmarshaler         = (*ExecutionHash)(nil)

	_ encoding.TextMarshaler   = (*ExecutionAddress)(nil)
	_ encoding.TextUnmarshaler = (*ExecutionAddress)(nil)
	_ json.Marshaler           = (*ExecutionAddress)(nil)
	_ json.Unmarshaler         = (*ExecutionAddress)(nil)
)

/* -------------------------------------------------------------------------- */
/*                                ExecutionHash                               */
/* -------------------------------------------------------------------------- */

// ExecutionHash represents the 32 byte Keccak256 hash of arbitrary data.
// We use this type to represent hashes of things that come from the execution
// layer.
type ExecutionHash [32]byte

// NewExecutionHashFromHex creates a new hash from a hex string.
func NewExecutionHashFromHex(input string) ExecutionHash {
	return ExecutionHash(hex.MustToBytes(input))
}

// Hex converts a hash to a hex string.
func (h ExecutionHash) Hex() string { return string(hex.EncodeBytes(h[:])) }

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h ExecutionHash) String() string {
	return h.Hex()
}

// MarshalText returns the hex representation of h.
func (h ExecutionHash) MarshalText() ([]byte, error) {
	return hex.EncodeBytes(h[:]), nil
}

// UnmarshalText parses a hash in hex syntax.
func (h *ExecutionHash) UnmarshalText(input []byte) error {
	return hex.DecodeFixedText(input, h[:])
}

// MarshalJSON returns the JSON representation of h.
func (h ExecutionHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Hex())
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *ExecutionHash) UnmarshalJSON(input []byte) error {
	return hex.DecodeFixedJSON(input, h[:])
}

/* -------------------------------------------------------------------------- */
/*                              ExecutionAddress                              */
/* -------------------------------------------------------------------------- */

// ExecutionAddress represents a 20-byte Ethereum address.
// We use this type to represent addresses that come from the execution layer.
// It is EIP-55 checksummed and compliant.
type ExecutionAddress [20]byte

// NewExecutionAddressFromHex creates a new address from a hex string.
func NewExecutionAddressFromHex(input string) ExecutionAddress {
	return ExecutionAddress(hex.MustToBytes(input))
}

// Hex converts an address to a hex string.
func (a ExecutionAddress) Hex() string { return string(a.checksumHex()) }

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (a ExecutionAddress) String() string {
	return a.Hex()
}

// MarshalText returns the hex representation of a.
func (a ExecutionAddress) MarshalText() ([]byte, error) {
	return []byte(a.Hex()), nil
}

// UnmarshalText parses an address in hex syntax.
func (a *ExecutionAddress) UnmarshalText(input []byte) error {
	return hex.DecodeFixedText(input, a[:])
}

// MarshalJSON returns the JSON representation of a.
func (a ExecutionAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Hex())
}

// UnmarshalJSON parses an address in hex syntax.
func (a *ExecutionAddress) UnmarshalJSON(input []byte) error {
	return a.UnmarshalText(input[1 : len(input)-1])
}

// checksumHex returns the checksummed hex representation of a.
func (a *ExecutionAddress) checksumHex() []byte {
	buf := hex.EncodeBytes(a[:])

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		//nolint:mnd // todo fix.
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte >>= 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf
}
