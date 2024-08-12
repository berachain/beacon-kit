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
	"sync"
	"time"
)

// FutureI is the interface that the future implements.
type FutureI[T any] interface {
	// Resolve returns the result of the future, blocking until it's available
	// or the context is done.
	Resolve() (result T, err error)
	// ResolveWithTimeout returns the result of the future, blocking until it's
	// available or the timeout is reached.
	ResolveWithTimeout(timeout time.Duration) (result T, err error)
	// SetResult is called by the router to set the result of the future.
	SetResult(result T, err error)
	// IsDone returns true if the future has completed (successfully or with an
	// error).
	IsDone() bool
}

// Future represents a value that will be available at some point in the future.
type Future[T any] struct {
	result T
	err    error
	done   chan struct{}
	once   sync.Once
}

// NewFuture creates a new Future and starts the given function in a goroutine.
func NewFuture[T any]() *Future[T] {
	f := &Future[T]{
		done: make(chan struct{}),
	}
	return f
}

// Resolve returns the result of the future, blocking until it's available or
// the context is done.
func (f *Future[T]) Resolve() (T, error) {
	<-f.done
	return f.result, f.err
}

// SetResult sets the result of the future.
func (f *Future[T]) SetResult(result T, err error) {
	f.result = result
	f.err = err
	f.once.Do(func() { close(f.done) })
}

// GetWithTimeout returns the result of the future, blocking until it's
// available or the timeout is reached.
func (f *Future[T]) ResolveWithTimeout(timeout time.Duration) (T, error) {
	select {
	case <-f.done:
		return f.result, f.err
	case <-time.After(timeout):
		var zero T
		return zero, ErrTimeout
	}
}

// IsDone returns true if the future has completed (successfully or with an
// error).
func (f *Future[T]) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}
