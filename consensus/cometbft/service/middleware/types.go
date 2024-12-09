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

package middleware

import (
	"time"

	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/transition"
)

// BeaconBlock is an interface for accessing the beacon block.
type BeaconBlock[BeaconBlockT any, BeaconBlockHeaderT any] interface {
	constraints.SSZMarshallable
	constraints.Nillable
	constraints.Empty[BeaconBlockT]
	NewFromSSZ([]byte, uint32) (BeaconBlockT, error)

	GetHeader() BeaconBlockHeaderT
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the given time.
	MeasureSince(key string, start time.Time, args ...string)
}

type BlobSidecars[T any] interface {
	constraints.SSZMarshallable
	constraints.Empty[T]
}

type validatorUpdates = transition.ValidatorUpdates
