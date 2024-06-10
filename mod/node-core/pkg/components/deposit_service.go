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

package components

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/ethereum/go-ethereum/event"
)

// DepositServiceIn is the input for the deposit service.
type DepositServiceIn struct {
	depinject.In

	Logger                log.Logger
	ChainSpec             primitives.ChainSpec
	EngineClient          *engineclient.EngineClient[*types.ExecutionPayload]
	TelemetrySink         *metrics.TelemetrySink
	DepositStore          *depositdb.KVStore[*types.Deposit]
	BeaconDepositContract *deposit.WrappedBeaconDepositContract[
		*types.Deposit, types.WithdrawalCredentials,
	]
	BlockFeed *event.FeedOf[*feed.Event[*types.BeaconBlock]]
}

// ProvideDepositService provides the deposit service to the depinject
// framework.
func ProvideDepositService(in DepositServiceIn) *deposit.Service[
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	*feed.Event[*types.BeaconBlock],
	*types.Deposit,
	*types.ExecutionPayload,
	event.Subscription,
	types.WithdrawalCredentials,
] {
	// Build the deposit service.
	return deposit.NewService[
		*types.BeaconBlockBody,
		*types.BeaconBlock,
		*feed.Event[*types.BeaconBlock],
		*depositdb.KVStore[*types.Deposit],
		*types.ExecutionPayload,
		event.Subscription,
	](
		in.Logger.With("service", "deposit"),
		math.U64(in.ChainSpec.Eth1FollowDistance()),
		in.EngineClient,
		in.TelemetrySink,
		in.DepositStore,
		in.BeaconDepositContract,
		in.BlockFeed,
	)
}
