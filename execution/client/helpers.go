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

package client

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/primitives/common"
)

// createContextWithTimeout creates a context with a timeout and returns it
// along with the cancel function.
func (s *EngineClient) createContextWithTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	dctx, cancel := context.WithTimeoutCause(
		ctx,
		s.cfg.RPCTimeout,
		engineerrors.ErrEngineAPITimeout,
	)
	return dctx, cancel
}

// processPayloadStatusResult processes the payload status result and
// returns the latest valid hash or an error.
func processPayloadStatusResult(
	result *engineprimitives.PayloadStatusV1,
) (*common.ExecutionHash, error) {
	switch result.Status {
	case engineprimitives.PayloadStatusAccepted:
		// NewPayload -- Returned then the payload does not extend the canonical chain. Already enforced on CL in state_processor_payload.
		// FCU -- NEVER returned
		// CL TODOs:
		// - Remove from FCU check because it is never returned (done in shota's PR)
		// - Distinguish metric between accepted/syncing
		// - Remove error-handling from state_processor.go because
		//    - in process we DO want to return error (as is)
		//    - in finalize we do NOT want to return error (as is)
		return nil, engineerrors.ErrAcceptedPayloadStatus
	case engineprimitives.PayloadStatusSyncing:
		// NewPayload --
		//	- EL does not have the parent block. The EL does NOT know if it is valid or invalid
		//  - EL has the parent block but is in Snap Sync (Geth).
		//  - CL TODOs:
		//     - Distinguish metric between accepted/syncing.
		//     - Keep as is for now - return Error in Process and no Error in Finalize.
		//     - Consider blocking in Finalize to catch up EL when syncing.
		// FCU --
		return nil, engineerrors.ErrSyncingPayloadStatus
	case engineprimitives.PayloadStatusInvalid:
		// NewPayload --
		//	- if the block is part of an invalid chain. We don't expect this in finalize, ALWAYS return err.
		//  - if the block's timestamp is <= the parent block's timestamp. We don't expect this in finalize, ALWAYS return err.
		//  - if invalid state transition. We don't expect this in finalize, ALWAYS return err.
		//  - CL TODOs: Keep as is.
		// FCU --
		return result.LatestValidHash, engineerrors.ErrInvalidPayloadStatus
	case engineprimitives.PayloadStatusValid:
		// NewPayload --
		//  - EL returns in happy case.
		//  - EL already has the payload, i.e. duplicate payload (should not happen technically).
		// 	- CL TODOs: Keep as is.
		// FCU --
		return result.LatestValidHash, nil
	default:
		// NewPayload -- Not expected ever, ALWAYS return err.
		// FCU -- Not expected ever, ALWAYS return err.
		// CL TODOs:
		//  - In FinalizeBlock, force return error to ensure EL is in sync.
		return nil, engineerrors.ErrUnknownPayloadStatus
	}
}
