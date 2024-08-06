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

package dispatcher

import (
	"context"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
)

// EventServer asyncronously dispatches events to subscribers.
type EventServer struct {
	feeds map[types.EventID]publisher
}

// NewEventServer creates a new event server.
func NewEventServer() *EventServer {
	return &EventServer{
		feeds: make(map[types.EventID]publisher),
	}
}

// Dispatch dispatches the given event to the feed with the given eventID.
func (es *EventServer) Dispatch(ctx context.Context, event *types.Event[any]) {
	es.feeds[event.Type()].Publish(ctx, *event)
}

// Subscribe subscribes the given channel to the feed with the given eventID.
// It will error if the channel type does not match the event type corresponding
// feed.
func (es *EventServer) Subscribe(eventID types.EventID, ch chan any) error {
	feed, ok := es.feeds[eventID]
	if !ok {
		return ErrFeedNotFound
	}
	return feed.Subscribe(ch)
}

// Start starts the event server.
func (es *EventServer) Start(ctx context.Context) {
	for _, feed := range es.feeds {
		go feed.Start(ctx)
	}
}

// RegisterFeed registers the given feed with the given eventID.
// Any subsequent events with <eventID> dispatched to this EventServer must be
// consistent with the type expected by <feed>.
func (es *EventServer) RegisterFeed(eventID types.EventID, feed publisher) {
	es.feeds[eventID] = feed
}
