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
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

//go:generate sszgen --path . --include ../../primitives/common,../../primitives/bytes,../../primitives/math,../../primitives/crypto --objs DepositMessage --output deposit_message_sszgen.go

// DepositMessage represents a deposit message as defined in the Ethereum 2.0
// specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#depositmessage
type DepositMessage struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey" ssz-size:"48"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials" ssz-size:"32"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
}

// CreateAndSignDepositMessage constructs and signs a deposit message.
func CreateAndSignDepositMessage(
	forkData *ForkData,
	domainType common.DomainType,
	signer crypto.BLSSigner,
	credentials WithdrawalCredentials,
	amount math.Gwei,
) (*DepositMessage, crypto.BLSSignature, error) {
	domain := forkData.ComputeDomain(domainType)
	depositMessage := &DepositMessage{
		Pubkey:      signer.PublicKey(),
		Credentials: credentials,
		Amount:      amount,
	}
	hashTreeRoot, err := depositMessage.HashTreeRoot()
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}
	signingRoot := ComputeSigningRootFromHTR(hashTreeRoot, domain)
	signature, err := signer.Sign(signingRoot[:])
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	return depositMessage, signature, nil
}

// VerifyCreateValidator verifies the deposit data when attempting to create a
// new validator from a given deposit.
func (dm *DepositMessage) VerifyCreateValidator(
	forkData *ForkData,
	signature crypto.BLSSignature,
	domainType common.DomainType,
	signatureVerificationFn func(
		pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
	) error,
) error {
	hashTreeRoot, err := dm.HashTreeRoot()
	if err != nil {
		return err
	}
	signingRoot := ComputeSigningRootFromHTR(
		hashTreeRoot, forkData.ComputeDomain(domainType))
	if err := signatureVerificationFn(
		dm.Pubkey, signingRoot[:], signature,
	); err != nil {
		return errors.Join(err, ErrDepositMessage)
	}

	return nil
}
