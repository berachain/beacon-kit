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

package sha256_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
	"github.com/stretchr/testify/require"
)

func Test_SafeMerkleizeVector(t *testing.T) {
	tests := []struct {
		name            string
		roots           [][32]byte
		maxRootsAllowed uint64
		expected        [32]byte
		wantErr         bool
	}{
		{
			name:            "empty roots list",
			roots:           make([][32]byte, 0),
			maxRootsAllowed: 16,
			expected:        tree.ZeroHashes[0],
			wantErr:         false,
		},
		{
			name:            "maxRootsAllowed is less than the number of roots",
			roots:           [][32]byte{{0x01}, {0x01}, {0x01}, {0x01}},
			maxRootsAllowed: 3,
			expected:        [32]byte{0x00},
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := sha256.SafeMerkleizeVector(tt.roots, tt.maxRootsAllowed)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.Equal(t, tt.expected, root)
		})
	}
}
