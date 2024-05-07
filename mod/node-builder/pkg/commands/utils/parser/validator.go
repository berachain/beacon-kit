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

package parser

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ConvertPubkey converts a string to a public key.
func ConvertPubkey(pubkey string) (crypto.BLSPubkey, error) {
	// convert the public key to a BLSPubkey.
	pubkeyBytes, err := DecodeFrom0xPrefixedString(pubkey)
	if err != nil {
		return crypto.BLSPubkey{}, err
	}
	if len(pubkeyBytes) != constants.BLSPubkeyLength {
		return crypto.BLSPubkey{}, ErrInvalidPubKeyLength
	}

	return crypto.BLSPubkey(pubkeyBytes), nil
}

// ConvertWithdrawalCredentials converts a string to a withdrawal credentials.
func ConvertWithdrawalCredentials(credentials string) (
	consensus.WithdrawalCredentials,
	error,
) {
	// convert the credentials to a WithdrawalCredentials.
	credentialsBytes, err := DecodeFrom0xPrefixedString(credentials)
	if err != nil {
		return consensus.WithdrawalCredentials{}, err
	}
	if len(credentialsBytes) != constants.RootLength {
		return consensus.WithdrawalCredentials{},
			ErrInvalidWithdrawalCredentialsLength
	}
	return consensus.WithdrawalCredentials(credentialsBytes), nil
}

// ConvertAmount converts a string to a deposit amount.
//
//nolint:mnd // lots of magic numbers
func ConvertAmount(amount string) (math.Gwei, error) {
	// Convert the amount to a Gwei.
	amountBigInt, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return 0, ErrInvalidAmount
	}
	return math.Gwei(amountBigInt.Uint64()), nil
}

// ConvertSignature converts a string to a signature.
func ConvertSignature(signature string) (crypto.BLSSignature, error) {
	// convert the signature to a BLSSignature.
	signatureBytes, err := DecodeFrom0xPrefixedString(signature)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	if len(signatureBytes) != constants.BLSSignatureLength {
		return crypto.BLSSignature{}, ErrInvalidSignatureLength
	}
	return crypto.BLSSignature(signatureBytes), nil
}

// ConvertVersion converts a string to a version.
func ConvertVersion(version string) (primitives.Version, error) {
	versionBytes, err := DecodeFrom0xPrefixedString(version)
	if err != nil {
		return primitives.Version{}, err
	}
	if len(versionBytes) != constants.DomainTypeLength {
		return primitives.Version{}, ErrInvalidVersionLength
	}
	return primitives.Version(versionBytes), nil
}

// ConvertGenesisValidatorRoot converts a string to a genesis validator root.
func ConvertGenesisValidatorRoot(root string) (primitives.Root, error) {
	rootBytes, err := DecodeFrom0xPrefixedString(root)
	if err != nil {
		return primitives.Root{}, err
	}
	if len(rootBytes) != constants.RootLength {
		return primitives.Root{}, ErrInvalidRootLength
	}
	return primitives.Root(rootBytes), nil
}
