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
	"crypto/ecdsa"
	"errors"
	"math/big"

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

// NewEthAccount creates a new Ethereum account.
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

// PrivateKey returns the private key of the account.
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
