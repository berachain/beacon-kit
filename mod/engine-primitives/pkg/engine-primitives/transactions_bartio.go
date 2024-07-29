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
	"github.com/karalabe/ssz"
)

// BartioTransactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload on the bArtio testnet. This is due to a
// mistake made during the initial implementation of BeaconKit. This type will
// be deprecated off of
// eventually.
type BartioTransactions [][]byte

// HashTreeRoot returns the hash tree root of the Transactions list.
//
// NOTE: Uses a new merkleizer for each call.
func (txs BartioTransactions) HashTreeRoot() common.Root {
	roots := make(Roots, len(txs))
	for i, tx := range txs {
		roots[i] = BartioTx(tx).HashTreeRoot()
	}
	return ssz.HashConcurrent(roots)
}

// BartioTx represents a single transaction in the Bartio format.
type BartioTx []byte

// SizeSSZ returns the SSZ sssize of the BartioTx.
func (tx BartioTx) SizeSSZ() uint32 {
	return ssz.SizeDynamicBytes(tx)
}

// DefineSSZ implements the SSZ encoding for BartioTx.
func (tx BartioTx) DefineSSZ(codec *ssz.Codec) {
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineStaticBytes(
			codec,
			(*[]byte)(&tx),
		)
	})
}

// HashTreeRoot returns the Merkle root hash of the BartioTx.
func (tx BartioTx) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(tx)
}

// Roots is a list of common.Roots.
type Roots []common.Root

// SizeSSZ returns the SSZ size of the Roots object.
func (roots Roots) SizeSSZ() uint32 {
	return ssz.SizeSliceOfStaticBytes(roots)
}

// DefineSSZ defines the SSZ encoding for the Roots object.
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
