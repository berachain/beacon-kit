// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package consensusprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives"
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./deposit.go -objs Deposit -include ./withdrawal_credentials.go,../primitives,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deposit.ssz.go
type Deposit struct {
	// Public key of the validator specified in the deposit.
	Pubkey primitives.BLSPubkey `json:"pubkey" ssz-max:"48"`

	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials" ssz-size:"32"`

	// Deposit amount in gwei.
	Amount primitives.Gwei `json:"amount"`

	// Signature of the deposit data.
	Signature primitives.BLSSignature `json:"signature" ssz-max:"96"`

	// Index of the deposit in the deposit contract.
	Index uint64 `json:"index"`
}

// NewDeposit creates a new Deposit instance.
func NewDeposit(
	pubkey primitives.BLSPubkey,
	credentials WithdrawalCredentials,
	amount primitives.Gwei,
	signature primitives.BLSSignature,
	index uint64,
) *Deposit {
	return &Deposit{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
		Signature:   signature,
		Index:       index,
	}
}
