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

package compare_test

import (
	"slices"
	"testing"
	"testing/quick"
	"unsafe"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz"
	zcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	zdeneb "github.com/protolambda/zrnt/eth2/beacon/deneb"
	zspec "github.com/protolambda/zrnt/eth2/configs"
	ztree "github.com/protolambda/ztyp/tree"
	zview "github.com/protolambda/ztyp/view"
)

var c = quick.Config{MaxCount: 1000000}
var hFn = ztree.GetHashFn()
var spec = zspec.Mainnet

func TestListHashTreeRootZtyp(t *testing.T) {
	f := func(s []byte, limit uint64) bool {
		a := ssz.ByteListFromBytes(s, limit)

		root1, err := a.HashTreeRoot()
		root2 := hFn.ByteListHTR(s, limit)
		if err != nil {
			t.Log("Failed to calculate HashTreeRoot on input: ", err)
			return false
		}
		return root1 == [32]byte(root2)
	}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}

func TestVectorHashTreeRootZTyp(t *testing.T) {
	f := func(s []byte) bool {
		a := ssz.ByteVectorFromBytes(s)

		root1, err := a.HashTreeRoot()
		root2 := hFn.ByteVectorHTR(s)
		if err != nil {
			t.Log("Failed to calculate HashTreeRoot on input: ", err)
			return false
		}
		return root1 == [32]byte(root2)
	}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}

func TestExecutableDataDenebHashTreeRootZrnt(t *testing.T) {
	f := func(payload *types.ExecutableDataDeneb, logsBloom [256]byte) bool {
		// skip in these cases lest we trigger a nil-pointer dereference in fastssz's HashTreeRootWith
		if payload == nil || payload.Withdrawals == nil || slices.Contains(payload.Withdrawals, nil) ||
			payload.Transactions == nil || slices.ContainsFunc(payload.Transactions, func(e []byte) bool {
			return e == nil
		}) {
			return true
		}

		payload.LogsBloom = logsBloom[:]
		typeRoot, err := payload.HashTreeRoot()
		if err != nil {
			t.Log("Failed to calculate HashTreeRoot on payload:", err)
			return false
		}

		baseFeePerGas := zview.Uint256View{}
		baseFeePerGas.SetBytes32(payload.BaseFeePerGas.Unwrap())
		zpayload := zdeneb.ExecutionPayload{
			ParentHash:    ztree.Root(payload.ParentHash),
			FeeRecipient:  zcommon.Eth1Address(payload.FeeRecipient),
			StateRoot:     ztree.Root(payload.StateRoot),
			ReceiptsRoot:  ztree.Root(payload.ReceiptsRoot),
			LogsBloom:     zcommon.LogsBloom(payload.LogsBloom),
			PrevRandao:    ztree.Root(payload.Random),
			BlockNumber:   zview.Uint64View(payload.Number),
			GasLimit:      zview.Uint64View(payload.GasLimit),
			GasUsed:       zview.Uint64View(payload.GasUsed),
			Timestamp:     zcommon.Timestamp(payload.Timestamp),
			ExtraData:     payload.ExtraData,
			BaseFeePerGas: baseFeePerGas,
			BlockHash:     ztree.Root(payload.BlockHash),
			Transactions:  *(*zcommon.PayloadTransactions)(unsafe.Pointer(&payload.Transactions)),
			Withdrawals:   *(*zcommon.Withdrawals)(unsafe.Pointer(&payload.Withdrawals)),
			BlobGasUsed:   zview.Uint64View(payload.BlobGasUsed.Unwrap()),
			ExcessBlobGas: zview.Uint64View(payload.ExcessBlobGas.Unwrap()),
		}
		zRoot := zpayload.HashTreeRoot(spec, hFn)

		container := ssz.ContainerFromElements(
			ssz.ByteVectorFromBytes(payload.ParentHash[:]),
			ssz.ByteVectorFromBytes(payload.FeeRecipient[:]),
			ssz.ByteVectorFromBytes(payload.StateRoot[:]),
			ssz.ByteVectorFromBytes(payload.ReceiptsRoot[:]),
			ssz.ByteVectorFromBytes(payload.LogsBloom),
			ssz.ByteVectorFromBytes(payload.Random[:]),
			payload.Number,
			payload.GasLimit,
			payload.GasUsed,
			payload.Timestamp,
			ssz.ByteListFromBytes(payload.ExtraData, 32),
			payload.BaseFeePerGas.UnwrapU256(),
			ssz.ByteVectorFromBytes(payload.BlockHash[:]),
			engineprimitives.ProperTransactionsFromBytes(payload.Transactions),
			ssz.ListFromElements(16, payload.Withdrawals...),
			payload.BlobGasUsed,
			payload.ExcessBlobGas,
		)
		containerRoot, err := container.HashTreeRoot()
		if err != nil {
			return false
		}
		return typeRoot == containerRoot && typeRoot == zRoot
	}
	if err := quick.Check(f, &c); err != nil {
		t.Error(err)
	}
}
