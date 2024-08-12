package types

import (
	"fmt"
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
		return zero, fmt.Errorf("future timed out")
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
