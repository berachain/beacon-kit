// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package http_test

import (
	"net/http"
	"testing"

	iohttp "github.com/berachain/beacon-kit/io/http"
	"github.com/berachain/beacon-kit/io/jwt"
	"github.com/ethereum/go-ethereum/node"
	"github.com/stretchr/testify/require"
)

func TestNewHeaderWithJWT(t *testing.T) {
	tests := []struct {
		name      string
		jwtSecret string
		wantErr   bool
	}{
		{
			name: "valid JWT secret w/0x",
			//nolint:lll // test case.
			jwtSecret: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr:   false,
		},
		{
			name: "valid JWT secret w/o 0x",
			//nolint:lll // test case.
			jwtSecret: "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr:   false,
		},
		{
			name:      "empty JWT secret",
			jwtSecret: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtSecret, err := jwt.NewFromHex(tt.jwtSecret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to create JWT secret: %v", err)
				}
				return
			}

			var headers http.Header
			if tt.wantErr {
				require.PanicsWithValue(
					t,
					"http.Header is nil",
					func() { headers = iohttp.NewHeaderWithJWT(jwtSecret) },
					"Expected panic for nil JWT secret",
				)
			} else {
				require.NotPanics(t,
					func() {
						headers = iohttp.NewHeaderWithJWT(jwtSecret)
					}, "Unexpected panic for valid JWT secret",
				)
				require.NotNil(t, headers)
				require.IsType(t, http.Header{}, headers)
				fn := node.NewJWTAuth(*jwtSecret)
				err = fn(headers)
				require.NoError(t, err)
				authHeader := headers.Get("Authorization")
				require.NotEmpty(t, authHeader)
				require.Contains(t, authHeader, "Bearer ")
				require.Regexp(t,
					`^Bearer [A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+$`, authHeader,
					"Authorization header does not match expected JWT format",
				)
			}
		})
	}
}

func TestAddJWTHeader(t *testing.T) {
	secret, err := jwt.NewFromHex(
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	)
	require.NoError(t, err)

	header := http.Header{}
	err = iohttp.AddJWTHeader(header, secret)
	require.NoError(t, err)

	authHeader := header.Get("Authorization")
	require.NotEmpty(t, authHeader)
	require.Contains(t, authHeader, "Bearer ")

	// Test with nil header
	err = iohttp.AddJWTHeader(nil, secret)
	require.EqualError(t, err, "http.Header is nil")
}
