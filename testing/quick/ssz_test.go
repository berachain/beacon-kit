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

package quick_test

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

func roll(n int, r *rand.Rand) int {
	k := r.Intn(n)
	if k%2 == 0 {
		return 0
	}
	return k
}

func rbytes(n int, r *rand.Rand) []byte {
	rbs := make([]byte, n)
	k := r.Intn(n)
	if k%2 == 0 {
		r.Read(rbs)
	}
	return rbs
}

type BbbDeneb struct {
	types.BeaconBlockBody
}

func (b *BbbDeneb) Generate(r *rand.Rand, size int) reflect.Value {
	b = &BbbDeneb{}
	b.RandaoReveal = crypto.BLSSignature(rbytes(96, r))
	b.Eth1Data = &types.Eth1Data{
		DepositRoot:  common.Root(rbytes(32, r)),
		DepositCount: math.U64(r.Uint64()),
		BlockHash:    common.ExecutionHash(rbytes(32, r)),
	}
	b.Graffiti = [32]byte(rbytes(32, r))
	k := roll(16, r)
	b.Deposits = make([]*types.Deposit, k)
	for i := 0; i < k; i++ {
		var proof [33][32]byte
		for j := 0; j < 33; j++ {
			proof[j] = [32]byte(rbytes(32, r))
		}
		b.Deposits[i] = &types.Deposit{
			Pubkey:      crypto.BLSPubkey(rbytes(48, r)),
			Credentials: types.WithdrawalCredentials(rbytes(32, r)),
			Amount:      math.Gwei(r.Uint64()),
			Signature:   crypto.BLSSignature(rbytes(96, r)),
			Index:       r.Uint64(),
		}
	}
	k = roll(10, r) // MaxTxsPerPayload 1048576 too big
	txs := make([][]byte, k)
	for i := 0; i < k; i++ {
		txs[i] = rbytes(1024, r) // MaxBytesPerTx 1073741824 too big
	}
	k = roll(16, r)
	withdrawals := make([]*engineprimitives.Withdrawal, k)
	for i := 0; i < k; i++ {
		withdrawals[i] = &engineprimitives.Withdrawal{
			math.U64(r.Uint64()),
			math.U64(r.Uint64()),
			common.ExecutionAddress(rbytes(20, r)),
			math.U64(r.Uint64()),
		}
	}
	b.ExecutionPayload = &types.ExecutionPayload{
		ParentHash:    common.ExecutionHash(rbytes(32, r)),
		FeeRecipient:  common.ExecutionAddress(rbytes(20, r)),
		StateRoot:     bytes.B32(rbytes(32, r)),
		ReceiptsRoot:  bytes.B32(rbytes(32, r)),
		LogsBloom:     bytes.B256(rbytes(256, r)),
		Random:        common.Bytes32(rbytes(32, r)),
		Number:        math.U64(r.Uint64()),
		GasLimit:      math.U64(r.Uint64()),
		GasUsed:       math.U64(r.Uint64()),
		Timestamp:     math.U64(r.Uint64()),
		ExtraData:     bytes.Bytes(rbytes(32, r)),
		BaseFeePerGas: math.NewU256(r.Uint64()),
		BlockHash:     common.ExecutionHash(rbytes(32, r)),
		Transactions:  txs,
		Withdrawals:   withdrawals,
		BlobGasUsed:   math.U64(r.Uint64()),
		ExcessBlobGas: math.U64(r.Uint64()),
	}
	k = roll(16, r) // 4096 in Deneb
	b.BlobKzgCommitments = make([]eip4844.KZGCommitment, k)
	for i := 0; i < k; i++ {
		b.BlobKzgCommitments[i] = eip4844.KZGCommitment(rbytes(48, r))
	}

	return reflect.ValueOf(b)
}

func pprint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func TestSSZRoundTripBeaconBodyDeneb(t *testing.T) {
	f := func(body *BbbDeneb, n uint) bool {
		bz, err := body.MarshalSSZ()
		if err != nil {
			t.Log("Serialize: could not serialize body --", err)
			return false
		}
		destBody := &BbbDeneb{}
		if err = destBody.UnmarshalSSZ(bz); err != nil {
			t.Log("Deserialize: could not deserialize back the serialized body --", err)
			return false
		}

		if !reflect.DeepEqual(body, destBody) {
			t.Log("Deserialize: deserialized body different than former body after serialization")
			t.Log(pprint(body))
			t.Log(pprint(destBody))
			return false
		}

		htr := body.HashTreeRoot()
		destHtr := destBody.HashTreeRoot()
		if !reflect.DeepEqual(htr, destHtr) {
			t.Log("Hash tree root differs after serialization-deserialization round trip")
		}

		destBz, err := destBody.MarshalSSZ()
		if err != nil {
			t.Log("Serialize: could not serialize back the body after deserialization --", err)
			return false
		}

		if !reflect.DeepEqual(bz, destBz) {
			t.Log("Serialize: serialized body different after a serialization-deserialization-serialization trip")
			return false
		}

		return true
	}

	if err := quick.Check(f, &Conf); err != nil {
		t.Error(err)
	}
}

var concurrencyThreshold uint64 = 65536

type Container struct {
	Deposits []*types.Deposit
}

func (c *Container) SizeSSZ() uint32 {
	return ssz.SizeSliceOfStaticObjects(c.Deposits)
}

func (c *Container) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsOffset(codec, &c.Deposits, concurrencyThreshold)
	ssz.DefineSliceOfStaticObjectsContent(codec, &c.Deposits, concurrencyThreshold)
}

func (c *Container) Generate(r *rand.Rand, size int) reflect.Value {
	deposits := make([]*types.Deposit,
			 uint32(concurrencyThreshold)/(&types.Deposit{}).SizeSSZ()+1)
	for i := 0; i < len(deposits); i++ {
		var proof [33][32]byte
		for j := 0; j < 33; j++ {
			proof[j] = [32]byte(rbytes(32, r))
		}
		deposits[i] = &types.Deposit{
			Pubkey:      crypto.BLSPubkey(rbytes(48, r)),
			Credentials: types.WithdrawalCredentials(rbytes(32, r)),
			Amount:      math.Gwei(r.Uint64()),
			Signature:   crypto.BLSSignature(rbytes(96, r)),
			Index:       r.Uint64(),
		}
	}
	c = &Container{Deposits: deposits}

	return reflect.ValueOf(c)
}

func TestHashConcurrent(t *testing.T) {
	f := func(c *Container) bool {
		htrSeq := ssz.HashSequential(c)
		htrC := ssz.HashConcurrent(c)
		if !reflect.DeepEqual(htrSeq, htrC) {
			t.Log("Sequential hash != Concurrent hash")
			t.Log(pprint(c))
			return false
		}
		return true
	}

	if err := quick.Check(f, &Conf); err != nil {
		t.Error(err)
	}
}
