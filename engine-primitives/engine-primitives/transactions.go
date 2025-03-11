// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/karalabe/ssz"
)

// ProperTransactions is a type alias for [][]byte, which is how
// transactions are received in the execution payload.
type Transactions [][]byte

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Transactions.
func (txs Transactions) SizeSSZ(siz *ssz.Sizer, _ bool) uint32 {
	return ssz.SizeSliceOfDynamicBytes(siz, txs)
}

// DefineSSZOffset defines the SSZ offset for the Transactions object.
func (txs Transactions) DefineSSZOffset(c *ssz.Codec) {
	ssz.DefineSliceOfDynamicBytesOffset(
		c,
		(*[][]byte)(&txs),
		constants.MaxTxsPerPayload,
		constants.MaxBytesPerTx,
	)
}

// DefineSSZContent defines the SSZ content for the Transactions object.
func (txs Transactions) DefineSSZContent(c *ssz.Codec) {
	ssz.DefineSliceOfDynamicBytesContent(
		c,
		(*[][]byte)(&txs),
		constants.MaxTxsPerPayload,
		constants.MaxBytesPerTx,
	)
}

// DefineSSZ defines the SSZ (en/de)coding and hashing for the Transactions object.
func (txs Transactions) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		txs.DefineSSZContent(codec)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		txs.DefineSSZContent(codec)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		txs.DefineSSZOffset(codec)
	})
}

// HashTreeRoot returns the hash tree root of the Transactions object.
func (txs Transactions) HashTreeRoot() common.Root {
	return ssz.HashConcurrent(txs)
}
