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

package appmanager

import (
	"time"

	appmanager "cosmossdk.io/core/app"
	"cosmossdk.io/core/transaction"
)

type abciRequest struct {
	height int64
	time   time.Time
	txs    [][]byte
}

func blockToABCIRequest[T transaction.Tx](
	block *appmanager.BlockRequest[T],
) *abciRequest {
	txs := make([][]byte, len(block.Txs))
	for i, tx := range block.Txs {
		txs[i] = tx.Bytes()
	}
	return &abciRequest{
		height: int64(block.Height),
		time:   block.Time,
		txs:    txs,
	}
}

func (req *abciRequest) GetHeight() int64 {
	return req.height
}

func (req *abciRequest) GetTime() time.Time {
	return req.time
}

func (req *abciRequest) GetTxs() [][]byte {
	return req.txs
}

func (req *abciRequest) SetTxs(txs [][]byte) {
	req.txs = txs
}
