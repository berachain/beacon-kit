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
	"github.com/berachain/beacon-kit/mod/errors"
	nodetypes "github.com/berachain/beacon-kit/mod/node-api/handlers/node/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/types"
)

// Syncing returns the syncing status of the beacon node.
func (h *Handler[ContextT]) Syncing(_ ContextT) (any, error) {
	node := h.backend.GetNode()
	if node == nil {
		return nil, errors.New("node is nil")
	}

	// Get blockStore for heights
	blockStore := node.BlockStore()
	latestHeight := blockStore.Height()
	baseHeight := blockStore.Base()

	// Get consensus reactor for sync status
	consensusReactor := node.ConsensusReactor()

	response := nodetypes.SyncingData{
		HeadSlot:     latestHeight,
		IsOptimistic: true,
		ELOffline:    false,
	}

	// Check if we're catching up by comparing heights
	if consensusReactor != nil {
		response.SyncDistance = latestHeight - baseHeight
	}
	response.IsSyncing = response.SyncDistance > 0

	return types.Wrap(&response), nil
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
