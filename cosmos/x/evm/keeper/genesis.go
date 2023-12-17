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

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/cosmos/x/evm/types"
)

var GenesisHashKey = []byte("eth1_genesis_hash")

func (k *Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) error {
	fmt.Println("CALLING INIT GENESIS", data.Eth1GenesisHash)
	fmt.Println("CALLING INIT GENESIS", []byte(data.Eth1GenesisHash))
	ctx.KVStore(k.storeKey).Set(GenesisHashKey, []byte(data.Eth1GenesisHash))
	return nil
}

func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	bingBong := common.Bytes2Hex(ctx.KVStore(k.storeKey).Get(GenesisHashKey))
	fmt.Println("RETRIEVE GENESIS IN EXPORT GENESIS", bingBong)
	return &types.GenesisState{
		Eth1GenesisHash: bingBong,
	}
}
