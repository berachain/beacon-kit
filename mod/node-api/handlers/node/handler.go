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
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/server/context"
	cmtclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
)

const RpcEndpoint = "tcp://localhost:26657"

// Handler is the handler for the node API.
type Handler[ContextT context.Context] struct {
	*handlers.BaseHandler[ContextT]
	backend   Backend
	rpcClient *cmtclient.HTTP
	clientCtx client.Context
}

func NewHandler[ContextT context.Context](backend Backend) *Handler[ContextT] {
	rpcClient, err := cmtclient.New(RpcEndpoint)
	if err != nil {
		// Not returning error to keep the pattern same across all handlers.
		rpcClient = nil
	}
	clientCtx := client.Context{}.WithClient(rpcClient)

	h := &Handler[ContextT]{
		BaseHandler: handlers.NewBaseHandler(
			handlers.NewRouteSet[ContextT](""),
		),
		backend:   backend,
		rpcClient: rpcClient,
		clientCtx: clientCtx,
	}
	return h
}
