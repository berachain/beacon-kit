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

//nolint:testpackage // private functions.
package merkle

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/stretchr/testify/require"
)

func TestDepositTreeSnapshot_CalculateRoot(t *testing.T) {
	hasher := merkle.NewHasher[[32]byte](sha256.Hash)
	tests := []struct {
		name         string
		finalized    int
		depositCount uint64
		want         [32]byte
	}{
		{
			name:         "empty",
			finalized:    0,
			depositCount: 0,
			want:         [32]byte{215, 10, 35, 71, 49, 40, 92, 104, 4, 194, 164, 245, 103, 17, 221, 184, 200, 44, 153, 116, 15, 32, 120, 84, 137, 16, 40, 175, 52, 226, 126, 94},
		},
		{
			name:         "1 Finalized",
			finalized:    1,
			depositCount: 2,
			want:         [32]byte{36, 118, 154, 57, 217, 109, 145, 116, 238, 1, 207, 59, 187, 28, 69, 187, 70, 55, 153, 180, 15, 150, 37, 72, 140, 36, 109, 154, 212, 202, 47, 59},
		},
		{
			name:         "many finalised",
			finalized:    6,
			depositCount: 20,
			want:         [32]byte{210, 63, 57, 119, 12, 5, 3, 25, 139, 20, 244, 59, 114, 119, 35, 88, 222, 88, 122, 106, 239, 20, 45, 140, 99, 92, 222, 166, 133, 159, 128, 72},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var finalized [][32]byte
			for i := range tt.finalized {
				finalized = append(
					finalized,
					hexString(t, fmt.Sprintf("%064d", i)),
				)
			}
			ds := &DepositTreeSnapshot{
				finalized:    finalized,
				depositCount: tt.depositCount,
				hasher:       hasher,
			}
			root := ds.CalculateRoot()
			if got := root; !reflect.DeepEqual(got, tt.want) {
				require.Equal(t, tt.want, got)
			}
		})
	}
}
