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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./deposit.go -objs Deposit -include ../../../primitives/pkg/common,./withdrawal_credentials.go,../../../primitives/pkg/math,../../../primitives/pkg/bytes,../../../primitives/pkg/crypto,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deposit.ssz.go
//nolint:lll // struct tags.
type Deposit struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"      ssz-max:"48"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials"              ssz-size:"32"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
	// Signature of the deposit data.
	Signature crypto.BLSSignature `json:"signature"   ssz-max:"96"`
	// Index of the deposit in the deposit contract.
	Index uint64 `json:"index"`
}

// NewDeposit creates a new Deposit instance.
func NewDeposit(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
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

// New creates a new Deposit instance.
func (d *Deposit) New(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
	index uint64,
) *Deposit {
	return NewDeposit(
		pubkey, credentials, amount, signature, index,
	)
}

// Deposits is a typealias for a list of Deposits.
type Deposits []*Deposit

// HashTreeRoot returns the hash tree root of the Withdrawals list.
func (d Deposits) HashTreeRoot() (common.Root, error) {
	// TODO: read max deposits from the chain spec.
	return ssz.MerkleizeListComposite[any, math.U64](
		d, constants.MaxDepositsPerBlock,
	)
}

// VerifySignature verifies the deposit data and signature.
func (d *Deposit) VerifySignature(
	forkData *ForkData,
	domainType common.DomainType,
	signatureVerificationFn func(
		pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
	) error,
) error {
	return (&DepositMessage{
		Pubkey:      d.Pubkey,
		Credentials: d.Credentials,
		Amount:      d.Amount,
	}).VerifyCreateValidator(
		forkData, d.Signature,
		domainType, signatureVerificationFn,
	)
}

// GetAmount returns the deposit amount in gwei.
func (d *Deposit) GetAmount() math.Gwei {
	return d.Amount
}

// GetPubkey returns the public key of the validator specified in the deposit.
func (d *Deposit) GetPubkey() crypto.BLSPubkey {
	return d.Pubkey
}

// GetIndex returns the index of the deposit in the deposit contract.
func (d *Deposit) GetIndex() uint64 {
	return d.Index
}

// GetSignature returns the signature of the deposit data.
func (d *Deposit) GetSignature() crypto.BLSSignature {
	return d.Signature
}

// GetWithdrawalCredentials returns the staking credentials of the deposit.
func (d *Deposit) GetWithdrawalCredentials() WithdrawalCredentials {
	return d.Credentials
}
