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

//nolint:lll // lots of hex strings.
package jwt_test

import (
	"strings"
	"testing"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	"github.com/stretchr/testify/require"
)

func TestNewFromHex(t *testing.T) {
	wantValid := jwt.Secret(
		gethprimitives.FromHex(
			"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		),
	)
	tests := []struct {
		name    string
		hexStr  string
		want    *jwt.Secret
		wantErr bool
	}{
		{
			name: "valid hex string w/ 0x prefix",

			hexStr:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			want:    &(wantValid),
			wantErr: false,
		},
		{
			name: "valid hex string no 0x prefix",

			hexStr:  "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			want:    &(wantValid),
			wantErr: false,
		},
		{
			name:    "invalid hex string",
			hexStr:  "0x123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty hex string",
			hexStr:  "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jwt.NewFromHex(tt.hexStr)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSecretString(t *testing.T) {
	tests := []struct {
		name   string
		secret jwt.Secret
		want   string
	}{
		{
			name: "mask secret correctly",
			secret: jwt.Secret(
				gethprimitives.FromHex(
					"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				),
			),
			want: "0x123456**********************************************************",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.secret.String(),
				"Secret.String() mismatch",
			)
		})
	}
}

func TestNewRandom(t *testing.T) {
	secret, err := jwt.NewRandom()
	require.NoError(t, err, "NewRandom() error")
	require.Len(t, secret.Bytes(), 32, "NewRandom() length mismatch")
}

func TestSecretBytes(t *testing.T) {
	expectedLength := 32 // Assuming the secret is expected to be 32 bytes long
	secret, _ := jwt.NewRandom()
	bytes := secret.Bytes()
	require.Len(t, bytes, expectedLength, "Bytes() length mismatch")
}

func TestSecretHexWithFixedInput(t *testing.T) {
	expectedHex := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	expectedHexLength := 64

	secret, err := jwt.NewFromHex(expectedHex)
	require.NoError(t, err, "NewFromHex() error")

	hexStr := secret.Hex()
	require.Equal(t, expectedHex, hexStr, "Hex() output mismatch")

	// Check if the hex string is of the expected length and format.
	require.Len(t, hexStr, expectedHexLength+2, "Hex() length mismatch")

	// Strip the '0x' prefix and check if the remaining string is valid hex.
	hexStr = strings.TrimPrefix(hexStr, "0x")
	require.Len(
		t,
		hexStr,
		expectedHexLength,
		"Hex() length after stripping '0x' mismatch",
	)
	require.True(
		t,
		jwt.HexRegexp.MatchString(hexStr),
		"Hex() output does not match hexadecimal format",
	)
}

func TestSecretRoundTripEncoding(t *testing.T) {
	originalSecret, err := jwt.NewRandom()
	require.NoError(t, err, "NewRandom() error")

	// Encode the original secret to hex string
	encodedSecret := hex.FromBytes(originalSecret.Bytes())

	// Decode the hex string back to secret
	decodedSecret, err := jwt.NewFromHex(encodedSecret.Unwrap())
	require.NoError(t, err, "NewFromHex() error")

	// Compare the original and decoded secrets
	require.Equal(
		t,
		originalSecret,
		decodedSecret,
		"Round trip encoding failed",
	)
}
