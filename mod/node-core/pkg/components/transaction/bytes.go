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

package transaction

import (
	"cosmossdk.io/core/transaction"
	"github.com/cosmos/gogoproto/proto"
)

type BytesTx []byte

var _ transaction.Tx = (*BytesTx)(nil)

func NewTxFromBytes[T transaction.Tx](bz []byte) T {
	return any(BytesTx(bz)).(T)
}

func (tx BytesTx) Bytes() []byte {
	return tx
}

func (tx BytesTx) Hash() [32]byte {
	return [32]byte{}
}

func (tx BytesTx) GetGasLimit() (uint64, error) { return 0, nil }

func (tx BytesTx) GetMessages() ([]proto.Message, error) { return nil, nil }

func (tx BytesTx) GetSenders() ([][]byte, error) { return nil, nil }

var _ transaction.Codec[transaction.Tx] = (*BytesTxCodec[transaction.Tx])(nil)

type BytesTxCodec[T transaction.Tx] struct{}

func (cdc *BytesTxCodec[T]) Decode(bz []byte) (T, error) {
	return NewTxFromBytes[T](bz), nil
}

// DecodeJSON decodes the tx JSON bytes into a DecodedTx
func (cdc *BytesTxCodec[T]) DecodeJSON(bz []byte) (T, error) {
	return NewTxFromBytes[T](bz), nil
}
