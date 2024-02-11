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
	"github.com/itsdevbear/bolaris/io/jwt"
)

func TestNewFromHex(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		want    jwt.Secret
		wantErr bool
	}{
		{
			name:   "valid hex string w/ 0x prefix",
			hexStr: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			want: jwt.Secret(
				common.FromHex("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			),
			wantErr: false,
		},
		{
			name:   "valid hex string no 0x prefix",
			hexStr: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			want: jwt.Secret(
				common.FromHex("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			),
			wantErr: false,
		},
		{
			name:    "invalid hex string",
			hexStr:  "0x123",
			want:    jwt.Secret{},
			wantErr: true,
		},
		{
			name:    "empty hex string",
			hexStr:  "",
			want:    jwt.Secret{},
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
