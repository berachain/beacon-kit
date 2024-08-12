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

// MessageID represents the type of a message.
type MessageID string

// EventID is a type alias for a MessageID.
type EventID = MessageID

// BaseMessage defines the minimal interface that the dispatcher expects from a
// message.
type BaseMessage interface {
	ID() MessageID
	Context() context.Context
}

// DataMessage defines the interface that the underlying route expects from a
// message with data.
type Message[DataT any] interface {
	BaseMessage
	Data() DataT
	Error() error
	Is(messageType MessageID) bool
}

// NewEvent creates a new Event with the given context and beacon event.
func NewMessage[
	DataT any,
](
	ctx context.Context, messageType MessageID, data DataT, errs ...error,
) Message[DataT] {
	return &message[DataT]{
		ctx:  ctx,
		id:   messageType,
		data: data,
		err:  errors.Join(errs...),
	}
}

// Event acts as a type alias for a Message that is meant to be broadcasted
// to all subscribers.
type Event[DataT any] struct{ Message[DataT] }

// NewEvent creates a new Event with the given context and beacon event.
func NewEvent[
	DataT any,
](
	ctx context.Context, messageType EventID, data DataT, errs ...error,
) *Event[DataT] {
	return &Event[DataT]{
		Message: NewMessage(ctx, messageType, data, errs...),
	}
}

// A Message is a hard type implementation of the Message and Event interfaces.
type message[DataT any] struct {
	// ctx is the context associated with the event.
	ctx context.Context
	// id is the name of the event.
	id MessageID
	// event is the actual beacon event.
	data DataT
	// err is the error associated with the event.
	err error
}

// ID returns the ID of the event.
func (m *message[DataT]) ID() MessageID {
	return m.id
}

// Context returns the context associated with the event.
func (m *message[DataT]) Context() context.Context {
	return m.ctx
}

// Data returns the data associated with the event.
func (m *message[DataT]) Data() DataT {
	return m.data
}

// Error returns the error associated with the event.
func (m *message[DataT]) Error() error {
	return m.err
}

// Is returns true if the event has the given type.
func (m *message[DataT]) Is(messageType MessageID) bool {
	return m.id == messageType
}
