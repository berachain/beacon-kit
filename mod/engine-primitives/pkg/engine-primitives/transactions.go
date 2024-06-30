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
	"unsafe"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkleizer"
)

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
//
// NOTE: We made a mistake on bArtio here. This shoudl've been ssz.ListBasic
// for the actual transactions type.
// TODO: We will deprecate this type in the future.
type BartioTransactions = ssz.ListComposite[ssz.VectorBasic[ssz.Byte]]

// BartioTransactionsFromBytes creates a Transactions object from a byte slice.
func BartioTransactionsFromBytes(data [][]byte) *BartioTransactions {
	return ssz.ListCompositeFromElements(
		// TODO: Move this value to chain spec.
		constants.MaxTxsPerPayload,
		*(*[]ssz.VectorBasic[ssz.Byte])(unsafe.Pointer(&data))...)
}

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
type Transactions = ssz.ListComposite[*ssz.ListBasic[ssz.Byte]]

// TransactionsFromBytes creates a Transactions object from a byte slice.
func TransactionsFromBytes(data [][]byte) *Transactions {
	d := *(*[][]ssz.Byte)(unsafe.Pointer(&data))
	txs := make([]*ssz.ListBasic[ssz.Byte], 0)
	for _, i := range d {
		txs = append(
			txs,
			ssz.ListBasicFromElements(constants.MaxBytesPerTransaction, i...),
		)
	}
	return ssz.ListCompositeFromElements(
		// TODO: Move this value to chain spec.
		constants.MaxTxsPerPayload, txs...,
	)
}

// TODO: make the ChainSpec a generic on this type.
type TxsMerkleizer merkleizer.
	Merkleizer[[32]byte, ssz.VectorBasic[ssz.Byte]]
