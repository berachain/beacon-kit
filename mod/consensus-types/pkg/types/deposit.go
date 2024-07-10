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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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
	merkleizer := merkle.NewMerkleizer[[32]byte, *Deposit]()
	return merkleizer.MerkleizeListComposite(
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
