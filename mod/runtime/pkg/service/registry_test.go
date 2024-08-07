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

package service_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/log/pkg/noop"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegistry_StartAll(t *testing.T) {
	logger := noop.NewLogger[any]()
	dispatcher := &mocks.Dispatcher{}
	registry := service.NewRegistry(dispatcher, service.WithLogger(logger))

	service1 := &mocks.Basic{}
	service1.On("Start", mock.Anything).Return(nil).Once()
	service1.On("Name").Return("Service1")

	service2 := &mocks.Basic{}
	service2.On("Start", mock.Anything).Return(nil).Once()
	service2.On("Name").Return("Service2")

	err := registry.RegisterService(service1)
	if err != nil {
		t.Fatalf("Failed to register Service1: %v", err)
	}

	err = registry.RegisterService(service2)
	if err != nil {
		t.Fatalf("Failed to register Service2: %v", err)
	}

	require.NoError(t, registry.StartAll(context.Background()))
	time.Sleep(25 * time.Millisecond)

	service1.AssertCalled(t, "Start", mock.Anything)
	service2.AssertCalled(t, "Start", mock.Anything)
}

func TestRegistry_FetchService(t *testing.T) {
	logger := noop.NewLogger[any]()
	dispatcher := &mocks.Dispatcher{}
	registry := service.NewRegistry(dispatcher, service.WithLogger(logger))

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
