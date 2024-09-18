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

package parser

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ConvertPubkey converts a string to a public key.
func ConvertPubkey(pubkey string) (crypto.BLSPubkey, error) {
	// convert the public key to a BLSPubkey.
	pubkeyBytes, err := hex.ToBytes(pubkey)
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
	types.WithdrawalCredentials,
	error,
) {
	// convert the credentials to a WithdrawalCredentials.
	credentialsBytes, err := hex.ToBytes(credentials)
	if err != nil {
		return types.WithdrawalCredentials{}, err
	}
	if len(credentialsBytes) != constants.RootLength {
		return types.WithdrawalCredentials{},
			ErrInvalidWithdrawalCredentialsLength
	}
	return types.WithdrawalCredentials(credentialsBytes), nil
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
	signatureBytes, err := hex.ToBytes(signature)
	if err != nil {
		return crypto.BLSSignature{}, err
	}
	if len(signatureBytes) != constants.BLSSignatureLength {
		return crypto.BLSSignature{}, ErrInvalidSignatureLength
	}
	return crypto.BLSSignature(signatureBytes), nil
}

// ConvertVersion converts a string to a version.
func ConvertVersion(version string) (common.Version, error) {
	versionBytes, err := hex.ToBytes(version)
	if err != nil {
		return common.Version{}, err
	}
	if len(versionBytes) != constants.DomainTypeLength {
		return common.Version{}, ErrInvalidVersionLength
	}
	return common.Version(versionBytes), nil
}

// ConvertGenesisValidatorRoot converts a string to a genesis validator root.
func ConvertGenesisValidatorRoot(root string) (common.Root, error) {
	rootBytes, err := hex.ToBytes(root)
	if err != nil {
		return common.Root{}, err
	}
	if len(rootBytes) != constants.RootLength {
		return common.Root{}, ErrInvalidRootLength
	}
	return common.Root(rootBytes), nil
}
