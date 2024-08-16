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
	"fmt"
)

// Subscription is a channel that receives events.
type Subscription[T BaseEvent] chan T

// NewSubscription creates a new subscription with an unbuffered channel as the
// underlying type.
func NewSubscription[T BaseEvent]() Subscription[T] {
	return make(chan T)
}

// Clear clears the subscription channel.
func (s Subscription[T]) Clear() {
	for len(s) > 0 {
		<-s
	}
}

// Listen will start a goroutine that will trigger the handler for each event
// in the subscription.
func (s Subscription[T]) Listen(
	ctx context.Context, handler func(T),
) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-s:
				go func(e T) {
					handler(e)
				}(event)
			}
		}
	}()
}

// Await will block until an event is received from the subscription or the
// context is canceled.
func (s Subscription[T]) Await(ctx context.Context) (T, error) {
	fmt.Println("AWAITING EVENT")
	select {
	case event := <-s:
		fmt.Println("GOT EVENT")
		return event, nil
	case <-ctx.Done():
		return *new(T), ctx.Err()
	}
}

// Chan returns the underlying channel of the subscription.
func (s Subscription[T]) Chan() chan T {
	return s
}
