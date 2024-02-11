// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package jwt_test

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itsdevbear/bolaris/io/jwt"
)

func TestNewFromHex(t *testing.T) {
	wantValid := jwt.Secret(
		common.FromHex("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
	)
	tests := []struct {
		name    string
		hexStr  string
		want    *jwt.Secret
		wantErr bool
	}{
		{
			name:    "valid hex string w/ 0x prefix",
			hexStr:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			want:    &(wantValid),
			wantErr: false,
		},
		{
			name:    "valid hex string no 0x prefix",
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
					"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
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
			len(secret.Bytes()), 32,
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
	for _, b := range bytes {
		if b == 0 {
			t.Errorf("Bytes() contains zero byte, want all bytes to be non-zero")
			break
		}
	}
}

func TestSecretRoundTripEncoding(t *testing.T) {
	originalSecret, err := jwt.NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() error = %v, wantErr %v", err, false)
	}

	// Encode the original secret to hex string
	encodedSecret := hexutil.Encode(originalSecret.Bytes())

	// Decode the hex string back to secret
	decodedSecret, err := jwt.NewFromHex(encodedSecret)
	if err != nil {
		t.Fatalf("NewFromHex() error = %v", err)
	}

	// Compare the original and decoded secrets
	if !reflect.DeepEqual(originalSecret, decodedSecret) {
		t.Errorf("Round trip encoding failed. Original: %v, Decoded: %v", originalSecret, decodedSecret)
	}
}
