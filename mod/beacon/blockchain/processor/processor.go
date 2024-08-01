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

package processor

import (
	"sync"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// Processor is the blockchain service.
type Processor[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BeaconBlockHeaderT BeaconBlockHeader,
	BeaconStateT ReadOnlyBeaconState[
		BeaconStateT, BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	],
	BlobSidecarsT BlobSidecars,
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	PayloadAttributesT interface {
		IsNil() bool
		Version() uint32
		GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress
	},
	WithdrawalT any,
] struct {
	// sb represents the backend storage for beacon states and associated
	// sidecars.
	sb StorageBackend[
		AvailabilityStoreT,
		BeaconBlockBodyT,
		BeaconStateT,
		BlobSidecarsT,
	]
	// logger is used for logging messages in the service.
	logger log.Logger[any]
	// cs holds the chain specifications.
	cs common.ChainSpec
	// ee is the execution engine responsible for processing execution payloads.

	ee ExecutionEngine[PayloadAttributesT]
	// lb is a local builder for constructing new beacon states.
	lb LocalBuilder[BeaconStateT]
	// sp is the state processor for beacon blocks and states.
	sp StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	]
	// metrics is the metrics for the service.
	metrics *chainMetrics
	// optimisticPayloadBuilds is a flag used when the optimistic payload
	// builder is enabled.
	optimisticPayloadBuilds bool
	// forceStartupSyncOnce is used to force a sync of the startup head.
	forceStartupSyncOnce *sync.Once
}

// New creates a new validator service.
func New[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, ExecutionPayloadT],
	BeaconBlockBodyT BeaconBlockBody[ExecutionPayloadT],
	BeaconBlockHeaderT BeaconBlockHeader,
	BeaconStateT ReadOnlyBeaconState[
		BeaconStateT, BeaconBlockHeaderT,
		ExecutionPayloadHeaderT,
	],
	BlobSidecarsT BlobSidecars,
	DepositT any,
	ExecutionPayloadT ExecutionPayload,
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	GenesisT Genesis[DepositT, ExecutionPayloadHeaderT],
	PayloadAttributesT interface {
		IsNil() bool
		Version() uint32
		GetSuggestedFeeRecipient() gethprimitives.ExecutionAddress
	},
	WithdrawalT any,
](
	sb StorageBackend[
		AvailabilityStoreT,
		BeaconBlockBodyT,
		BeaconStateT,
		BlobSidecarsT,
	],
	logger log.Logger[any],
	cs common.ChainSpec,
	ee ExecutionEngine[PayloadAttributesT],
	lb LocalBuilder[BeaconStateT],
	sp StateProcessor[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
		*transition.Context,
		DepositT,
		ExecutionPayloadHeaderT,
	],
	ts TelemetrySink,
	optimisticPayloadBuilds bool,
) *Processor[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, DepositT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, GenesisT, PayloadAttributesT, WithdrawalT,
] {
	return &Processor[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobSidecarsT, DepositT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, GenesisT, PayloadAttributesT, WithdrawalT,
	]{
		sb:                      sb,
		logger:                  logger,
		cs:                      cs,
		ee:                      ee,
		lb:                      lb,
		sp:                      sp,
		metrics:                 newChainMetrics(ts),
		optimisticPayloadBuilds: optimisticPayloadBuilds,
		forceStartupSyncOnce:    new(sync.Once),
	}
}
