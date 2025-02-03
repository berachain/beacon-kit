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

package randomize_test

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
)

var concurrencyThreshold uint64 = 65536

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

func (b *BbbDeneb) Generate(r *rand.Rand, _ int) reflect.Value {
	b = &BbbDeneb{}
	b.RandaoReveal = crypto.BLSSignature(rbytes(96, r))
	b.Eth1Data = &types.Eth1Data{
		DepositRoot:  common.Root(rbytes(32, r)),
		DepositCount: math.U64(r.Uint64()),
		BlockHash:    common.ExecutionHash(rbytes(32, r)),
	}
	b.Graffiti = [32]byte(rbytes(32, r))
	k := roll(16, r)
	sizer := &ssz.Sizer{}
	b.Deposits = make([]*types.Deposit,
		uint32(concurrencyThreshold)/(&types.Deposit{}).SizeSSZ(sizer)+1)
	for i := 0; i < len(b.Deposits); i++ {
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
	for i := range k {
		txs[i] = rbytes(1024, r) // MaxBytesPerTx 1073741824 too big
	}
	k = roll(16, r)
	withdrawals := make([]*engineprimitives.Withdrawal, k)
	for i := range k {
		withdrawals[i] = &engineprimitives.Withdrawal{
			Index:     math.U64(r.Uint64()),
			Validator: math.U64(r.Uint64()),
			Address:   common.ExecutionAddress(rbytes(20, r)),
			Amount:    math.U64(r.Uint64()),
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
	k = roll(4096, r)
	b.BlobKzgCommitments = make([]eip4844.KZGCommitment, k)
	for i := range k {
		b.BlobKzgCommitments[i] = eip4844.KZGCommitment(rbytes(48, r))
	}

	return reflect.ValueOf(b)
}

func pprint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func TestSSZRoundTripBeaconBodyDeneb(t *testing.T) {
	f := func(body *BbbDeneb) bool {
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

		/* TODO investigate why this fails
		if !reflect.DeepEqual(body, destBody) {
			t.Log("Deserialized body different than former body after serialization")
			t.Log(pprint(body))
			t.Log(pprint(destBody))
			return false
		}
		*/

		if destBody.ExecutionPayload.GetWithdrawals() == nil {
			t.Log("Withdrawals is nil after deserialization")
			return false
		}

		htr := body.HashTreeRoot()
		destHtr := destBody.HashTreeRoot()
		if !reflect.DeepEqual(htr, destHtr) {
			t.Log("HTR differs after serialization-deserialization round trip")
			t.Log(htr)
			t.Log(destHtr)
		}

		destBz, err := destBody.MarshalSSZ()
		if err != nil {
			t.Log("Could not serialize back the body after deserialization --", err)
			return false
		}

		if !reflect.DeepEqual(bz, destBz) {
			t.Log("Serialized body different after a",
				"serialization-deserialization-serialization trip")
			t.Log(pprint(body))
			t.Log(pprint(destBody))
			t.Log(bz)
			t.Log(destBz)
			return false
		}

		htrSeq := ssz.HashSequential(body)
		htrC := ssz.HashConcurrent(body)
		if !reflect.DeepEqual(htrSeq, htrC) {
			t.Log("Sequential hash != Concurrent hash")
			t.Log(pprint(body))
			t.Log(htrSeq)
			t.Log(htrC)
			return false
		}

		return true
	}

	if err := quick.Check(f, &Conf); err != nil {
		t.Error(err)
	}
}
