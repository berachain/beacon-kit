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

package primitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SigVerificationFn is a function that verifies a signature.
type SigVerificationFn func(pubkey, message, signature []byte) bool

// DepositMessage as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#depositmessage
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./deposit_message.go -objs DepositMessage -include ./pkg/crypto,./withdrawal_credentials.go,./pkg/math/,./primitives.go,./pkg/bytes,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deposit_message.ssz.go
type DepositMessage struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey" ssz-max:"48"`

	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials" ssz-size:"32"`

	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
}

// CreateAndSignDepositMessage creates and signs a deposit message.
func CreateAndSignDepositMessage(
	forkData *ForkData,
	domainType DomainType,
	signer crypto.BLSSigner,
	credentials WithdrawalCredentials,
	amount math.Gwei,
) (*DepositMessage, crypto.BLSSignature, error) {
	domain, err := forkData.ComputeDomain(domainType)
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	depositMessage := DepositMessage{
		Pubkey:      signer.PublicKey(),
		Credentials: credentials,
		Amount:      amount,
	}

	signingRoot, err := ComputeSigningRoot(&depositMessage, domain)
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	signature, err := signer.Sign(signingRoot[:])
	if err != nil {
		return nil, crypto.BLSSignature{}, err
	}

	return &depositMessage, signature, nil
}

// VerifyDeposit verifies the deposit data when attempting to create a
// new validator from a given deposit.
func (d *DepositMessage) VerifyCreateValidator(
	forkData *ForkData,
	signature crypto.BLSSignature,
	isSignatureValid SigVerificationFn,
	domainType DomainType,
) error {
	domain, err := forkData.ComputeDomain(domainType)
	if err != nil {
		return err
	}

	signingRoot, err := ComputeSigningRoot(d, domain)
	if err != nil {
		return err
	}

	if !isSignatureValid(
		d.Pubkey[:],
		signingRoot[:],
		signature[:],
	) {
		return ErrDepositMessage
	}

	return nil
}
