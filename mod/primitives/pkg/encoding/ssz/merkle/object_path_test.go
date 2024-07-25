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

package merkle_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/stretchr/testify/require"
)

func Test_ObjectPath(t *testing.T) {
	nested := schema.DefineContainer(
		schema.NewField("bytes32", schema.B32()),
		schema.NewField("uint64", schema.U64()),
		schema.NewField("list_bytes32", schema.DefineList(schema.B32(), 10)),
		schema.NewField("bytes256", schema.DefineVector(schema.U8(), 256)),
	)
	root := schema.DefineContainer(
		schema.NewField("bytes32", schema.B32()),
		schema.NewField("uint32", schema.U32()),
		schema.NewField("list_uint64", schema.DefineList(schema.U64(), 1000)),
		schema.NewField("list_nested", schema.DefineList(nested, 1000)),
		schema.NewField("nested", nested),
		schema.NewField("vector_uint64", schema.DefineVector(schema.U64(), 40)),
	)

	cases := []struct {
		path   string
		gindex uint64
		offset uint8
		error  string
	}{
		// happy paths
		{path: "bytes32", gindex: 8},
		{path: "bytes32/3", gindex: 8, offset: 3},
		{path: "uint32", gindex: 9},
		{path: "list_nested", gindex: 11},
		{path: "list_nested/__len__", gindex: 2*11 + 1},
		{path: "list_nested/12", gindex: 11*2*1024 + 12},
		{path: "list_nested/12/uint64", gindex: (11*2*1024+12)*4 + 1},
		{path: "nested", gindex: 12},
		{path: "nested/uint64", gindex: 12*4 + 1},
		{path: "nested/bytes256", gindex: 12*4 + 3},
		{path: "nested/bytes256/30", gindex: (12*4 + 3) * 8, offset: 30},
		{path: "vector_uint64", gindex: 13},
		// 40 64-bit ints occupy 320 bytes (10 chunks), nextPowerOfTwo(10) = 16
		{path: "vector_uint64/5", gindex: 13*16 + (5 / 4), offset: 8},

		// error cases
		{path: "nested/__len__", error: "__len__ is only valid"},
	}
	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc.path, "/", "."), func(t *testing.T) {
			objectPath := merkle.ObjectPath[uint64, [32]byte](tc.path)
			typ, gindex, offset, err := objectPath.GetGeneralizedIndex(root)

			if tc.error != "" {
				require.ErrorContains(t, err, tc.error)
				return
			}

			require.NoError(t, err)
			require.NotNil(
				t, typ, "Type should not be nil")
			require.Equal(
				t, tc.gindex, gindex, "Unexpected generalized index",
			)
			require.Equal(t, tc.offset, offset, "Unexpected offset")
		})
	}
}

func TestPaths(t *testing.T) {
	bs := schema.DefineContainer(
		schema.NewField("GenesisValidatorsRoot", schema.B32()),
		schema.NewField("Slot", schema.U64()),
		schema.NewField("Fork", schema.DefineContainer(
			schema.NewField("PreviousVersion", schema.B4()),
			schema.NewField("CurrentVersion", schema.B4()),
			schema.NewField("Epoch", schema.U64()),
		)),
		schema.NewField("LatestBlockHeader", schema.DefineContainer(
			schema.NewField("Slot", schema.U64()),
			schema.NewField("ProposerIndex", schema.U64()),
			schema.NewField("ParentBlockRoot", schema.B32()),
			schema.NewField("StateRoot", schema.B32()),
			schema.NewField("BodyRoot", schema.B32()),
		)),
		schema.NewField("BlockRoots", schema.DefineList(schema.B32(), 8192)),
		schema.NewField("StateRoots", schema.DefineList(schema.B32(), 8192)),
		schema.NewField("Eth1Data", schema.DefineContainer(
			schema.NewField("DepositRoot", schema.B32()),
			schema.NewField("DepositCount", schema.U64()),
			schema.NewField("BlockHash", schema.B32()),
		)),
		schema.NewField("Eth1DepositIndex", schema.U64()),
		schema.NewField("LatestExecutionPayloadHeader", schema.DefineContainer(
			schema.NewField("ParentHash", schema.B32()),
			schema.NewField("FeeRecipient", schema.B20()),
			schema.NewField("StateRoot", schema.B32()),
			schema.NewField("ReceiptsRoot", schema.B32()),
			schema.NewField("LogsBloom", schema.B256()),
			schema.NewField("Random", schema.U64()),
			schema.NewField("Number", schema.U64()),
			schema.NewField("GasLimit", schema.U64()),
			schema.NewField("GasUsed", schema.U64()),
			schema.NewField("Timestamp", schema.U64()),
			schema.NewField("ExtraData", schema.DefineByteList(32)),
			schema.NewField("BaseFeePerGas", schema.B32()),
			schema.NewField("BlockHash", schema.B32()),
			schema.NewField("TransactionsRoot", schema.B32()),
			schema.NewField("WithdrawalsRoot", schema.B32()),
			schema.NewField("BlobGasUsed", schema.U64()),
			schema.NewField("ExcessBlobGas", schema.U64()),
		)),
		schema.NewField("Validators", schema.DefineList(schema.DefineContainer(
			schema.NewField("Pubkey", schema.B48()),
			schema.NewField("WithdrawalCredentials", schema.B32()),
			schema.NewField("EffectiveBalance", schema.U64()),
			schema.NewField("Slashed", schema.Bool()),
			schema.NewField("ActivationEligibilityEpoch", schema.U64()),
			schema.NewField("ActivationEpoch", schema.U64()),
			schema.NewField("ExitEpoch", schema.U64()),
			schema.NewField("WithdrawableEpoch", schema.U64()),
		), 1099511627776)),
		schema.NewField("Balances", schema.DefineList(schema.U64(), 1099511627776)),
		schema.NewField("RandaoMixes", schema.DefineList(schema.B32(), 65536)),
		schema.NewField("NextWithdrawalIndex", schema.U64()),
		schema.NewField("NextWithdrawalValidatorIndex", schema.U64()),
		schema.NewField("Slashings", schema.DefineList(schema.U64(), 1099511627776)),
		schema.NewField("TotalSlashing", schema.U64()),
	)

	blockHeader := schema.DefineContainer(
		schema.NewField("Slot", schema.U64()),
		schema.NewField("ProposerIndex", schema.U64()),
		schema.NewField("ParentBlockRoot", schema.B32()),
		schema.NewField("State", bs),
		schema.NewField("BodyRoot", schema.B32()),
	)

	val0PubKeyPath := merkle.ObjectPath[merkle.GeneralizedIndex, [32]byte](
		"State/Validators/2",
	)
	_, gindex, offset, err := val0PubKeyPath.GetGeneralizedIndex(blockHeader)
	require.NoError(t, err)
	fmt.Println("gIndex", gindex)
	fmt.Println("offset", offset)

	panic("see logs")
}
