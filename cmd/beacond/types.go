// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package main

import (
	"github.com/berachain/beacon-kit/beacon/validator"
	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	dablob "github.com/berachain/beacon-kit/da/blob"
	engineclient "github.com/berachain/beacon-kit/execution/client"
	execution "github.com/berachain/beacon-kit/execution/engine"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-api/engines/echo"
	"github.com/berachain/beacon-kit/node-api/server"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/node-core/services/version"
	"github.com/berachain/beacon-kit/payload/attributes"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/transition"
	"github.com/berachain/beacon-kit/state-transition/core"
	"github.com/berachain/beacon-kit/storage/filedb"
)

/* -------------------------------------------------------------------------- */
/*                                  Services                                  */
/* -------------------------------------------------------------------------- */

type (
	// AttributesFactory is a type alias for the attributes factory.
	AttributesFactory = attributes.Factory

	// BlobProcessor is a type alias for the blob processor.
	BlobProcessor = dablob.Processor

	// CometBFTService is a type alias for the CometBFT service.
	CometBFTService = cometbft.Service[*Logger]

	// EngineClient is a type alias for the engine client.
	EngineClient = engineclient.EngineClient

	// EngineClient is a type alias for the engine client.
	ExecutionEngine = execution.Engine

	// IndexDB is a type alias for the range DB.
	IndexDB = filedb.RangeDB

	// LocalBuilder is a type alias for the local builder.
	LocalBuilder = payloadbuilder.PayloadBuilder

	// NodeAPIEngine is a type alias for the node API engine.
	NodeAPIEngine = echo.Engine

	// NodeAPIServer is a type alias for the node API server.
	NodeAPIServer = server.Server

	// ReportingService is a type alias for the reporting service.
	ReportingService = version.ReportingService

	// SidecarFactory is a type alias for the sidecar factory.
	SidecarFactory = dablob.SidecarFactory

	// StateProcessor is the type alias for the state processor interface.
	StateProcessor = core.StateProcessor[*Context]

	// StorageBackend is the type alias for the storage backend interface.
	StorageBackend = storage.Backend

	// ValidatorService is a type alias for the validator service.
	ValidatorService = validator.Service
)

/* -------------------------------------------------------------------------- */
/*                                    Types                                   */
/* -------------------------------------------------------------------------- */

type (
	// Context is a type alias for the transition context.
	Context = transition.Context

	// Logger is a type alias for the logger.
	Logger = phuslu.Logger

	// LoggerConfig is a type alias for the logger config.
	LoggerConfig = phuslu.Config
)
