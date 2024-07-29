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

package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/karalabe/ssz"
)

type BartioTx []byte

func (tx *BartioTx) SizeSSZ() uint32 {
	return ssz.SizeDynamicBytes(*tx)
}

func (tx *BartioTx) DefineSSZ(codec *ssz.Codec) {
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineDynamicBytesOffset(
			codec,
			(*[]byte)(tx),
			constants.MaxTxsPerPayload,
		)
	})
}

type Roots []common.Root

func (roots Roots) SizeSSZ() uint32 {
	return ssz.SizeSliceOfStaticBytes(roots)
}

func (roots Roots) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticBytesContent(
			codec,
			(*[]common.Root)(&roots),
			constants.MaxTxsPerPayload,
		)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticBytesContent(
			codec,
			(*[]common.Root)(&roots),
			constants.MaxTxsPerPayload,
		)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticBytesOffset(
			codec, (*[]common.Root)(&roots), constants.MaxTxsPerPayload,
		)
	})
}

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
//
// TODO: Remove and deprecate this type once migrated to ProperTransactions.
type BartioTransactions [][]byte

// HashTreeRoot returns the hash tree root of the Transactions list.
//
// NOTE: Uses a new merkleizer for each call.
func (txs BartioTransactions) HashTreeRoot() common.Root {
	return txs.HashTreeRootWith(
		merkle.NewMerkleizer[[32]byte, common.Root](),
	)
}

func (txs BartioTransactions) HashTreeRoot2() common.Root {
	roots := make(Roots, len(txs))
	merkleizer := merkle.NewMerkleizer[[32]byte, common.Root]()
	for i, tx := range txs {
		var err error
		roots[i], err = merkleizer.MerkleizeByteSlice(tx)
		if err != nil {
			panic(err)
		}
	}
	return ssz.HashConcurrent(roots)
}

// HashTreeRootWith returns the hash tree root of the Transactions list
// using the given merkle.
func (txs BartioTransactions) HashTreeRootWith(
	merkleizer *merkle.Merkleizer[[32]byte, common.Root],
) common.Root {
	var (
		err   error
		root  common.Root
		roots = make([]common.Root, len(txs))
	)

	for i, tx := range txs {
		roots[i], err = merkleizer.MerkleizeByteSlice(tx)
		if err != nil {
			panic(err)
		}
	}
	// fmt.Println("roots1", roots)
	root, err = merkleizer.MerkleizeListComposite(
		roots,
		constants.MaxTxsPerPayload,
	)
	if err != nil {
		panic(err)
	}
	return root
}

// ProperTransactions is a type alias for [][]byte, which is how transactions
// are
// received in the execution payload.
type Transactions [][]byte

// SizeSSZ returns the SSZ encoded size in bytes for the Transactions.
func (txs Transactions) SizeSSZ(fixed bool) uint32 {
	if fixed || txs == nil {
		return 0
	}
	return ssz.SizeSliceOfDynamicBytes(txs)
}

// DefineSSZ defines the SSZ encoding for the Transactions object.
// TODO: This can accidentally decouple from the definition in
// ExecutionPayload and we should be cognizant of that and/or
// make a PR to allow for them to be defined in one place.
func (txs Transactions) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfDynamicBytesContent(
			codec,
			(*[][]byte)(&txs),
			constants.MaxTxsPerPayload,
			constants.MaxBytesPerTx,
		)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfDynamicBytesContent(
			codec,
			(*[][]byte)(&txs),
			constants.MaxTxsPerPayload,
			constants.MaxBytesPerTx,
		)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfDynamicBytesOffset(
			codec,
			(*[][]byte)(&txs),
			constants.MaxTxsPerPayload,
			constants.MaxBytesPerTx,
		)
	})
}

// HashTreeRoot returns the hash tree root of the Transactions object.
func (txs Transactions) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(txs)
}
