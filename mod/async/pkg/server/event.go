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

package server

import (
	"context"

	"github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
)

// EventServer asyncronously dispatches events to subscribers.
type EventServer struct {
	publishers map[types.MessageID]types.Publisher
	logger     log.Logger[any]
}

// NewEventServer creates a new event server.
func NewEventServer() *EventServer {
	return &EventServer{
		publishers: make(map[types.MessageID]types.Publisher),
	}
}

// Dispatch dispatches the given event to the publisher with the given eventID.
func (es *EventServer) Publish(event types.MessageI) error {
	return es.publishers[event.ID()].Publish(event)
}

// Subscribe subscribes the given channel to the publisher with the given eventID.
// It will error if the channel type does not match the event type corresponding
// publisher.
func (es *EventServer) Subscribe(eventID types.MessageID, ch any) error {
	publisher, ok := es.publishers[eventID]
	if !ok {
		return ErrPublisherNotFound
	}
	return publisher.Subscribe(ch)
}

// Start starts the event server.
func (es *EventServer) Start(ctx context.Context) {
	for _, publisher := range es.publishers {
		go publisher.Start(ctx)
	}
}

// RegisterPublisher registers the given publisher with the given eventID.
// Any subsequent events with <eventID> dispatched to this EventServer must be
// consistent with the type expected by <publisher>.
func (es *EventServer) RegisterPublisher(
	eventID types.MessageID, publisher types.Publisher,
) {
	es.publishers[eventID] = publisher
}

func (es *EventServer) SetLogger(logger log.Logger[any]) {
	es.logger = logger
}
