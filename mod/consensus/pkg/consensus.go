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

package consensus

import (
	"context"

	"cosmossdk.io/core/transaction"
	"cosmossdk.io/server/v2/cometbft/handlers"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/types"
	"github.com/cosmos/gogoproto/proto"
)

// Engine is the interface that must be implemented by all consensus
// engines.
type Engine[T transaction.Tx, ValidatorUpdateT any] interface {
	InitGenesis(ctx context.Context, genesisBz []byte) ([]ValidatorUpdateT, error)
	Prepare(context.Context, handlers.AppManager[T], []T, proto.Message) ([]T, error)
	Process(context.Context, handlers.AppManager[T], []T, proto.Message) error
	PreBlock(ctx context.Context, msg proto.Message) error
	EndBlock(ctx context.Context) ([]ValidatorUpdateT, error)
}

var _ Engine[transaction.Tx, any] = (*cometbft.ConsensusEngine[transaction.Tx, any])(nil)

func NewEngine[T transaction.Tx, ValidatorUpdateT any](
	version types.ConsensusVersion,
	txCodec transaction.Codec[T],
	m types.Middleware,
) Engine[T, ValidatorUpdateT] {
	switch version {
	case types.CometBFTConsensus:
		return cometbft.NewConsensusEngine[T, ValidatorUpdateT](
			txCodec, m,
		)
	case types.RollKitConsensus:
		panic("not implemented")
	default:
		panic("unknown consensus version")
	}
}
