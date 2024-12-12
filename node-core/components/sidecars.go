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
	"github.com/berachain/beacon-kit/consensus-types/types"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
)

type SidecarFactoryInput struct {
	depinject.In
	ChainSpec     common.ChainSpec
	TelemetrySink *metrics.TelemetrySink
}

func ProvideSidecarFactory[
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, *AttestationData, DepositT,
		ExecutionPayloadT, *SlashingInfo,
	],
	DepositT any,
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader[ExecutionPayloadHeaderT],
	WithdrawalsT Withdrawals[WithdrawalT],
	WithdrawalT Withdrawal[WithdrawalT],
](in SidecarFactoryInput) *dablob.SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT,
] {
	return dablob.NewSidecarFactory[
		BeaconBlockT,
		BeaconBlockBodyT,
	](
		in.ChainSpec,
		types.KZGPositionDeneb,
		in.TelemetrySink,
	)
}
