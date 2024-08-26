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

package state

import (
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
)

type Manager struct {
	cms storetypes.CommitMultiStore
}

// NewManager creates a new Manager.
func NewManager(
	db dbm.DB,
	logger log.Logger,
	opts ...func(*Manager),
) *Manager {
	sm := &Manager{cms: store.NewCommitMultiStore(
		db,
		logger,
		storemetrics.NewNoOpMetrics(),
	)}
	for _, opt := range opts {
		opt(sm)
	}
	return sm
}

// CommitMultiStore returns the CommitMultiStore of the Manager.
// TODO:REMOVE
func (sm *Manager) CommitMultiStore() storetypes.CommitMultiStore {
	return sm.cms
}
