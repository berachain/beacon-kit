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

package ssz_test

import (
	"testing"

	ssz "github.com/itsdevbear/bolaris/encoding/ssz"
	"github.com/protolambda/ztyp/tree"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
)

func TestWithrawalSliceRoot(t *testing.T) {
	tests := []struct {
		name  string
		input []*enginev1.Withdrawal
		want  tree.Root
	}{
		{
			name:  "empty",
			input: make([]*enginev1.Withdrawal, 0),
			want: tree.Root{0x79, 0x29, 0x30, 0xbb, 0xd5, 0xba,
				0xac, 0x43, 0xbc, 0xc7, 0x98, 0xee, 0x49, 0xaa,
				0x81, 0x85, 0xef, 0x76, 0xbb, 0x3b, 0x44, 0xba,
				0x62, 0xb9, 0x1d, 0x86, 0xae, 0x56, 0x9e, 0x4b,
				0xb5, 0x35},
		},
		{
			name: "non-empty",
			input: []*enginev1.Withdrawal{{
				Index:          123,
				ValidatorIndex: 123123,
				Address:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				Amount:         50,
			},
			},
			want: tree.Root{0x10, 0x34, 0x29, 0xd1, 0x34, 0x30,
				0xa0, 0x1c, 0x4, 0xdd, 0x3, 0xed, 0xe6, 0xa6,
				0x33, 0xb2, 0xc9, 0x24, 0x23, 0x5c, 0x43, 0xca,
				0xb2, 0x32, 0xaa, 0xed, 0xfe, 0xd5, 0x9, 0x78,
				0xd1, 0x6f},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ssz.WithdrawalsRoot(tt.input, 16)
			require.NoError(t, err)
			require.DeepSSZEqual(t, tt.want, got)
		})
	}
}
