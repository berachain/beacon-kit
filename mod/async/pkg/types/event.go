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
	"errors"
)

// EventID represents the type of an event.
type EventID string

// Event represents a generic event in the beacon chain.
type Event[DataT any] struct {
	// ctx is the context associated with the event.
	ctx context.Context
	// eventType is the name of the event.
	eventType EventID
	// event is the actual beacon event.
	data DataT
	// error is the error associated with the event.
	err error
}

// NewEvent creates a new Event with the given context and beacon event.
func NewEvent[
	DataT any,
](
	ctx context.Context, eventType EventID, data DataT, errs ...error,
) *Event[DataT] {
	return &Event[DataT]{
		ctx:       ctx,
		eventType: eventType,
		data:      data,
		err:       errors.Join(errs...),
	}
}

// Type returns the type of the event.
func (e Event[DataT]) Type() EventID {
	return e.eventType
}

// Context returns the context associated with the event.
func (e Event[DataT]) Context() context.Context {
	return e.ctx
}

// Data returns the data associated with the event.
func (e Event[DataT]) Data() DataT {
	return e.data
}

// Error returns the error associated with the event.
func (e Event[DataT]) Error() error {
	return e.err
}

// Is returns true if the event has the given type.
func (e Event[DataT]) Is(eventType EventID) bool {
	return e.eventType == eventType
}
