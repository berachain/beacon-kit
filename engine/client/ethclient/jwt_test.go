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

package ethclient_test

import (
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/node"
	"github.com/itsdevbear/bolaris/engine/client/ethclient"
	"github.com/itsdevbear/bolaris/io/jwt"
	"github.com/stretchr/testify/require"
)

func TestBuildHeaders(t *testing.T) {
	tests := []struct {
		name      string
		jwtSecret string
		wantErr   bool
	}{
		{
			name: "valid JWT secret w/0x",
			//nolint:lll
			jwtSecret: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr:   false,
		},
		{
			name: "valid JWT secret w/o 0x",
			//nolint:lll
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
			// Create a new JWT secret
			jwtSecret, err := jwt.NewFromHex(tt.jwtSecret)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to create JWT secret: %v", err)
				} else {
					return
				}
			}

			// Create a new Eth1Client with the JWT secret
			client := &ethclient.Eth1Client{}
			if err = ethclient.WithJWTSecret(jwtSecret)(client); err != nil {
				if !tt.wantErr {
					t.Errorf("Failed to set JWT secret option: %v", err)
				}
			}
			headers, err := client.BuildHeaders()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
