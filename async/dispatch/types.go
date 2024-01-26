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

package dispatch

import (
	"time"

	"github.com/itsdevbear/bolaris/async/dispatch/queues"
)

// Queue represents a queue of work items to be executed. It's interface is inspired by
// Apple's Grand Central Dispatch (GCD) API.
// https://developer.apple.com/documentation/dispatch/dispatchqueue
type Queue interface {
	Async(queues.WorkItem)
	AsyncAfter(time.Duration, queues.WorkItem)
	Sync(queues.WorkItem)
	AsyncAndWait(queues.WorkItem)
}

// Event represents actions that occur during consensus. Listeners can
// register callbacks with event handlers for specific event types.
type Event interface {
	Type() string
	Source() any
	Value() any
}
