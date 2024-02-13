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
	"reflect"
	"testing"

	"github.com/itsdevbear/bolaris/encoding/ssz"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

func TestTransactionsRoot(t *testing.T) {
	tests := []struct {
		name    string
		txs     []ssz.Bytes
		want    [32]byte
		wantErr bool
	}{
		{
			name: "nil",
			txs:  nil,
			want: [32]byte{127, 254, 36, 30, 166, 1, 135, 253, 176, 24,
				123, 250, 34, 222, 53, 209, 249, 190, 215, 171, 6, 29,
				148, 1, 253, 71, 227, 74, 84, 251, 237, 225},
		},
		{
			name: "empty",
			txs:  []ssz.Bytes{},
			want: [32]byte{127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123,
				250, 34, 222, 53, 209, 249, 190, 215, 171, 6, 29, 148, 1,
				253, 71, 227, 74, 84, 251, 237, 225},
		},
		{
			name: "one tx",
			txs:  []ssz.Bytes{{1, 2, 3}},
			want: [32]byte{102, 209, 140, 87, 217, 28, 68, 12, 133, 42, 77, 136,
				191, 18, 234, 105, 166, 228, 216, 235, 230, 95, 200, 73,
				85, 33, 134, 254, 219, 97, 82, 209},
		},
		{
			name: "max txs",
			txs: func() []ssz.Bytes {
				var txs []ssz.Bytes
				for i := 0; i < primitives.MaxTxsPerPayloadLength; i++ {
					txs = append(txs, []byte{})
				}
				return txs
			}(),
			want: [32]byte{13, 66, 254, 206, 203, 58, 48, 133, 78, 218, 48, 231, 120,
				90, 38, 72, 73, 137, 86, 9, 31, 213, 185, 101, 103, 144, 0,
				236, 225, 57, 47, 244},
		},
		{
			name: "exceed max txs",
			txs: func() []ssz.Bytes {
				var txs []ssz.Bytes
				for i := 0; i < primitives.MaxTxsPerPayloadLength+1; i++ {
					txs = append(txs, []byte{})
				}
				return txs
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ssz.TransactionsRoot(tt.txs)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionsRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionsRoot() got = %v, want %v", got, tt.want)
			}
		})
	}
}
