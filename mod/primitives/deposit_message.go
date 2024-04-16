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

// SigVerificationFn is a function that verifies a signature.
type SigVerificationFn func(pubkey, message, signature []byte) bool

// DepositMessage as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#depositmessage
//
//nolint:lll
//go:generate go run github.com/ferranbt/fastssz/sszgen --path ./deposit_message.go -objs DepositMessage -include ./withdrawal_credentials.go,./u64.go,./primitives.go,./bytes.go,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deposit_message.ssz.go
type DepositMessage struct {
	// Public key of the validator specified in the deposit.
	Pubkey BLSPubkey `json:"pubkey" ssz-max:"48"`

	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials" ssz-size:"32"`

	// Deposit amount in gwei.
	Amount Gwei `json:"amount"`
}

// VerifyDeposit verifies the deposit data when attempting to create a
// new validator from a given deposit.
func (d *DepositMessage) VerifyCreateValidator(
	forkData *ForkData,
	signature BLSSignature,
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
