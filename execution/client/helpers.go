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

package client

import (
	"context"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/primitives/common"
)

// createContextWithTimeout creates a context with a timeout and returns it
// along with the cancel function.
func (s *EngineClient[
	_, _,
]) createContextWithTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	startTime := time.Now()
	dctx, cancel := context.WithTimeoutCause(
		ctx,
		s.cfg.RPCTimeout,
		engineerrors.ErrEngineAPITimeout,
	)
	s.metrics.measureNewPayloadDuration(startTime)
	return dctx, cancel
}

// processPayloadStatusResult processes the payload status result and
// returns the latest valid hash or an error.
func processPayloadStatusResult(
	result *engineprimitives.PayloadStatusV1,
) (*common.ExecutionHash, error) {
	switch result.Status {
	case engineprimitives.PayloadStatusAccepted:
		return nil, engineerrors.ErrAcceptedPayloadStatus
	case engineprimitives.PayloadStatusSyncing:
		return nil, engineerrors.ErrSyncingPayloadStatus
	case engineprimitives.PayloadStatusInvalid:
		return result.LatestValidHash, engineerrors.ErrInvalidPayloadStatus
	case engineprimitives.PayloadStatusValid:
		return result.LatestValidHash, nil
	default:
		return nil, engineerrors.ErrUnknownPayloadStatus
	}
}
