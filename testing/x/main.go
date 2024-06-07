// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package main

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func genKey(key string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func main() {
	chainID := big.NewInt(80084)

	keyStrs := []string{
		"8bdb3424a89528f4b6eaa6b86a9afff947b02b2a6069d0cec203d9d381fd9c2a",
		"a0684317c46e1bb1faa9050f74bedbeecde68b51ba861df503d9abb8c2e4545f",
		"1261bffac4755ebe6c2006e44e71863828b77be1c264653343d58b580afee70e",
		"82ccf4904ceeb7d4f4ae9151868fc7c4dfa64fe09f1e9f544231d430fe8ef845",
		"cbb5a7450acd848c5888d50eec3f9ecfc924953ab88df915b02a425df9e12f51",
		"467dd98370170dfb5f5088c9802b49ade4a2c97fdc71c81aaae511e2bb8883b8",
		"358cf3d3cf27f36ae752e773190046ddfeffb73b4156f0e30acd909211006564",
	}
	dials := []string{
		"https://restless-thrumming-county.bera-bartio.quiknode.pro/6f5c8dc2120be6048421ac6d84c1f700e5875e50",
		"https://bartio-eth-rpc-internal.berachain-devnet.com/",
		// "https://artio-internal-rpc.berachain.com",
		// "https://artio-sentry-rpc.berachain.com",
		// "http://localhost:10545",
		// "https://bera-testnet-evm-rpc.staketab.org:443",
		// "https://berachain-testnet-evm.synergynodes.com",
		// "https://artio.rpc.berachain.com",
		// "https://rpc-beratestnet-1.cosmos-spaces.cloud",
		// "https://rpc-beratestnet-2.cosmos-spaces.cloud",
	}
	clients := make([]*ethclient.Client, len(dials))
	ctx := context.Background()
	for i, url := range dials {
		var err error
		clients[i], err = ethclient.Dial(url)
		if err != nil {
			panic(err)
		}
	}

	signer := types.NewLondonSigner(chainID)

	to := common.HexToAddress("0x7Ad4F487Fb23902bDAa1B885E64d5893d81dEA98")
	wg := sync.WaitGroup{}

	keys := make([]*ecdsa.PrivateKey, len(keyStrs))
	for i, key := range keyStrs {
		keys[i] = genKey(key)
	}
	wg.Add(len(keys))

	for _, key := range keys {
		nonce, err := clients[1].NonceAt(
			ctx, crypto.PubkeyToAddress(key.PublicKey),
			big.NewInt(int64(rpc.LatestBlockNumber)))
		if err != nil {
			panic(err)
		}

		go func(_nonce uint64) {
			for i := range [50000]int{} {
				tip, _ := clients[1].SuggestGasTipCap(ctx)
				fmt.Println(tip)
				gas, _ := clients[1].SuggestGasPrice(ctx)

				defer wg.Done()
				tx, err := types.SignNewTx(key, signer, &types.LegacyTx{
					Nonce:    _nonce + uint64(i),
					GasPrice: big.NewInt(0).Add(gas, big.NewInt(1e10)),
					Gas:      1954688,
					To:       &to,
					Data: func() []byte {
						b := make([]byte, 12000)
						rand.Read(b)
						return b
					}(),
				})
				if err != nil {
					panic(err)
				}

				if err := clients[i%len(clients)].SendTransaction(ctx, tx); err != nil {
					fmt.Println("ERROR: ", err)
				} else {
					fmt.Println("sent: ", tx.Hash(), tx.Nonce())
				}
			}
		}(nonce)
	}

	wg.Wait()
}
