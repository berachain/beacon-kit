// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"fmt"

	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/node/types"
)

// Syncing returns the syncing status of the beacon node.
func (h *Handler) Syncing(handlers.Context) (any, error) {
	latestHeight, syncingToHeight := h.backend.GetSyncData()
	syncDistance := syncingToHeight - latestHeight
	response := &types.SyncingData{
		HeadSlot:     latestHeight,
		SyncDistance: syncDistance,

		// BeaconKit has two operation modes: syncing from genesis
		// and normal operations. Every block finalized is verified
		// by the EL so for every purpose IsSyncing is equivalent to syncDistance > 0
		IsSyncing: syncDistance > 0,

		// BeaconKit always verifies blocks payload, whether
		// it is syncing or it's in normal operation mode.
		// So IsOptimistic is always false
		IsOptimistic: false,

		// BeaconKit fails to verify and finalize blocks if
		// the EL is not reachable, resulting in a panic state.
		// TODO: consider exposing the EL state here
		ELOffline: false,
	}

	return types.Wrap(response), nil
}

// Note: version comes from git describe (via the build process)
// Git describe usually returns string like <v1.2.3-14-g2414721>, which can be understood as follow:
// - v1.2.3 → the most recent tag reachable from this commit
// - 14 → number of commits since that tag
// - g2414721 → the abbreviated commit hash, **with a leading g**.
// That g stands for “git”. It’s a prefix Git uses to distinguish the commit hash from other possible identifiers.
func (h *Handler) Version(handlers.Context) (any, error) {
	appName, version, os, arch := h.backend.GetVersionData()
	r := &types.VersionData{
		Version: fmt.Sprintf("%s/%s (%s %s)",
			appName,
			version,
			os,
			arch,
		),
	}

	return types.Wrap(r), nil
}
