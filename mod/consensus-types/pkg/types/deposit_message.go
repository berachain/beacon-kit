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
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// SigVerificationFn defines a function type for verifying a signature.
type SigVerificationFn func(
	pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
) error

// DepositMessage represents a deposit message as defined in the Ethereum 2.0
// specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#depositmessage
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./deposit_message.go -objs DepositMessage -include ./withdrawal_credentials.go,../../../primitives/pkg/math,../../../primitives/pkg/crypto,./fork_data.go,../../../primitives/pkg/bytes,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deposit_message.ssz.go
type DepositMessage struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"      ssz-max:"48"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials"              ssz-size:"32"`
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
	domain, err := forkData.ComputeDomain(domainType)
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	depositMessage := &DepositMessage{
		Pubkey:      signer.PublicKey(),
		Credentials: credentials,
		Amount:      amount,
	}

	signingRoot, err := ssz.ComputeSigningRoot(depositMessage, domain)
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	signature, err := signer.Sign(signingRoot[:])
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	return depositMessage, signature, nil
}

// VerifyDeposit verifies the deposit data when attempting to create a
// new validator from a given deposit.
func (d *DepositMessage) VerifyCreateValidator(
	forkData *ForkData,
	signature crypto.BLSSignature,
	signatureVerificationFn SigVerificationFn,
	domainType common.DomainType,
) error {
	domain, err := forkData.ComputeDomain(domainType)
	if err != nil {
		return err
	}

	signingRoot, err := ssz.ComputeSigningRoot(d, domain)
	if err != nil {
		return err
	}

	if err = signatureVerificationFn(
		d.Pubkey, signingRoot[:], signature,
	); err != nil {
		return errors.Join(err, ErrDepositMessage)
	}

	return nil
}
