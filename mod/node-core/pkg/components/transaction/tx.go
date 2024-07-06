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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/cosmos/gogoproto/proto"
)

var _ transaction.Tx = (*SSZTx)(nil)

type SSZTx struct {
	constraints.SSZMarshallable
}

func NewTxFromSSZ[T transaction.Tx](bz []byte) T {
	tx := &SSZTx{}
	if err := tx.UnmarshalSSZ(bz); err != nil {
		panic(err)
	}

	return any(tx).(T)
}

// func (tx *SSZTx) Unwrap() *SSZTx {
// 	return tx
// }

func (tx *SSZTx) Bytes() []byte {
	bz, err := tx.MarshalSSZ()
	if err != nil {
		panic(err)
	}
	return bz
}

func (tx *SSZTx) Hash() [32]byte {
	bz, err := tx.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	return bz
}

func (tx *SSZTx) GetGasLimit() (uint64, error) { return 0, nil }

func (tx *SSZTx) GetMessages() ([]proto.Message, error) { return nil, nil }

func (tx *SSZTx) GetSenders() ([][]byte, error) { return nil, nil }
