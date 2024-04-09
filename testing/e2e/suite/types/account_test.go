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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestNewEthAccount(t *testing.T) {
	name := "testAccount"
	hexPk := "f561d45d1df30a6556d30e39f97011faa3632e43cd378224ad6cc83bb8aea3e6"

	account := types.NewEthAccountFromHex(name, hexPk)
	require.NotNil(t, account.PrivateKey())
	require.Equal(t, name, account.Name())
}

func TestEthAccount_PublicKey(t *testing.T) {
	hexPk := "f561d45d1df30a6556d30e39f97011faa3632e43cd378224ad6cc83bb8aea3e6"
	expectedPrivateKey, _ := crypto.HexToECDSA(hexPk)
	require.NotNil(t, expectedPrivateKey)

	account := types.NewEthAccountFromHex("testAccount", hexPk)
	publicKey := account.PublicKey()

	require.Equal(t, expectedPrivateKey.PublicKey, *publicKey)
}

func TestEthAccount_Address(t *testing.T) {
	hexPk := "f561d45d1df30a6556d30e39f97011faa3632e43cd378224ad6cc83bb8aea3e6"
	account := types.NewEthAccountFromHex("testAccount", hexPk)
	address := account.Address()

	expectedAddress := common.HexToAddress(
		"0x9e0437E86fE55b4D6Bb1191CbACdA384775fF63D",
	)
	require.Equal(t, expectedAddress, address)
}
