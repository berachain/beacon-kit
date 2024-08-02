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

package app

import (
	"context"

	consensusengine "github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
)

type App[
	StorageBackendT any,
	StateProcessorT any,
] struct {
	Logger log.Logger[any]

	// The consensus engine client is responsible
	// for communicating with the consensus engine
	// for the chain.
	consensusengine.Client
	// The execution engine client is responsible
	// for communicating with the execution engine
	// for the chain.
	// ExecutionClient executionengine.Engine
	// The backend is the central data access layer for
	// the application.
	backend StorageBackendT
	// The state processor is the component that is
	// responsible for transitioning the state of the
	// chain.
	stateProcessor StateProcessorT
	// Services contained within the service registry
	// are defined as pieces of app-specific tooling
	// that are used to support the core functionality
	// of the application.
	services *service.Registry
}

func (a *App[_, _]) Start(ctx context.Context) error {
	return a.services.StartAll(ctx)
}
