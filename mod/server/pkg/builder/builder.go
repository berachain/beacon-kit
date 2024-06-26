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

package serverbuilder

import (
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	serverv2 "cosmossdk.io/server/v2"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
)

type ServerBuilder[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
] struct {
	server *serverv2.Server[NodeT, T]

	depInjectCfg depinject.Config
	Components   []any
}

func New[
	NodeT types.Node[T], T transaction.Tx, ValidatorUpdateT any,
](
	opts ...Opt[NodeT, T, ValidatorUpdateT],
) *ServerBuilder[NodeT, T, ValidatorUpdateT] {
	b := &ServerBuilder[NodeT, T, ValidatorUpdateT]{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (sb *ServerBuilder[
	NodeT, T, ValidatorUpdateT,
]) Build(logger log.Logger) *serverv2.Server[NodeT, T] {
	// var (
	// logger    log.Logger
	// cmtServer *components.CometBFTServer[NodeT, T, ValidatorUpdateT]
	// )
	// if err := depinject.Inject(
	// 	depinject.Configs(
	// 		sb.depInjectCfg,
	// 		depinject.Provide(
	// 			sb.components...,
	// 		),
	// 	),
	// 	&logger,
	// 	&cmtServer,
	// ); err != nil {
	// 	panic(err)
	// }

	// sb.server = serverv2.NewServer(
	// 	logger,
	// 	sb.components...,
	// )

	return sb.server
}
