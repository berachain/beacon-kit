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

package blocktime

import (
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// NextPayloadTimeFromSuccess calculates the
// next timestamp for an execution payload
// once parent block has successfully verified or
// has been accepted
//
// TODO: This is hood and needs to be improved.
func NextPayloadTimeFromSuccess(
	chainSpec common.ChainSpec,
	parentPayloadTime math.U64,
) uint64 {
	//#nosec:G701 // not an issue in practice.
	return max(
		uint64(time.Now().Unix())+chainSpec.TargetSecondsPerEth1Block(),
		uint64(parentPayloadTime+1),
	)
}

// NextPayloadTimeFromFailure calculates the
// next timestamp for an execution payload
// once parent block has not verified.
//
// TODO: this is hood as fuck.
func NextPayloadTimeFromFailure(parentPayloadTime math.U64) uint64 {
	//#nosec:G701 // not an issue in practice.
	return max(
		uint64(time.Now().Add(time.Second).Unix()),
		uint64(parentPayloadTime+1),
	)
}
