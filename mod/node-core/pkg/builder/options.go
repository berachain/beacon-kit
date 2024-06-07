// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package nodebuilder

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

type Opt[NodeT types.NodeI] func(*NodeBuilder[NodeT])

func WithName[NodeT types.NodeI](name string) Opt[NodeT] {
	return func(nb *NodeBuilder[NodeT]) {
		nb.node.SetAppName(name)
		nb.name = name
	}
}

func WithDescription[NodeT types.NodeI](description string) Opt[NodeT] {
	return func(nb *NodeBuilder[NodeT]) {
		nb.node.SetAppDescription(description)
		nb.description = description
	}
}

func WithDepInjectConfig[NodeT types.NodeI](cfg depinject.Config) Opt[NodeT] {
	return func(nb *NodeBuilder[NodeT]) {
		nb.depInjectCfg = cfg
	}
}

func WithChainSpec[NodeT types.NodeI](cs primitives.ChainSpec) Opt[NodeT] {
	return func(nb *NodeBuilder[NodeT]) {
		nb.chainSpec = cs
	}
}
