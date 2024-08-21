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
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
)

// ValidatorServiceInput is the input for the validator service provider.
type ValidatorServiceInput[
	AvailabilityStoreT any,
	BeaconBlockT any,
	BeaconStateT any,
	BlobSidecarsT any,
	DepositT any,
	LoggerT any,
	StorageBackendT any,
] struct {
	depinject.In
	Cfg            *config.Config
	ChainSpec      common.ChainSpec
	Dispatcher     *Dispatcher
	LocalBuilder   LocalBuilder[BeaconStateT, *ExecutionPayload]
	Logger         LoggerT
	StateProcessor StateProcessor[
		BeaconBlockT, BeaconStateT, *Context, DepositT, *ExecutionPayloadHeader,
	]
	StorageBackend StorageBackendT
	Signer         crypto.BLSSigner
	SidecarFactory SidecarFactory[BeaconBlockT, BlobSidecarsT]
	TelemetrySink  *metrics.TelemetrySink
}

// ProvideValidatorService is a depinject provider for the validator service.
func ProvideValidatorService[
	AvailabilityStoreT any,
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		*Eth1Data, *ExecutionPayload, *SlashingInfo,
	],
	BeaconBlockHeaderT any,
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, BeaconStateMarshallableT,
		*Eth1Data, *ExecutionPayloadHeader, *Fork, KVStoreT,
		*Validator, Validators, *Withdrawal,
	],
	BeaconStateMarshallableT any,
	BeaconBlockStoreT any,
	BlobSidecarsT any,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	KVStoreT any,
	LoggerT log.AdvancedLogger[any, LoggerT],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BeaconBlockStoreT, DepositStoreT,
	],
](
	in ValidatorServiceInput[
		AvailabilityStoreT, BeaconBlockT, BeaconStateT,
		BlobSidecarsT, DepositT, LoggerT, StorageBackendT,
	],
) (*validator.Service[
	*AttestationData, BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT, DepositT, DepositStoreT,
	*Eth1Data, *ExecutionPayload, *ExecutionPayloadHeader,
	*ForkData, *SlashingInfo, *SlotData,
], error) {
	// Build the builder service.
	return validator.NewService[
		*AttestationData,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconStateT,
		BlobSidecarsT,
		DepositT,
		DepositStoreT,
		*Eth1Data,
		*ExecutionPayload,
		*ExecutionPayloadHeader,
		*ForkData,
		*SlashingInfo,
		*SlotData,
	](
		&in.Cfg.Validator,
		in.Logger.With("service", "validator"),
		in.ChainSpec,
		in.StorageBackend,
		in.StateProcessor,
		in.Signer,
		in.SidecarFactory,
		in.LocalBuilder,
		[]validator.PayloadBuilder[BeaconStateT, *ExecutionPayload]{
			in.LocalBuilder,
		},
		in.TelemetrySink,
		in.Dispatcher,
	), nil
}
