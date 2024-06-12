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

package components

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	consensus "github.com/berachain/beacon-kit/mod/sync/pkg/consensus"
	"github.com/ethereum/go-ethereum/event"
)

// CLSyncServiceInput is the input for the CLSyncService provider.
type CLSyncServiceInput[SubscriptionT interface {
	// Unsubscribe terminates the subscription.
	Unsubscribe()
}] struct {
	depinject.In
	Logger   log.Logger
	SyncFeed *event.FeedOf[*feed.Event[bool]]
}

// ProvideCLSyncService is a depinject provider that returns a CLSyncService.
func ProvideCLSyncService[
	SubscriptionT interface {
		Unsubscribe()
	}](
	in CLSyncServiceInput[SubscriptionT],
) (*consensus.SyncService[SubscriptionT], error) {
	// Build the CLSyncService.
	return consensus.NewSyncService[SubscriptionT](
		in.SyncFeed,
		in.Logger.With("service", "cl-sync"),
	), nil
}
