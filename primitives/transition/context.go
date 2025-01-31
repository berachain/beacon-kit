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

package transition

import (
	"context"

	"github.com/berachain/beacon-kit/primitives/math"
)

// Context is the context for the state transition.
type Context struct {
	// ConsensusCtx is the context passed by CometBFT callbacks
	// We pass it down to be able to cancel processing (although
	// currently CometBFT context is set to TODO)
	ConsensusCtx context.Context
	// MeterGas controls whether gas data related to the execution
	// layer payload should be meter or not. We currently meter only
	// finalized blocks.
	MeterGas bool
	// OptimisticEngine indicates whether to optimistically assume
	// the execution client has the correct state certain errors
	// are returned by the execution engine.
	OptimisticEngine bool
	// VerifyPayload indicates whether to call NewPayload on the
	// execution client. This can be done when the node is not
	// syncing, and the payload is already known to the execution client.
	VerifyPayload bool
	// ValidateRandao indicates whether to validate the Randao mix.
	ValidateRandao bool
	// ValidateResult indicates whether to validate the result of
	// the state transition.
	ValidateResult bool
	// Address of current block proposer
	ProposerAddress []byte
	// ConsensusTime returns the timestamp of current consensus request.
	// It is used to build next payload and to validate currentpayload.
	ConsensusTime math.U64
}
