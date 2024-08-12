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
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/cli/pkg/flags"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
	"github.com/spf13/cast"
)

// BlobProofVerifierInput is the input for the
// dep inject framework.
type BlobProofVerifierInput struct {
	depinject.In
	AppOpts          servertypes.AppOptions
	JSONTrustedSetup *gokzg4844.JSONTrustedSetup
}

// ProvideBlobProofVerifier is a function that provides the module to the
// application.
func ProvideBlobProofVerifier(
	in BlobProofVerifierInput,
) (kzg.BlobProofVerifier, error) {
	return kzg.NewBlobProofVerifier(
		cast.ToString(in.AppOpts.Get(flags.KZGImplementation)),
		in.JSONTrustedSetup,
	)
}

// BlobVerifierInput is the input for the BlobVerifier.
type BlobVerifierInput[
	BlobProofVerifierT kzg.BlobProofVerifier,
] struct {
	depinject.In
	BlobProofVerifier BlobProofVerifierT
	TelemetrySink     *metrics.TelemetrySink
}

// ProvideBlobVerifier is a function that provides the BlobVerifier to the
// depinject framework.
func ProvideBlobVerifier[
	BlobProofVerifierT kzg.BlobProofVerifier,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BlobSidecarT BlobSidecar[BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecars[BlobSidecarT],
](in BlobVerifierInput[BlobProofVerifierT]) *dablob.Verifier[
	BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT,
] {
	return dablob.NewVerifier[
		BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT,
	](in.BlobProofVerifier, in.TelemetrySink)
}

// BlobProcessorIn is the input for the BlobProcessor.
type BlobProcessorIn[
	BlobProofVerifierT kzg.BlobProofVerifier,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BlobSidecarT BlobSidecar[BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecars[BlobSidecarT],
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In

	BlobVerifier *dablob.Verifier[
		BeaconBlockHeaderT, BlobSidecarT, BlobSidecarsT,
	]
	ChainSpec     common.ChainSpec
	Logger        LoggerT
	TelemetrySink *metrics.TelemetrySink
}

// ProvideBlobProcessor is a function that provides the BlobProcessor to the
// depinject framework.
func ProvideBlobProcessor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BlobProofVerifierT kzg.BlobProofVerifier,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BlobSidecarT BlobSidecar[BeaconBlockHeaderT],
	BlobSidecarsT BlobSidecars[BlobSidecarT],
	LoggerT log.AdvancedLogger[any, LoggerT],
](in BlobProcessorIn[
	BlobProofVerifierT, BeaconBlockHeaderT, BlobSidecarT,
	BlobSidecarsT, LoggerT,
]) *dablob.Processor[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BlobSidecarT, BlobSidecarsT,
] {
	return dablob.NewProcessor[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BlobSidecarT, BlobSidecarsT,
	](
		in.Logger.With("service", "blob-processor"),
		in.ChainSpec,
		in.BlobVerifier,
		types.BlockBodyKZGOffset,
		in.TelemetrySink,
	)
}

// DAServiceIn is the input for the BlobService.
type DAServiceIn[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BlobProcessorT BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	],
	BlobSidecarsT any,
	ExecutionPayloadT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In

	AvailabilityStore AvailabilityStoreT
	SidecarsBroker    *broker.Broker[*asynctypes.Event[BlobSidecarsT]]
	BlobProcessor     BlobProcessorT
	Logger            LoggerT
}

// ProvideDAService is a function that provides the BlobService to the
// depinject framework.
func ProvideDAService[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT any,
	BlobProcessorT BlobProcessor[
		AvailabilityStoreT, BeaconBlockBodyT,
		BlobSidecarsT, ExecutionPayloadT,
	],
	BlobSidecarT any,
	BlobSidecarsT BlobSidecars[BlobSidecarT],
	ExecutionPayloadT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
](in DAServiceIn[
	AvailabilityStoreT, BeaconBlockBodyT, BlobProcessorT,
	BlobSidecarsT, ExecutionPayloadT, LoggerT,
]) *da.Service[
	AvailabilityStoreT, BeaconBlockBodyT, BlobProcessorT,
	BlobSidecarsT, ExecutionPayloadT,
] {
	return da.NewService[
		AvailabilityStoreT, BeaconBlockBodyT, BlobProcessorT,
		BlobSidecarsT, ExecutionPayloadT,
	](
		in.AvailabilityStore,
		in.BlobProcessor,
		in.SidecarsBroker,
		in.Logger.With("service", "da"),
	)
}
