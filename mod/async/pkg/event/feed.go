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

package event

import "github.com/ethereum/go-ethereum/event"

// Subscription is a subscription to a feed.
type Subscription = event.Subscription

// FeedOf is a feed of events.
// It is a wrapper around the event.FeedOf type.
type FeedOf[
	E ~uint8,
	T interface {
		Type() E
	},
] struct {
	event.FeedOf[T]
	ids map[E]struct{}
}

// NewFeedOf creates a new feed of events with the given event IDs.
func NewFeedOf[
	E ~uint8, T interface {
		Type() E
	},
](ids ...E) *FeedOf[E, T] {
	_ids := make(map[E]struct{})
	for _, id := range ids {
		_ids[id] = struct{}{}
	}
	return &FeedOf[E, T]{
		FeedOf: event.FeedOf[T]{},
		ids:    _ids,
	}
}

// TODO: Upgrade later to avoid the caller having to allocate the channel, good
// dev ux. // Subscribe subscribes to the feed and returns a channel and a
// subscription.=
// func (f *FeedOf[E, T]) Subscribe() (chan<- T, event.Subscription) {
// 	ch := make(chan T, 1)
// 	sub := f.FeedOf.Subscribe(make(chan T, 1))
// 	return ch, sub
// }

// Send sends an event to the feed.
func (f *FeedOf[E, T]) Send(event T) int {
	// Add a safety check to ensure that the event is registered with
	// the feed.
	if _, ok := f.ids[event.Type()]; !ok {
		return 0
	}
	return f.FeedOf.Send(event)
}
