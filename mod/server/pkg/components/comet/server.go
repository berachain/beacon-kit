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

package comet

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	sdkcomet "cosmossdk.io/server/v2/cometbft"
	"cosmossdk.io/server/v2/cometbft/mempool"
	consensus "github.com/berachain/beacon-kit/mod/consensus/pkg"
	consensustypes "github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/spf13/viper"
)

// assert that CometBFTServer implements the ServerComponent interface
var _ serverv2.ServerComponent[
	types.Node[transaction.Tx], transaction.Tx,
] = (*Server[
	types.Node[transaction.Tx], transaction.Tx, any,
])(nil)

// Server is a wrapper around the Server from the Cosmos SDK.
type Server[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
] struct {
	*sdkcomet.CometBFTServer[NodeT, T]
	TxCodec transaction.Codec[T]
}

// Init wraps the default Init method and sets the PrepareProposal and
// ProcessProposal handlers.
func (s *Server[NodeT, T, ValidatorUpdateT]) Init(
	node NodeT, v *viper.Viper, logger log.Logger,
) error {
	if err := s.CometBFTServer.Init(node, v, logger); err != nil {
		return err
	}
	var middleware *nodecomponents.ABCIMiddleware
	registry := node.GetServiceRegistry()
	if err := registry.FetchService(&middleware); err != nil {
		return err
	}

	engine := consensus.NewEngine[T, ValidatorUpdateT](
		consensustypes.CometBFTConsensus,
		s.TxCodec,
		middleware,
	)

	s.CometBFTServer.App.SetMempool(mempool.NoOpMempool[T]{})
	s.CometBFTServer.App.SetPrepareProposalHandler(engine.Prepare)
	s.CometBFTServer.App.SetProcessProposalHandler(engine.Process)
	return nil
}
