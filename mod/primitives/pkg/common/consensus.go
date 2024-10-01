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
	"github.com/berachain/beacon-kit/mod/chain-spec/pkg/chain"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

//nolint:lll
type (
	// Bytes32 defines the commonly used 32-byte array.
	Bytes32 = bytes.B32

	// ChainSpec defines an interface for chain-specific parameters.
	ChainSpec = chain.Spec[DomainType, math.Epoch, ExecutionAddress, math.Slot, any]

	// Domain as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	//nolint:lll
	Domain = bytes.B32

	// DomainType as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	//nolint:lll
	DomainType = bytes.B4

	// Hash32 as er the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Hash32 = bytes.B32

	// Version as per the Ethereum 2.0 specification.
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	Version = bytes.B4

	// ForkDigest as per the Ethereum 2.0 Specification:
	// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#custom-types
	ForkDigest = bytes.B4
)

/* -------------------------------------------------------------------------- */
/*                                    Root                                    */
/* -------------------------------------------------------------------------- */

// Root represents a 32-byte Merkle root.
// We use this type to represent roots that come from the consensus layer.
type Root [32]byte

// NewRootFromHex creates a new root from a hex string.
func NewRootFromHex(input string) (Root, error) {
	val, err := hex.ToBytes(input)
	if err != nil {
		return Root{}, err
	}
	return Root(val), nil
}

// NewRootFromBytes creates a new root from a byte slice.
func NewRootFromBytes(input []byte) Root {
	var root Root
	copy(root[:], input)
	return root
}

// Hex converts a root to a hex string.
func (r Root) Hex() string { return string(hex.EncodeBytes(r[:])) }

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (r Root) String() string {
	return r.Hex()
}

// MarshalText returns the hex representation of r.
func (r Root) MarshalText() ([]byte, error) {
	return []byte(r.Hex()), nil
}

// UnmarshalText parses a root in hex syntax.
func (r *Root) UnmarshalText(input []byte) error {
	var err error
	*r, err = NewRootFromHex(string(input))
	return err
}

// MarshalJSON returns the JSON representation of r.
func (r Root) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Hex())
}

// UnmarshalJSON parses a root in hex syntax.
func (r *Root) UnmarshalJSON(input []byte) error {
	return r.UnmarshalText(input[1 : len(input)-1])
}
