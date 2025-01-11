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

package spec

const (
	// Gwei value constants.
	DefaultMaxEffectiveBalance       = 32e9
	DefaultEjectionBalance           = 16e9
	DefaultEffectiveBalanceIncrement = 1e9

	DefaultHysteresisQuotient           = 4
	DefaultHysteresisDownwardMultiplier = 1
	DefaultHysteresisUpwardMultiplier   = 5

	// Time parameters constants.
	DefaultSlotsPerEpoch                = 32
	DefaultSlotsPerHistoricalRoot       = 8
	DefaultMinEpochsToInactivityPenalty = 4

	// Default DepositContractAddress is the default address of the pre-deployed
	// beacon deposit contract.
	DefaultDepositContractAddress = "0x4242424242424242424242424242424242424242"
)
