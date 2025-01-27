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

// NOTE: Most of these default values are taken from ETH2.0 spec.
// Some values (mentioned below) are modified to better suit Berachain's system.

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

	// Signature domains.
	DefaultDomainTypeProposer          = 0
	DefaultDomainTypeAttester          = 1
	DefaultDomainTypeRandao            = 2
	DefaultDomainTypeDeposit           = 3
	DefaultDomainTypeVoluntaryExit     = 4
	DefaultDomainTypeSelectionProof    = 5
	DefaultDomainTypeAggregateAndProof = 6
	DefaultDomainTypeApplicationMask   = 16777216 // "0x00000001" in little endian

	// Eth1-related values.
	DefaultDepositContractAddress    = "0x4242424242424242424242424242424242424242" // Berachain specific.
	DefaultMaxDepositsPerBlock       = 16
	DefaultDepositEth1ChainID        = 1
	DefaultEth1FollowDistance        = 1 // Berachain specific.
	DefaultTargetSecondsPerEth1Block = 2 // Berachain specific.

	// Fork-related values.
	DefaultDeneb1ForkEpoch  = 9999999999999998 // Set as a future epoch as not yet determined.
	DefaultElectraForkEpoch = 9999999999999999 // Set as a future epoch as not yet determined.

	// State list length constants.
	DefaultEpochsPerHistoricalVector = 8
	DefaultEpochsPerSlashingsVector  = 8
	DefaultHistoricalRootsLimit      = 8
	DefaultValidatorRegistryLimit    = 1099511627776

	// Slashing.
	DefaultInactivityPenaltyQuotient      = 16777216
	DefaultProportionalSlashingMultiplier = 1

	// Capella values.
	DefaultMaxWithdrawalsPerPayload         = 16
	DefaultMaxValidatorsPerWithdrawalsSweep = 1 << 14

	// Deneb values.
	DefaultMinEpochsForBlobsSidecarsRequest = 4096
	DefaultMaxBlobCommitmentsPerBlock       = 4096
	DefaultMaxBlobsPerBlock                 = 6
	DefaultFieldElementsPerBlob             = 4096
	DefaultBytesPerBlob                     = 131072
	DefaultKZGCommitmentInclusionProofDepth = 17

	// Berachain values.
	DefaultValidatorSetCap      = 256
	DefaultEVMInflationAddress  = "0x0000000000000000000000000000000000000000"
	DefaultEVMInflationPerBlock = 0
)
