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
	"crypto/ecdsa"
	"math/big"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// EthAccount represents an Ethereum account.
type EthAccount struct {
	name string
	pk   *ecdsa.PrivateKey
}

// NewEthAccountFromHex creates a new Ethereum account from hex private key.
func NewEthAccountFromHex(
	name string,
	hexPk string,
) *EthAccount {
	pk, err := crypto.HexToECDSA(hexPk)
	if err != nil {
		panic(err)
	}
	return &EthAccount{
		name: name,
		pk:   pk,
	}
}

// NewEthAccount creates a new Ethereum account.
func NewEthAccount(
	name string,
	pk *ecdsa.PrivateKey,
) *EthAccount {
	return &EthAccount{
		name: name,
		pk:   pk,
	}
}

// Name returns the name of the account.
func (a EthAccount) Name() string {
	return a.name
}

// SignTx signs a transaction with the account's private key.
func (a EthAccount) SignTx(
	chainID *big.Int, tx *types.Transaction,
) (*types.Transaction, error) {
	return types.SignTx(
		tx, types.LatestSignerForChainID(chainID), a.pk,
	)
}

// SignerFunc returns a signer function for the account.
func (a EthAccount) SignerFunc(chainID *big.Int) bind.SignerFn {
	return func(
		addr common.Address, tx *types.Transaction,
	) (*types.Transaction, error) {
		if addr != a.Address() {
			return nil, errors.New("account not authorized to sign")
		}
		return a.SignTx(chainID, tx)
	}
}

// PublicKey returns the public key of the account.
func (a EthAccount) PublicKey() *ecdsa.PublicKey {
	return &a.pk.PublicKey
}

// PrivateKey returns the private key of the account.
func (a EthAccount) PrivateKey() *ecdsa.PrivateKey {
	return a.pk
}

// Address returns the address of the account.
func (a EthAccount) Address() common.Address {
	return crypto.PubkeyToAddress(a.pk.PublicKey)
}
