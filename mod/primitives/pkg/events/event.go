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

package events

import (
	"context"
	"errors"
)

// BaseEvent defines the minimal interface that the dispatcher expects from a
// message.
type BaseEvent interface {
	ID() EventID
	Context() context.Context
}

// Event defines the interface that the underlying route expects from a
// message with data.
type Event[DataT any] interface {
	BaseEvent
	Data() DataT
	Error() error
	Is(id EventID) bool
}

// NewEvent creates a new Event with the given context and beacon event.
func NewEvent[
	DataT any,
](
	ctx context.Context, id EventID, data DataT, errs ...error,
) Event[DataT] {
	return &event[DataT]{
		ctx:  ctx,
		id:   id,
		data: data,
		err:  errors.Join(errs...),
	}
}

// An event is a hard type implementation of the Event interface.
type event[DataT any] struct {
	// ctx is the context associated with the event.
	ctx context.Context
	// id is the name of the event.
	id EventID
	// event is the actual beacon event.
	data DataT
	// err is the error associated with the event.
	err error
}

// ID returns the ID of the event.
func (m *event[DataT]) ID() EventID {
	return m.id
}

// Context returns the context associated with the event.
func (m *event[DataT]) Context() context.Context {
	return m.ctx
}

// Data returns the data associated with the event.
func (m *event[DataT]) Data() DataT {
	return m.data
}

// Error returns the error associated with the event.
func (m *event[DataT]) Error() error {
	return m.err
}

// Is returns true if the event has the given type.
func (m *event[DataT]) Is(messageType EventID) bool {
	return m.id == messageType
}
