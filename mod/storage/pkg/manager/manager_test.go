// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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

package manager_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/storage/pkg/interfaces/mocks"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/stretchr/testify/require"
)

func TestDBManager_Start(t *testing.T) {
	mockPrunable := new(mocks.Prunable)
	logger := log.NewNopLogger()
	p1 := pruner.NewPruner(logger, mockPrunable, "pruner1")
	p2 := pruner.NewPruner(logger, mockPrunable, "pruner2")

	m, err := manager.NewDBManager(
		manager.WithPruner(p1),
		manager.WithPruner(p2),
		manager.WithLogger(logger),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = m.Start(ctx)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	mockPrunable.AssertNotCalled(t, "Prune")
}

func TestDBManager_NotifyAll(t *testing.T) {
	mockPruner := new(mocks.Prunable)
	logger := log.NewNopLogger()
	p1 := pruner.NewPruner(logger, mockPruner, "pruner1")
	p2 := pruner.NewPruner(logger, mockPruner, "pruner2")

	// set up the expectation for the Prune method
	mockPruner.On("Prune", uint64(123)).Return(nil).Twice()

	m, err := manager.NewDBManager(
		manager.WithPruner(p1),
		manager.WithPruner(p2),
		manager.WithLogger(logger),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = m.Start(ctx)
	require.NoError(t, err)
	m.NotifyAll(123)

	time.Sleep(100 * time.Millisecond)

	// Assert that both pruners received the Notify call
	mockPruner.AssertCalled(t, "Prune", uint64(123))
	mockPruner.AssertNumberOfCalls(t, "Prune", 2)
}

func TestDBManager_Notify(t *testing.T) {
	mockPruner := new(mocks.Prunable)
	logger := log.NewNopLogger()
	p1 := pruner.NewPruner(logger, mockPruner, "pruner1")
	p2 := pruner.NewPruner(logger, mockPruner, "pruner2")

	// set up the expectation for the Prune method
	mockPruner.On("Prune", uint64(123)).Return(nil).Twice()

	m, err := manager.NewDBManager(
		manager.WithPruner(p1),
		manager.WithPruner(p2),
		manager.WithLogger(logger),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = m.Start(ctx)
	require.NoError(t, err)
	m.Notify("pruner1", 123)

	time.Sleep(100 * time.Millisecond)

	// Assert that only pruner1 received the Notify call
	mockPruner.AssertCalled(t, "Prune", uint64(123))
	mockPruner.AssertNumberOfCalls(t, "Prune", 1)
}
