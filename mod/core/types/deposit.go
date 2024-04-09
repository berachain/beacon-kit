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
	"github.com/berachain/beacon-kit/mod/forks"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/davecgh/go-spew/spew"
	"github.com/itsdevbear/comet-bls12-381/bls/blst"
)

// Deposits is a typealias for a slice of Deposit.
type Deposits []*Deposit

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
//
//go:generate go run github.com/ferranbt/fastssz/sszgen --path deposit.go -objs Deposit,DepositMessage -include withdrawal_credentials.go,../../primitives,$GETH_PKG_INCLUDE/common -output deposit.ssz.go
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

// String returns a string representation of the Deposit.
func (d *Deposit) String() string {
	return spew.Sdump(d)
}

// DepositMessage as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#depositmessage
//
//nolint:lll
type DepositMessage struct {
	// Public key of the validator specified in the deposit.
	Pubkey primitives.BLSPubkey `json:"pubkey" ssz-max:"48"`

	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials" ssz-size:"32"`

	// Deposit amount in gwei.
	Amount primitives.Gwei `json:"amount"`
}

// VerifyDeposit verifies the deposit data when attempting to create a
// new validator from a given deposit.
func (d *DepositMessage) VerifyCreateValidator(
	forkData *forks.ForkData,
	signature primitives.BLSSignature,
) error {
	domain, err := forkData.ComputeDomain(primitives.DomainTypeDeposit)
	if err != nil {
		return err
	}

	signingRoot, err := primitives.ComputeSigningRoot(d, domain)
	if err != nil {
		return err
	}

	if !blst.VerifySignaturePubkeyBytes(
		d.Pubkey[:],
		signingRoot[:],
		signature[:],
	) {
		return ErrDepositMessage
	}

	return nil
}
