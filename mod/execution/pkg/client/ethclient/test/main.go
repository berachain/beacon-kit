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

// package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
// 	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
// 	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient/rpc"
// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
// 	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
// 	"github.com/davecgh/go-spew/spew"
// )

// func main() {
// 	c := rpc.NewClient("https://restless-thrumming-county.bera-bartio.quiknode.pro/6f5c8dc2120be6048421ac6d84c1f700e5875e50")
// 	cl := ethclient.New[*types.ExecutionPayload](c)

// 	logs, err := cl.GetLogsAtBlockNumber(
// 		context.TODO(),
// 		math.U64(2870518),
// 		common.NewExecutionAddressFromHex("0x4242424242424242424242424242424242424242"),
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	for _, l := range logs {
// 		spew.Dump(l)
// 		deposit := new(types.Deposit)
// 		deposit.UnmarshalLog(l)
// 		spew.Dump(deposit)
// 	}
// }
