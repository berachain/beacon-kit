// SPDX-License-Identifier: MIT
//
// # Copyright (c) 2023 Berachain Foundation
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
//

// nolint
//
//nolint:nolintlint // testing file.
package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/itsdevbear/bolaris/app/contracts"
	"github.com/itsdevbear/bolaris/beacon/execution/logs/callback"
	evmv1 "github.com/itsdevbear/bolaris/types/evm/v1"
)

func main() {
	ssc := &contracts.StakingCallbacks{}

	sc, err := callback.NewFrom(ssc)
	if err != nil {
		panic(err)
	}

	ethclient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	logs, err := ethclient.FilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress("0xB0ce0be267f1B1db9b30CD3E61DF1C6937129A84")},
		FromBlock: big.NewInt(135),
		ToBlock:   big.NewInt(1000),
	})

	if err != nil {
		panic(err)
	}

	for _, log := range logs {
		// Handle the log
		err = sc.HandleLog(context.Background(), evmv1.NewLogFromGethLog(log))
		if err != nil {
			panic(err)
		}
	}

	_ = sc
}
