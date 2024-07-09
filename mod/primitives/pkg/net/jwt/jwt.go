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

package jwt

import (
	"crypto/rand"
	"regexp"
	"strings"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

// HexRegexp is a regular expression to match hexadecimal characters.
var HexRegexp = regexp.MustCompile(`^(?:0x)?[0-9a-fA-F]*$`)

// EthereumJWTLength defines the length of the JWT byte array to be 32 bytes as
// defined the Engine API specification.
// https://github.com/ethereum/execution-apis/blob/main/src/engine/authentication.md
//
//nolint:lll // link.
const EthereumJWTLength = 32

// Secret represents a JSON Web Token as a fixed-size byte array.
type Secret [EthereumJWTLength]byte

// NewFromHex creates a new JWT secret from a hexadecimal string.
func NewFromHex(hexStr string) (*Secret, error) {
	// Ensure the hex string contains only hexadecimal characters.
	if !HexRegexp.MatchString(hexStr) {
		return nil, ErrContainsIllegalCharacter
	}

	// Convert the hex string to a byte array.
	bz := gethprimitives.FromHex(hexStr)
	if bz == nil || len(bz) != EthereumJWTLength {
		return nil, ErrLengthMismatch
	}
	s := Secret(bz)
	return &s, nil
}

// NewRandom creates a new random JWT secret.
func NewRandom() (*Secret, error) {
	secret := make([]byte, EthereumJWTLength)
	// We don't need to check n since:
	// n == len(b) if and only if err == nil.
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return NewFromHex(hex.FromBytes(secret).Unwrap())
}

// String returns the JWT secret as a string with the first 8 characters
// visible and the rest masked out for security.
func (s *Secret) String() string {
	secret := hex.FromBytes(s[:]).Unwrap()
	return secret[:8] + strings.Repeat("*", len(secret[8:]))
}

// Hex returns the JWT secret as a hexadecimal string.
func (s *Secret) Hex() string {
	return hex.FromBytes(s[:]).Unwrap()
}

// Bytes returns the JWT secret as a byte array.
func (s *Secret) Bytes() []byte {
	return s[:]
}
