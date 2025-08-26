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
	"runtime"

	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/node/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// Syncing is a placeholder so that beacon API clients don't break.
//
// TODO: Implement with real data.
func (h *Handler) Syncing(handlers.Context) (any, error) {
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
	response.Data.HeadSlot = "0"
	response.Data.SyncDistance = "1"
	response.Data.IsSyncing = false
	response.Data.IsOptimistic = true
	response.Data.ELOffline = false

	return response, nil
}

// Version is a placeholder so that beacon API clients don't break.
// Note: cometVersionInfo comes from git describe (via the build process)
// Git describe usually returns string like <v1.2.3-14-g2414721>, which can be understood as follow:
// - v1.2.3 → the most recent tag reachable from this commit
// - 14 → number of commits since that tag
// - g2414721 → the abbreviated commit hash, **with a leading g**.
// That g stands for “git”. It’s a prefix Git uses to distinguish the commit hash from other possible identifiers.
func (h *Handler) Version(handlers.Context) (any, error) {
	cometVersionInfo := version.NewInfo() // same used in beacond version command
	r := types.VersionData{
		Version: fmt.Sprintf("%s/%s (%s %s)",
			cometVersionInfo.AppName,
			cometVersionInfo.Version,
			runtime.GOOS,
			runtime.GOARCH,
		),
	}

	return types.Wrap(r), nil
}
