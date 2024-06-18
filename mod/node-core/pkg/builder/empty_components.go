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

package builder

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/middleware"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
)

// These functions return the addresses to empty node components.
// Typically used when supplying empty values to depinject, or as receiving
// addresses for constructed depinject values.

// emptyABCIMiddleware return an address pointing to an empty BeaconState.
func emptyABCIMiddleware() *components.ABCIMiddleware {
	return &middleware.ABCIMiddleware[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		components.BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*types.ExecutionPayload,
		*genesis.Genesis[*types.Deposit, *types.ExecutionPayloadHeader],
	]{}
}

// emptyAppBuilder returns an address pointing to an empty AppBuilder.
func emptyAppBuilder() *sdkruntime.AppBuilder {
	return &sdkruntime.AppBuilder{}
}

// emptyServiceRegistry returns an address pointing to an empty
// service.Registry.
func emptyServiceRegistry() *service.Registry {
	return &service.Registry{}
}
