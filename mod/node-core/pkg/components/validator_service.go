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

package components

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// ValidatorServiceInput is the input for the validator service provider.
type ValidatorServiceInput struct {
	depinject.In
	BlobProcessor *dablob.Processor[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlockBody,
	]
	Cfg          *config.Config
	ChainSpec    primitives.ChainSpec
	LocalBuilder *payloadbuilder.PayloadBuilder[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	]
	Logger         log.Logger
	StateProcessor StateProcessor
	StorageBackend StorageBackend
	Signer         crypto.BLSSigner
	TelemetrySink  *metrics.TelemetrySink
}

// ProvideValidatorService is a depinject provider for the validator service.
func ProvideValidatorService(
	in ValidatorServiceInput,
) *validator.Service[
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	BeaconState,
	*datypes.BlobSidecars,
	*depositdb.KVStore[*types.Deposit],
	*types.ForkData,
] {
	// Build the builder service.
	return validator.NewService[
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*depositdb.KVStore[*types.Deposit],
		*types.ForkData,
	](
		&in.Cfg.Validator,
		in.Logger.With("service", "validator"),
		in.ChainSpec,
		in.StorageBackend,
		in.BlobProcessor,
		in.StateProcessor,
		in.Signer,
		dablob.NewSidecarFactory[
			*types.BeaconBlock,
			*types.BeaconBlockBody,
		](
			in.ChainSpec,
			types.KZGPositionDeneb,
			in.TelemetrySink,
		),
		in.LocalBuilder,
		[]validator.PayloadBuilder[BeaconState, *types.ExecutionPayload]{
			in.LocalBuilder,
		},
		in.TelemetrySink,
	)
}
