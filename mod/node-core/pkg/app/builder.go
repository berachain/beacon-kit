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
	"github.com/berachain/beacon-kit/mod/config"
	consensusengine "github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/depinject"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components"
)

type Builder[
	StorageBackendT any,
	StateProcessorT any,
] struct {
	app *App[StorageBackendT, StateProcessorT]

	// components is a slice of components to be added to the app
	// These components will be depinjected into the app and added
	// to the app at runtime.
	components []any
}

func NewBuilder[
	StorageBackendT any,
	StateProcessorT any,
]() *Builder[StorageBackendT, StateProcessorT] {
	return &Builder[StorageBackendT, StateProcessorT]{
		app: &App[StorageBackendT, StateProcessorT]{},
	}
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) Build(
	logger *phuslu.Logger,
	appOpts *components.AppOptions,
	config *config.Config,
) (*App[StorageBackendT, StateProcessorT], error) {
	var err error
	// depinject components into the app.
	container := depinject.NewContainer()
	if err = container.Supply(
		logger,
		appOpts,
		config,
	); err != nil {
		return nil, err
	}
	if err = container.Provide(b.components...); err != nil {
		return nil, err
	}

	var client *components.RuntimeApp
	// Resolve dependencies and construct the app.
	if err = container.Inject(
		&b.app.backend,
		&b.app.stateProcessor,
		&b.app.services,
		&client,
	); err != nil {
		return nil, err
	}
	b.app.Logger = logger
	b.app.Client = client

	return b.app, nil
}

// BUILDER OPTIONS

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithComponents(
	components ...any,
) *Builder[StorageBackendT, StateProcessorT] {
	b.components = append(b.components, components...)
	return b
}

// For these methods, the paradigm will be to pass a pointer to the eventually
// constructed component. This allows the app to force the inclusion of the
// backend
// and state processor in the app without being strict on the actual components
// included, think of it as a minimal set of required components.

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithConsensusClient(
	consensusClient consensusengine.Client,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.Client = consensusClient
	return b
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithStorageBackend(
	storageBackend StorageBackendT,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.backend = storageBackend
	return b
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) WithStateProcessor(
	stateProcessor StateProcessorT,
) *Builder[StorageBackendT, StateProcessorT] {
	b.app.stateProcessor = stateProcessor
	return b
}

func (b *Builder[
	StorageBackendT, StateProcessorT,
]) App() *App[StorageBackendT, StateProcessorT] {
	return b.app
}
