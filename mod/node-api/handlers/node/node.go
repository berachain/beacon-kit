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

package node

import (
	"context"
	"fmt"
	"strconv"

	nodetypes "github.com/berachain/beacon-kit/mod/node-api/handlers/node/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
	cmtclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

// Syncing returns the syncing status of the beacon node.
func (h *Handler[ContextT]) Syncing(_ ContextT) (any, error) {
	// Create a new RPC client
	rpcClient, err := cmtclient.New("tcp://localhost:26657")
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client: %w", err)
	}

	// Create client context with the RPC client
	clientCtx := client.Context{}
	clientCtx = clientCtx.WithClient(rpcClient)

	// Query CometBFT status
	status, err := cmtservice.GetNodeStatus(context.Background(), clientCtx)
	if err != nil {
		return nil, fmt.Errorf("err in getting node status %w", err)
	}

	type SyncingResponse struct {
		Data struct {
			HeadSlot     string `json:"head_slot"`
			SyncDistance string `json:"sync_distance"`
			IsSyncing    bool   `json:"is_syncing"`
			IsOptimistic bool   `json:"is_optimistic"`
			ELOffline    bool   `json:"el_offline"`
		} `json:"data"`
	}

	response := SyncingResponse{}
	response.Data.HeadSlot = strconv.FormatInt(
		status.SyncInfo.LatestBlockHeight,
		10,
	)

	// Calculate sync distance
	if status.SyncInfo.LatestBlockHeight < status.SyncInfo.EarliestBlockHeight {
		syncDistance := status.SyncInfo.EarliestBlockHeight -
			status.SyncInfo.LatestBlockHeight
		response.Data.SyncDistance = strconv.FormatInt(syncDistance, 10)
		response.Data.IsSyncing = status.SyncInfo.CatchingUp
	} else {
		response.Data.SyncDistance = "0"
		response.Data.IsSyncing = false
	}

	// Keep existing values for these fields
	response.Data.IsOptimistic = true
	response.Data.ELOffline = false

	return response, nil
}

// Version returns the version of the beacon node.
func (h *Handler[ContextT]) Version(_ ContextT) (any, error) {
	version, err := h.backend.GetNodeVersion()
	if err != nil {
		return nil, err
	}
	return types.Wrap(nodetypes.VersionData{
		Version: version,
	}), nil
}
