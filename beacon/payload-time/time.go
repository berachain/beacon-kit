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

package payloadtime

import (
	"errors"
	"fmt"

	"github.com/berachain/beacon-kit/primitives/pkg/math"
)

// ErrTooFarInTheFuture is returned when the payload timestamp
// in a block exceeds the time bound.
var ErrTooFarInTheFuture = errors.New("timestamp too far in the future")

func Verify(
	consensusTime,
	parentPayloadTimestamp,
	payloadTimestamp math.U64,
) error {
	bound := max(
		consensusTime+1,
		parentPayloadTimestamp+1,
	)
	if payloadTimestamp > bound {
		return fmt.Errorf(
			"%w: timestamp bound: %d, got: %d",
			ErrTooFarInTheFuture,
			bound, payloadTimestamp,
		)
	}
	return nil
}

func Next(
	consensusTime,
	parentPayloadTimestamp math.U64,
	buildOptimistically bool,
) math.U64 {
	delta := math.U64(0)
	if buildOptimistically {
		// we're building a payload to be included into next block.
		// We estimate it to be included next second. If this estimate
		// turns out wrong (cause consensus block are finalized faster or
		// slower than consensusTime+1 sec), we're still fine as long as
		// Verify pass which should always to since:
		// Next.consensusTime <= Verify.consensusTime
		delta = 1
	}
	return max(
		consensusTime+delta,
		parentPayloadTimestamp+1,
	)
}
