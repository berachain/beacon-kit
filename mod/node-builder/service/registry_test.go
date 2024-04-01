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

package service_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/berachain/beacon-kit/mod/node-builder/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRegistry_StartAll(t *testing.T) {
	logger := log.NewLogger(os.Stdout)
	registry := service.NewRegistry(service.WithLogger(logger))

	service1 := &mocks.Basic{}
	service1.On("Start", mock.Anything).Return().Once()
	service1.On("Name").Return("Service1")

	service2 := &mocks.Basic{}
	service2.On("Start", mock.Anything).Return().Once()
	service2.On("Name").Return("Service2")

	err := registry.RegisterService(service1)
	if err != nil {
		t.Fatalf("Failed to register Service1: %v", err)
	}

	err = registry.RegisterService(service2)
	if err != nil {
		t.Fatalf("Failed to register Service2: %v", err)
	}

	registry.StartAll(context.Background())
	time.Sleep(25 * time.Millisecond)

	service1.AssertCalled(t, "Start", mock.Anything)
	service2.AssertCalled(t, "Start", mock.Anything)
}

func TestRegistry_Statuses(t *testing.T) {
	logger := log.NewNopLogger()
	registry := service.NewRegistry(service.WithLogger(logger))

	service1 := &mocks.Basic{}
	service1.On("Status").Return(nil).Twice()
	service1.On("Name").Return("Service1")

	service2 := &mocks.Basic{}
	service2.On("Status").Return(nil).Twice()
	service2.On("Name").Return("Service2")

	err := registry.RegisterService(service1)
	if err != nil {
		t.Fatalf("Failed to register Service1: %v", err)
	}

	err = registry.RegisterService(service2)
	if err != nil {
		t.Fatalf("Failed to register Service2: %v", err)
	}

	statuses := registry.Statuses()
	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}

	for _, err := range statuses {
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	service1.On("Start", mock.Anything).Return().Once()
	service2.On("Start", mock.Anything).Return().Once()
	registry.StartAll(context.Background())

	statuses = registry.Statuses()

	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}

	service1.AssertNumberOfCalls(t, "Status", 2)
	service2.AssertNumberOfCalls(t, "Status", 2)
}

func TestRegistry_FetchService(t *testing.T) {
	logger := log.NewNopLogger()
	registry := service.NewRegistry(service.WithLogger(logger))

	service1 := new(mocks.Basic)
	service1.On("Name").Return("Service1")
	if err := registry.RegisterService(service1); err != nil {
		t.Fatalf("Failed to register Service1: %v", err)
	}

	var fetchedService *mocks.Basic
	if err := registry.FetchService(&fetchedService); err != nil {
		t.Fatalf("Failed to fetch service: %v", err)
	}

	if reflect.TypeOf(fetchedService) != reflect.TypeOf(service1) {
		t.Errorf("Fetched service type mismatch")
	}
}
