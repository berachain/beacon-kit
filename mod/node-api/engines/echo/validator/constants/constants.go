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

package constants

const (
	StateIDHead      = "head"
	StateIDGenesis   = "genesis"
	StateIDFinalized = "finalized"
	StateIDJustified = "justified"
)

const (
	ValidatorIDRegex = `^0x[0-9a-fA-F]{1,96}$`
	RootRegex        = `^0x[0-9a-fA-F]{64}$`
)

const (
	StatusPendingInitialized = "pending_initialized"
	StatusPendingQueued      = "pending_queued"
	StatusActiveOngoing      = "active_ongoing"
	StatusActiveExiting      = "active_exiting"
	StatusActiveSlashed      = "active_slashed"
	StatusExitedUnslashed    = "exited_unslashed"
	StatusExitedSlashed      = "exited_slashed"
	StatusWithdrawalPossible = "withdrawal_possible"
	StatusWithdrawalDone     = "withdrawal_done"
)
