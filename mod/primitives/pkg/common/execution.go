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
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type (
	ExecutionAddress = common.Address
)

var (
	_ encoding.TextMarshaler   = (*ExecutionHash)(nil)
	_ encoding.TextUnmarshaler = (*ExecutionHash)(nil)
	_ json.Marshaler           = (*ExecutionHash)(nil)
	_ json.Unmarshaler         = (*ExecutionHash)(nil)
)

// ExecutionHash represents the 32 byte Keccak256 hash of arbitrary data.
// We use this type to represent hashes of things that come from the execution
// layer.
type ExecutionHash [32]byte

// NewExecutionHashFromHex creates a new hash from a hex string.
func NewExecutionHashFromHex(hex string) ExecutionHash {
	return ExecutionHash(hexutil.MustDecode(hex))
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
