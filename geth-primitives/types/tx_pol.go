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

package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// PoLTx represents an BRIP-0004 transaction. No gas is consumed for execution.
type PoLTx struct {
	ChainID  *big.Int
	From     common.Address // system address
	To       common.Address // address of the PoL Distributor contract
	Nonce    uint64         // block number distributing for
	GasLimit uint64         // artificial gas limit for the PoL tx, not consumed against the block gas limit
	GasPrice *big.Int       // gas price is set to the baseFee to make the tx valid for EIP-1559 rules
	Data     []byte         // encodes the pubkey distributing for
}

func (*PoLTx) txType() byte { return PoLTxType }

func (tx *PoLTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *PoLTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
