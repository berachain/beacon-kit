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

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	cmttypes "github.com/cometbft/cometbft/types"
)

const (
	// Default DepositContractAddress is the default address of the pre-deployed
	// beacon deposit contract.
	DefaultDepositContractAddress = "0x4242424242424242424242424242424242424242"
)

// CommonSpec returns a chain spec with default values.
//
//nolint:mnd // bet.
func CommonSpec() Common[any] {
	cmtConsensusParams := cmttypes.DefaultConsensusParams()
	cmtConsensusParams.Validator.PubKeyTypes = []string{crypto.CometBLSType}

	return Common[any]{
		// Gwei value constants.
		MinDepositAmount:             1e9,
		MaxEffectiveBalance:          32e9,
		EjectionBalance:              16e9,
		EffectiveBalanceIncrement:    1e9,
		HysteresisQuotient:           4,
		HysteresisDownwardMultiplier: 1,
		HysteresisUpwardMultiplier:   5,

		// Time parameters constants.
		SlotsPerEpoch:                32,
		MinEpochsToInactivityPenalty: 4,
		SlotsPerHistoricalRoot:       8,
		// Signature domains.
		DomainTypeProposer: common.DomainType{
			0x00, 0x00, 0x00, 0x00,
		},
		DomainTypeAttester: common.DomainType{
			0x01, 0x00, 0x00, 0x00,
		},
		DomainTypeRandao: common.DomainType{
			0x02, 0x00, 0x00, 0x00,
		},
		DomainTypeDeposit: common.DomainType{
			0x03, 0x00, 0x00, 0x00,
		},
		DomainTypeVoluntaryExit: common.DomainType{
			0x04, 0x00, 0x00, 0x00,
		},
		DomainTypeSelectionProof: common.DomainType{
			0x05, 0x00, 0x00, 0x00,
		},
		DomainTypeAggregateAndProof: common.DomainType{
			0x06, 0x00, 0x00, 0x00,
		},
		DomainTypeApplicationMask: common.DomainType{
			0x00, 0x00, 0x00, 0x01,
		},
		// Eth1-related values.
		DepositContractAddress: common.NewExecutionAddressFromHex(
			DefaultDepositContractAddress,
		),
		DepositEth1ChainID:        1,
		Eth1FollowDistance:        1,
		TargetSecondsPerEth1Block: 3,
		// Fork-related values.
		DenebPlusForkEpoch: 9999999999999998,
		ElectraForkEpoch:   9999999999999999,
		// State list length constants.
		EpochsPerHistoricalVector: 8,
		EpochsPerSlashingsVector:  8,
		HistoricalRootsLimit:      8,
		ValidatorRegistryLimit:    1099511627776,
		// Max operations per block constants.
		MaxDepositsPerBlock: 16,
		// Slashing
		ProportionalSlashingMultiplier: 1,
		// Capella values.
		MaxWithdrawalsPerPayload:         16,
		MaxValidatorsPerWithdrawalsSweep: 1 << 14,
		// Deneb values.
		MinEpochsForBlobsSidecarsRequest: 4096,
		MaxBlobCommitmentsPerBlock:       16,
		MaxBlobsPerBlock:                 6,
		FieldElementsPerBlob:             4096,
		BytesPerBlob:                     131072,
		KZGCommitmentInclusionProofDepth: 17,
		// Comet values.
		CometValues: cmtConsensusParams,
	}
}
