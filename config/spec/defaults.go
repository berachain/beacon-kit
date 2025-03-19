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

package spec

// NOTE: Most of these default values are taken from ETH2.0 spec.
// Some values (mentioned below) are modified to better suit Berachain's system.

//nolint:unused // Keeping values here for reference.
const (
	// Gwei value constants.
	defaultMaxEffectiveBalance       = 32e9
	defaultEjectionBalance           = 16e9
	defaultEffectiveBalanceIncrement = 1e9

	defaultHysteresisQuotient           = 4
	defaultHysteresisDownwardMultiplier = 1
	defaultHysteresisUpwardMultiplier   = 5

	// Time parameters constants.
	defaultSlotsPerEpoch                = 32
	defaultSlotsPerHistoricalRoot       = 8
	defaultMinEpochsToInactivityPenalty = 4

	// Signature domains.
	defaultDomainTypeProposer          = 0
	defaultDomainTypeAttester          = 1
	defaultDomainTypeRandao            = 2
	defaultDomainTypeDeposit           = 3
	defaultDomainTypeVoluntaryExit     = 4
	defaultDomainTypeSelectionProof    = 5
	defaultDomainTypeAggregateAndProof = 6
	defaultDomainTypeApplicationMask   = 16777216 // "0x00000001" in little endian

	// Eth1-related values.
	defaultDepositContractAddress    = "0x4242424242424242424242424242424242424242" // Berachain specific.
	defaultMaxDepositsPerBlock       = 16
	defaultDepositEth1ChainID        = 1
	defaultEth1FollowDistance        = 1 // Berachain specific.
	defaultTargetSecondsPerEth1Block = 2 // Berachain specific.

	// Fork-related values.
	defaultDeneb1ForkEpoch  = 9999999999999998 // Set as a future epoch as not yet determined.
	defaultElectraForkEpoch = 9999999999999999 // Set as a future epoch as not yet determined.

	// State list length constants.
	defaultEpochsPerHistoricalVector = 8
	defaultEpochsPerSlashingsVector  = 8
	defaultHistoricalRootsLimit      = 8
	defaultValidatorRegistryLimit    = 1099511627776

	// Slashing.
	defaultProportionalSlashingMultiplier = 1

	// Capella values.
	defaultMaxWithdrawalsPerPayload         = 16
	defaultMaxValidatorsPerWithdrawalsSweep = 1 << 14

	// Deneb values.
	defaultMinEpochsForBlobsSidecarsRequest = 4096
	defaultMaxBlobCommitmentsPerBlock       = 4096
	defaultMaxBlobsPerBlock                 = 6
	defaultFieldElementsPerBlob             = 4096
	defaultBytesPerBlob                     = 131072
	defaultKZGCommitmentInclusionProofDepth = 17

	// Berachain values.
	defaultValidatorSetCap      = 256
	defaultEVMInflationAddress  = "0x0000000000000000000000000000000000000000"
	defaultEVMInflationPerBlock = 0
)
