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
	"slices"
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

type BbbDeneb struct {
	types.BeaconBlockBody
}

type BodyTypes interface {
	types.Eth1Data | types.Deposit | engineprimitives.Withdrawal |
	[]byte | types.ExecutionPayload | eip4844.KZGCommitment
}

func orNil[T interface{ *E | []*E | []E }, E BodyTypes](
	what T, r *rand.Rand,
) T {
	if r.Intn(2) == 0 {
		return what
	}
	return nil
}

func (b *BbbDeneb) Generate(r *rand.Rand, _ int) reflect.Value {
	roll := func(n int) int {
		k := r.Intn(n)
		if k%2 == 0 {
			return 0
		}
		return k
	}
	rbytes := func(n int) []byte {
		rbs := make([]byte, n)
		r.Read(rbs)
		return rbs
	}

	sizer := &ssz.Sizer{}
	b = &BbbDeneb{}
	b.RandaoReveal = crypto.BLSSignature(rbytes(96))
	b.Eth1Data = orNil[*types.Eth1Data, types.Eth1Data](&types.Eth1Data{
		DepositRoot:  common.Root(rbytes(32)),
		DepositCount: math.U64(r.Uint64()),
		BlockHash:    common.ExecutionHash(rbytes(32)),
	}, r)
	b.Graffiti = [32]byte(rbytes(32))

	k := roll(16)
	var depositsLen = uint32(k)
	if k != 0 {
		depositsLen = uint32(concurrencyThreshold)/(
			&types.Deposit{}).SizeSSZ(sizer) + 1
	}
	b.Deposits = make([]*types.Deposit, depositsLen)
	for i := range depositsLen {
		b.Deposits[i] = &types.Deposit{
			Pubkey:      crypto.BLSPubkey(rbytes(48)),
			Credentials: types.WithdrawalCredentials(rbytes(32)),
			Amount:      math.Gwei(r.Uint64()),
			Signature:   crypto.BLSSignature(rbytes(96)),
			Index:       r.Uint64(),
		}
	}
	b.Deposits = orNil[[]*types.Deposit, types.Deposit](b.Deposits, r)

	k = roll(10) // MaxTxsPerPayload 1048576 too big
	txs := make([][]byte, k)
	for i := range k {
		txs[i] = rbytes(1024) // MaxBytesPerTx 1073741824 too big
	}
	txs = orNil[[][]byte, []byte](txs, r)

	k = roll(16)
	withdrawals := make([]*engineprimitives.Withdrawal, k)
	for i := range k {
		withdrawals[i] = &engineprimitives.Withdrawal{
			Index:     math.U64(r.Uint64()),
			Validator: math.U64(r.Uint64()),
			Address:   common.ExecutionAddress(rbytes(20)),
			Amount:    math.U64(r.Uint64()),
		}
	}
	withdrawals = orNil[
		[]*engineprimitives.Withdrawal, engineprimitives.Withdrawal,
	](withdrawals, r)

	b.ExecutionPayload = orNil[
		*types.ExecutionPayload, types.ExecutionPayload,
	](&types.ExecutionPayload{
		ParentHash:    common.ExecutionHash(rbytes(32)),
		FeeRecipient:  common.ExecutionAddress(rbytes(20)),
		StateRoot:     bytes.B32(rbytes(32)),
		ReceiptsRoot:  bytes.B32(rbytes(32)),
		LogsBloom:     bytes.B256(rbytes(256)),
		Random:        common.Bytes32(rbytes(32)),
		Number:        math.U64(r.Uint64()),
		GasLimit:      math.U64(r.Uint64()),
		GasUsed:       math.U64(r.Uint64()),
		Timestamp:     math.U64(r.Uint64()),
		ExtraData:     bytes.Bytes(rbytes(32)),
		BaseFeePerGas: math.NewU256(r.Uint64()),
		BlockHash:     common.ExecutionHash(rbytes(32)),
		Transactions:  txs,
		Withdrawals:   withdrawals,
		BlobGasUsed:   math.U64(r.Uint64()),
		ExcessBlobGas: math.U64(r.Uint64()),
	}, r)

	k = roll(4096)
	b.BlobKzgCommitments = make([]eip4844.KZGCommitment, k)
	for i := range k {
		b.BlobKzgCommitments[i] = eip4844.KZGCommitment(rbytes(48))
	}
	b.BlobKzgCommitments = orNil[
		[]eip4844.KZGCommitment, eip4844.KZGCommitment,
	](b.BlobKzgCommitments, r)

	return reflect.ValueOf(b)
}

func pprint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func anyNil(a *BbbDeneb) bool {
	var r bool
	r = a == nil
	r = r || a.Eth1Data == nil
	r = r || a.Deposits == nil
	r = r || a.ExecutionPayload == nil
	r = r || a.ExecutionPayload.Transactions == nil
	r = r || a.ExecutionPayload.Withdrawals == nil
	r = r || a.BlobKzgCommitments == nil
	return r
}

func compare(a, b *BbbDeneb) bool {
	var r bool
	r = a.RandaoReveal == b.RandaoReveal
	r = r && (a.Eth1Data != nil && *a.Eth1Data == *b.Eth1Data) ||
		(a.Eth1Data == nil && *b.Eth1Data == types.Eth1Data{})
	r = r && a.Graffiti == b.Graffiti
	if a.Deposits != nil {
		for i, depA := range a.Deposits {
			r = r && *depA == *b.Deposits[i]
		}
	} else {
		r = r && len(b.Deposits) == 0
	}
	if a.ExecutionPayload != nil {
		for i, txA := range a.ExecutionPayload.Transactions {
			r = r && slices.Equal(
				txA, b.ExecutionPayload.Transactions[i])
		}
		for i, wA := range a.ExecutionPayload.Withdrawals {
			r = r && *wA == *b.ExecutionPayload.Withdrawals[i]
		}
	} else {
		r = r && len(b.ExecutionPayload.Transactions) == 0
		r = r && len(b.ExecutionPayload.Withdrawals) == 0
	}
	if a.BlobKzgCommitments != nil {
		for i, cA := range a.BlobKzgCommitments {
			r = r && cA == b.BlobKzgCommitments[i]
		}
	} else {
		r = r && len(b.BlobKzgCommitments) == 0
	}
	return r
}

func TestSSZRoundTripBeaconBodyDeneb(t *testing.T) {
	t.Parallel()
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
		if anyNil(destBody) {
			t.Log("Deserialize: nil found in deserialized body")
			t.Log(pprint(body))
			t.Log(pprint(destBody))
			return false
		}

		if !compare(body, destBody) {
			t.Log("Deserialized body different than former body after serialization")
			t.Log(pprint(body))
			t.Log(pprint(destBody))
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

		if len(body.GetDeposits()) > 0 {
			htrSeq := ssz.HashSequential(body)
			htrC := ssz.HashConcurrent(body)
			if !reflect.DeepEqual(htrSeq, htrC) {
				t.Log("Sequential hash != Concurrent hash")
				t.Log(pprint(body))
				t.Log(htrSeq)
				t.Log(htrC)
				return false
			}
		}

		return true
	}

	if err := quick.Check(f, &Conf); err != nil {
		t.Error(err)
	}
}
