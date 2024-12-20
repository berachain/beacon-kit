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
