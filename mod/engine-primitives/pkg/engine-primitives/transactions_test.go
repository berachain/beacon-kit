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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
)

// Tests consistency of our HashTreeRoot implementation with Prysm's.
// https://github.com/prysmaticlabs/prysm/blob/8070fc8ecef3c2003dea7e1fe6dba179ddf76151/encoding/ssz/htrutils_test.go#L67
//
//nolint:lll // link.
var prysmConsistencyTests = []struct {
	name    string
	txs     [][]byte
	want    [32]byte
	wantErr bool
}{
	{
		name: "nil",
		txs:  nil,
		want: [32]byte{
			127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123, 250, 34, 222,
			53, 209, 249, 190, 215, 171, 6, 29, 148, 1, 253, 71, 227, 74,
			84, 251, 237, 225,
		},
	},
	{
		name: "empty",
		txs:  [][]byte{},
		want: [32]byte{
			127, 254, 36, 30, 166, 1, 135, 253, 176, 24, 123, 250, 34, 222,
			53, 209, 249, 190, 215, 171, 6, 29, 148, 1, 253, 71, 227, 74,
			84, 251, 237, 225,
		},
	},
	{
		name: "3 non-nil txs",
		txs: [][]byte{
			[]byte("transaction1"),
			[]byte("transaction2"),
			[]byte("transaction3"),
		},
		want: [32]byte{
			139, 213, 123, 109, 253, 176, 23, 93, 101, 51, 142, 198, 119,
			250, 13, 242, 79, 219, 180, 165, 254, 181, 9, 178, 4, 253,
			110, 75, 50, 25, 17, 141,
		},
	},
	{
		name: "one tx",
		txs:  [][]byte{{1, 2, 3}},
		want: [32]byte{
			102, 209, 140, 87, 217, 28, 68, 12, 133, 42, 77, 136, 191, 18,
			234, 105, 166, 228, 216, 235, 230, 95, 200, 73, 85, 33, 134,
			254, 219, 97, 82, 209,
		},
	},
	{
		name: "max txs",
		txs: func() [][]byte {
			var txs [][]byte
			for range int(constants.MaxTxsPerPayload) {
				txs = append(txs, []byte{
					0x01,
				})
			}
			return txs
		}(),
		want: [32]byte{
			168, 19, 62, 29, 232, 106, 28, 81, 99,
			73, 236, 102, 94, 160, 44, 191, 122, 176,
			38, 39, 139, 100, 136, 5, 48, 242, 34, 31, 60, 104, 191, 171,
		},
	},
}

// NOTE: not testing legacy and Bartio transactions types
// (engineprimitives.Transactions and engine.primitivesBartioTransactions
// respectively) since those will be deprecated soon.
func TestProperTransactions(t *testing.T) {
	for _, tt := range prysmConsistencyTests {
		t.Run(tt.name, func(t *testing.T) {
			got := engineprimitives.Transactions(
				tt.txs,
			).HashTreeRoot()

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf(
						"TransactionsRoot() got = %v, want %v, off at byte %d",
						[32]byte(got), tt.want, i,
					)
					return
				}
			}
		})
	}
}
