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
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// Node represents the node application.
type Node struct {
	*app.BeaconApp

	// name and description of the application.
	name        string
	description string
}

// New returns a new Node.
func New[NodeT types.NodeI]() NodeT {
	return types.NodeI(&Node{}).(NodeT)
}

// SetAppName sets the name of the application.
func (n *Node) SetAppName(name string) {
	n.name = name
}

// SetAppDescription sets the description of the application.
func (n *Node) SetAppDescription(description string) {
	n.description = description
}

// SetApplication sets the application.
func (n *Node) SetApplication(a servertypes.Application) {
	//nolint:errcheck // BeaconApp is our servertypes.Application
	n.BeaconApp = a.(*app.BeaconApp)
}
