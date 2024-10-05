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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	"github.com/stretchr/testify/require"
)

func TestNewFromHex(t *testing.T) {
	wantValid := jwt.Secret(
		hex.MustToBytes(
			"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
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
			wantErr: true,
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
				hex.MustToBytes(
					"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
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
	encodedSecret := hex.EncodeBytes((originalSecret.Bytes()))

	// Decode the hex string back to secret
	decodedSecret, err := jwt.NewFromHex(encodedSecret)
	require.NoError(t, err, "NewFromHex() error")

	// Compare the original and decoded secrets
	require.Equal(
		t,
		originalSecret,
		decodedSecret,
		"Round trip encoding failed",
	)
}

func TestBuildSignedToken(t *testing.T) {
	secret, err := jwt.NewRandom()
	require.NoError(t, err, "NewRandom() error")

	token, err := secret.BuildSignedToken()
	require.NoError(t, err, "BuildSignedToken() error")
	require.NotEmpty(t, token, "BuildSignedToken() returned empty token")

	// Verify the token structure (header.payload.signature)
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3, "Token should have three parts")
}

func TestNewFromHexEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		wantErr bool
	}{
		{
			name:    "lowercase hex string",
			hexStr:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			wantErr: false,
		},
		{
			name:    "uppercase hex string",
			hexStr:  "0xABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890",
			wantErr: false,
		},
		{
			name:    "mixed case hex string",
			hexStr:  "0xaBcDeF1234567890aBcDeF1234567890aBcDeF1234567890aBcDeF1234567890",
			wantErr: false,
		},
		{
			name:    "invalid characters",
			hexStr:  "0x123G567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr: true,
		},
		{
			name:    "too short",
			hexStr:  "0x1234567890abcdef",
			wantErr: true,
		},
		{
			name:    "too long",
			hexStr:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwt.NewFromHex(tt.hexStr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHexRegexp(t *testing.T) {
	validHexStrings := []string{
		"0x1234567890abcdef",
		"1234567890ABCDEF",
		"0xABCDEF1234567890",
		"abcdef1234567890",
	}

	invalidHexStrings := []string{
		"0x123G567890abcdef",
		"GHIJKLMNOPQRSTUV",
		"0xABCDEF12345678@0",
		"abcdef123456789g",
	}

	for _, validHex := range validHexStrings {
		require.True(
			t,
			jwt.HexRegexp.MatchString(validHex),
			"Valid hex string not matched: %s",
			validHex,
		)
	}

	for _, invalidHex := range invalidHexStrings {
		require.False(
			t,
			jwt.HexRegexp.MatchString(invalidHex),
			"Invalid hex string matched: %s",
			invalidHex,
		)
	}
}

func TestSecretComparison(t *testing.T) {
	secret1, err := jwt.NewRandom()
	require.NoError(t, err, "NewRandom() error for secret1")

	secret2, err := jwt.NewRandom()
	require.NoError(t, err, "NewRandom() error for secret2")

	require.NotEqual(
		t,
		secret1,
		secret2,
		"Two random secrets should not be equal",
	)

	secret3 := *secret1
	require.Equal(
		t,
		secret1,
		&secret3,
		"Copied secret should be equal to original",
	)
}
