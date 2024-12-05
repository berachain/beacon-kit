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

package types

import (
	"context"

	"github.com/berachain/beacon-kit/primitives/async"
)

// Broker is the unique publisher for broadcasting all events corresponding
// to the <eventID> to all registered client channels.
// There should be exactly one broker responsible for every eventID.
type Broker interface {
	// EventID returns the event ID that the publisher is responsible for.
	EventID() async.EventID
	// Start starts the broker loop in a goroutine to listen and broadcast
	// events to all subscribers.
	Start(ctx context.Context)
	// Publish publishes a msg to all subscribers.
	// Errors if the message is not of type T, or if the context is canceled.
	Publish(event async.BaseEvent) error
	// Subscribe registers the provided channel to the broker.
	// Errors if the channel is not of type chan T.
	// Contract: the channel must be a Chan[T], where T is the expected
	// type of the event data.
	Subscribe(ch any) error
	// Unsubscribe removes a client from the broker.
	// Returns an error if the provided channel is not of type chan T.
	Unsubscribe(ch any) error
}
