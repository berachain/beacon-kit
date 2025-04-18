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

package types

import (
	"context"

	"cosmossdk.io/store"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	service "github.com/berachain/beacon-kit/node-core/services/registry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Node defines the API for the node application.
// It extends the Application interface from the Cosmos SDK.
type Node interface {
	CommitMultistoreAccessor
	StorageBackendAccessor

	Start(context.Context) error
}

// CommitMultistoreAccessor allows access to the commit multistore.
// This is required by commands like rollback.
type CommitMultistoreAccessor interface {
	CommitMultiStore() store.CommitMultiStore
}

// StorageBackendAccessor allows access to the storage backend.
// This is required by commands like db-check.
type StorageBackendAccessor interface {
	StorageBackend() blockchain.StorageBackend
}

// ConsensusService defines everything we utilise externally from CometBFT.
type ConsensusService interface {
	service.Basic
	CreateQueryContext(
		height int64,
		prove bool,
	) (sdk.Context, error)
	LastBlockHeight() int64
}
