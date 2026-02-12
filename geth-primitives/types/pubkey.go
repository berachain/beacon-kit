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

package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ExecutionPubkey represents a 48-byte BLS12-381 public key on the execution layer.
// JSON and text serialization use 0x-prefixed hex strings.
type ExecutionPubkey [constants.BLSPubkeyLength]byte

// Bytes returns a copy of the underlying byte slice.
func (p ExecutionPubkey) Bytes() []byte { return p[:] }

// String returns the hex-encoded string representation of the pubkey.
func (p ExecutionPubkey) String() string { return hexutil.Encode(p[:]) }

// Format implements fmt.Formatter.
// Pubkey supports the %v, %s, %q, %x, %X and %d format verbs.
func (p ExecutionPubkey) Format(s fmt.State, c rune) {
	hexb := make([]byte, 2+len(p)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], p[:])

	switch c {
	case 'x', 'X':
		if !s.Flag('#') {
			hexb = hexb[2:]
		}
		if c == 'X' {
			hexb = bytes.ToUpper(hexb)
		}
		fallthrough
	case 'v', 's':
		s.Write(hexb)
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(hexb)
		s.Write(q)
	case 'd':
		fmt.Fprint(s, ([len(p)]byte)(p))
	default:
		fmt.Fprintf(s, "%%!%c(pubkey=%x)", c, p)
	}
}

// MarshalText encodes the pubkey as a 0x-prefixed hex string.
func (p ExecutionPubkey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(p[:]).MarshalText()
}

// UnmarshalText decodes a 0x-prefixed hex string into the pubkey.
func (p *ExecutionPubkey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Pubkey", input, p[:])
}

// UnmarshalJSON decodes a JSON string containing the 0x-prefixed hex pubkey.
func (p *ExecutionPubkey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(ExecutionPubkey{}), input, p[:])
}
