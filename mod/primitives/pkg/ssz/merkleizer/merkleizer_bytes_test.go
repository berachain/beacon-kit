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

package merkleizer_test

import (
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
)

var prysmConsistencyTests = []struct {
	name      string
	slice     []byte
	maxLength uint64
	want      [32]byte
	wantErr   bool
}{
	{
		name:  "nil",
		slice: nil,
		want: [32]byte{
			245, 165, 253, 66, 209, 106, 32, 48, 39, 152, 239, 110, 211, 9,
			151, 155, 67, 0, 61, 35, 32, 217, 240, 232, 234, 152, 49, 169,
			39, 89, 251, 75,
		},
	},
	{
		name:  "empty",
		slice: []byte{},
		want: [32]byte{
			245, 165, 253, 66, 209, 106, 32, 48, 39, 152, 239, 110, 211, 9,
			151, 155, 67, 0, 61, 35, 32, 217, 240, 232, 234, 152, 49, 169,
			39, 89, 251, 75,
		},
	},
	{
		name:  "byte slice 3 values",
		slice: []byte{1, 2, 3},
		want: [32]byte{
			20, 159, 26, 252, 247, 204, 44, 159, 161, 135, 211, 195, 106,
			59, 220, 149, 199, 163, 228, 155, 113, 118, 64, 126, 173, 223,
			102, 1, 241, 158, 164, 185,
		},
	},
	{
		name: "byte slice 32 values",
		slice: []byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18,
			19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
		},
		want: [32]byte{
			7, 30, 46, 77, 237, 240, 59, 126, 232, 232, 232, 6, 145, 210,
			31, 18, 117, 12, 217, 40, 204, 141, 90, 236, 241, 128, 221, 45,
			126, 39, 39, 202,
		},
	},
	{
		name:    "over max length",
		slice:   make([]byte, constants.RootLength+1),
		want:    [32]byte{},
		wantErr: true,
	},
}

func TestByteSliceMerkleization(t *testing.T) {
	merkleizer := merkleizer.New[[32]byte, ssz.Byte]()

	for _, tt := range prysmConsistencyTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.maxLength == 0 {
				tt.maxLength = constants.RootLength
			}
			byteList := ssz.ByteListFromBytes(tt.slice, tt.maxLength)
			got, err := byteList.HashTreeRootWith(merkleizer)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ByteSliceRoot() error = %v, wantErr %v", err, tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteSliceRoot() got = %v, want %v", got, tt.want)
			}
		})
	}
}
