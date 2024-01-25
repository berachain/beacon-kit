// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package notify

import "github.com/itsdevbear/bolaris/async/dispatch"

// GrandCentralDispatch is an interface that wraps the basic GetQueue method.
// It is used to retrieve a dispatch queue by its ID.
type GrandCentralDispatch interface {
	// GetQueue returns a queue with the provided ID.
	GetQueue(id string) dispatch.Queue
}

// EventHandler is the interface that wraps the basic Handle method.
type EventHandler interface {
	HandleNotification(event interface{})
}

// eventHandlerQueuePair is a struct that holds an event handler and a queue ID.
type eventHandlerQueuePair struct {
	// handler is an object that implements the EventHandler interface.
	handler EventHandler
	// queueID is a string that identifies a dispatch queue.
	queueID string
}
