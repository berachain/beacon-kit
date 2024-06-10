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
	"reflect"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
)

func TestNewFromHex(t *testing.T) {
	wantValid := jwt.Secret(
		common.FromHex(
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
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromHex() = %v, want %v", got, tt.want)
			}
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
				common.FromHex(
					"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				),
			),
			want: "0x123456**********************************************************",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.secret.String(); got != tt.want {
				t.Errorf("Secret.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRandom(t *testing.T) {
	secret, err := jwt.NewRandom()
	if err != nil {
		t.Errorf("NewRandom() error = %v, wantErr %v", err, false)
	}
	if len(secret) == 0 {
		t.Errorf("NewRandom() generated an empty secret")
	}

	if len(secret.Bytes()) != 32 {
		t.Errorf(
			"NewRandom() generated a secret of incorrect length: got %d, want %d",
			len(secret.Bytes()),
			32,
		)
	}
}

func TestSecretBytes(t *testing.T) {
	expectedLength := 32 // Assuming the secret is expected to be 32 bytes long
	secret, _ := jwt.NewRandom()
	bytes := secret.Bytes()
	if len(bytes) != expectedLength {
		t.Errorf("Bytes() length = %d, want %d", len(bytes), expectedLength)
	}
}

func TestSecretHexWithFixedInput(t *testing.T) {
	expectedHex := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	// Since the secret is 32 bytes, its hex representation should be 64
	// characters
	// long
	expectedHexLength := 64
	secret, err := jwt.NewFromHex(expectedHex)
	if err != nil {
		t.Fatalf("NewFromHex() error = %v", err)
	}
	hexStr := secret.Hex()
	if hexStr != expectedHex {
		t.Errorf("Hex() = %s, want %s", hexStr, expectedHex)
	}

	// Check if the hex string is of the expected length and format.
	if len(hexStr) != expectedHexLength+2 {
		t.Errorf("Hex() length = %d, want %d", len(hexStr), expectedHexLength)
	}

	// Strip the '0x' prefix and check if the remaining string is valid hex.
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr) != expectedHexLength {
		t.Errorf(
			"Hex() length after stripping '0x' = %d, want %d",
			len(hexStr), expectedHexLength)
	}

	if !jwt.HexRegexp.MatchString(hexStr) {
		t.Errorf(
			"Hex() output does not match hexadecimal format, got: %s", hexStr,
		)
	}
}

func TestSecretRoundTripEncoding(t *testing.T) {
	originalSecret, err := jwt.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v, wantErr %v", err, false)
	}

	// Encode the original secret to hex string
	encodedSecret := hex.FromBytes(originalSecret.Bytes())

	// Decode the hex string back to secret
	decodedSecret, err := jwt.NewFromHex(encodedSecret.Unwrap())
	if err != nil {
		t.Fatalf("NewFromHex() error = %v", err)
	}

	// Compare the original and decoded secrets
	if !reflect.DeepEqual(originalSecret, decodedSecret) {
		t.Errorf(
			"Round trip encoding failed. Original: %v, Decoded: %v",
			originalSecret, decodedSecret,
		)
	}
}
