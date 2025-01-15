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

package chain

import "github.com/berachain/beacon-kit/primitives/common"

// SpecData is the underlying data structure for chain-specific parameters.
type SpecData[
	DomainTypeT ~[4]byte,
	EpochT ~uint64,
	SlotT ~uint64,
	CometBFTConfigT any,
] struct {
	// Gwei value constants.
	//
	// MinDepositAmount is the minimum deposit amount per deposit
	// transaction.
	MinDepositAmount uint64 `mapstructure:"min-deposit-amount"`
	// MaxEffectiveBalance is the maximum effective balance allowed for a
	// validator before the upgrade.
	MaxEffectiveBalancePreUpgrade uint64 `mapstructure:"max-effective-balance-pre-upgrade"`
	// MaxEffectiveBalancePostUpgrade is the maximum effective balance allowed
	// for a validator after the upgrade.
	MaxEffectiveBalancePostUpgrade uint64 `mapstructure:"max-effective-balance-post-upgrade"`
	// EjectionBalance is the balance at which a validator is ejected.
	EjectionBalance uint64 `mapstructure:"ejection-balance"`
	// EffectiveBalanceIncrement is the effective balance increment.
	EffectiveBalanceIncrement uint64 `mapstructure:"effective-balance-increment"`

	// HysteresisQuotient is the quotient used in effective balance calculations
	HysteresisQuotient uint64 `mapstructure:"hysteresis-quotient"`
	// HysteresisDownwardMultiplier is the multiplier for downward balance
	// adjustments.
	HysteresisDownwardMultiplier uint64 `mapstructure:"hysteresis-downward-multiplier"`
	// HysteresisUpwardMultiplier is the multiplier for upward balance
	// adjustments.
	HysteresisUpwardMultiplier uint64 `mapstructure:"hysteresis-upward-multiplier"`

	// Time parameters constants.
	//
	// SlotsPerEpoch is the number of slots per epoch.
	SlotsPerEpoch uint64 `mapstructure:"slots-per-epoch"`
	// SlotsPerHistoricalRoot is the number of slots per historical root.
	SlotsPerHistoricalRoot uint64 `mapstructure:"slots-per-historical-root"`
	// MinEpochsToInactivityPenalty is the minimum number of epochs before a
	// validator is penalized for inactivity.
	MinEpochsToInactivityPenalty uint64 `mapstructure:"min-epochs-to-inactivity-penalty"`

	// Signature domains.
	//
	// DomainDomainTypeProposerProposer is the domain for beacon proposer
	// signatures.
	DomainTypeProposer DomainTypeT `mapstructure:"domain-type-beacon-proposer"`
	// DomainTypeAttester is the domain for beacon attester signatures.
	DomainTypeAttester DomainTypeT `mapstructure:"domain-type-beacon-attester"`
	// DomainTypeRandao is the domain for RANDAO reveal signatures.
	DomainTypeRandao DomainTypeT `mapstructure:"domain-type-randao"`
	// DomainTypeDeposit is the domain for deposit contract signatures.
	DomainTypeDeposit DomainTypeT `mapstructure:"domain-type-deposit"`
	// DomainTypeVoluntaryExit is the domain for voluntary exit signatures.
	DomainTypeVoluntaryExit DomainTypeT `mapstructure:"domain-type-voluntary-exit"`
	// DomainTypeSelectionProof is the domain for selection proof signatures.
	DomainTypeSelectionProof DomainTypeT `mapstructure:"domain-type-selection-proof"`
	// DomainTypeAggregateAndProof is the domain for aggregate and proof
	// signatures.
	DomainTypeAggregateAndProof DomainTypeT `mapstructure:"domain-type-aggregate-and-proof"`
	// DomainTypeApplicationMask is the domain for the application mask.
	DomainTypeApplicationMask DomainTypeT `mapstructure:"domain-type-application-mask"`

	// Eth1-related values.
	//
	// DepositContractAddress is the address of the deposit contract.
	DepositContractAddress common.ExecutionAddress `mapstructure:"deposit-contract-address"`
	// MaxDepositsPerBlock specifies the maximum number of deposit operations
	// allowed per block.
	MaxDepositsPerBlock uint64 `mapstructure:"max-deposits-per-block"`
	// DepositEth1ChainID is the chain ID of the execution client.
	DepositEth1ChainID uint64 `mapstructure:"deposit-eth1-chain-id"`
	// Eth1FollowDistance is the distance between the eth1 chain and the beacon
	// chain with respect to reading deposits.
	Eth1FollowDistance uint64 `mapstructure:"eth1-follow-distance"`
	// TargetSecondsPerEth1Block is the target time between eth1 blocks.
	TargetSecondsPerEth1Block uint64 `mapstructure:"target-seconds-per-eth1-block"`

	// Fork-related values.
	//
	// DenebPlus is the epoch at which the Deneb+ fork is activated.
	DenebPlusForkEpoch EpochT `mapstructure:"deneb-plus-fork-epoch"`
	// ElectraForkEpoch is the epoch at which the Electra fork is activated.
	ElectraForkEpoch EpochT `mapstructure:"electra-fork-epoch"`

	// State list lengths
	//
	// EpochsPerHistoricalVector is the number of epochs in the historical
	// vector.
	EpochsPerHistoricalVector uint64 `mapstructure:"epochs-per-historical-vector"`
	// EpochsPerSlashingsVector is the number of epochs in the slashings vector.
	EpochsPerSlashingsVector uint64 `mapstructure:"epochs-per-slashings-vector"`
	// HistoricalRootsLimit is the maximum number of historical roots.
	HistoricalRootsLimit uint64 `mapstructure:"historical-roots-limit"`
	// ValidatorRegistryLimit is the maximum number of validators in the
	// registry.
	ValidatorRegistryLimit uint64 `mapstructure:"validator-registry-limit"`

	// Rewards and penalties constants.
	//
	// InactivityPenaltyQuotient is the inactivity penalty quotient.
	InactivityPenaltyQuotient uint64 `mapstructure:"inactivity-penalty-quotient"`
	// ProportionalSlashingMultiplier is the slashing multiplier relative to the
	// base penalty.
	ProportionalSlashingMultiplier uint64 `mapstructure:"proportional-slashing-multiplier"`

	// Capella Values
	//
	// MaxWithdrawalsPerPayload indicates the maximum number of withdrawal
	// operations allowed in a single payload.
	MaxWithdrawalsPerPayload uint64 `mapstructure:"max-withdrawals-per-payload"`
	// MaxValidatorsPerWithdrawalsSweepPreUpgrade specifies the maximum number
	// of validator withdrawals allowed per sweep. Before the upgrade, this
	// value is consistent with ETH2.0 spec, set to 2^14.
	MaxValidatorsPerWithdrawalsSweepPreUpgrade uint64 `mapstructure:"max-validators-per-withdrawals-sweep-pre-upgrade"`
	// MaxValidatorsPerWithdrawalsSweepPostUpgrade specifies the maximum number
	// of validator withdrawals allowed per sweep. After the upgrade, this value
	// is set to the largest prime number smaller than the minimum possible
	// amount of total validators.
	MaxValidatorsPerWithdrawalsSweepPostUpgrade uint64 `mapstructure:"max-validators-per-withdrawals-sweep-post-upgrade"`

	// Deneb Values
	//
	// MinEpochsForBlobsSidecarsRequest is the minimum number of epochs the node
	// will keep the blobs for.
	MinEpochsForBlobsSidecarsRequest uint64 `mapstructure:"min-epochs-for-blobs-sidecars-request"`
	// MaxBlobCommitmentsPerBlock specifies the maximum number of blob
	// commitments allowed per block.
	MaxBlobCommitmentsPerBlock uint64 `mapstructure:"max-blob-commitments-per-block"`
	// MaxBlobsPerBlock specifies the maximum number of blobs allowed per block.
	MaxBlobsPerBlock uint64 `mapstructure:"max-blobs-per-block"`
	// FieldElementsPerBlob specifies the number of field elements per blob.
	FieldElementsPerBlob uint64 `mapstructure:"field-elements-per-blob"`
	// BytesPerBlob denotes the size of EIP-4844 blobs in bytes.
	BytesPerBlob uint64 `mapstructure:"bytes-per-blob"`
	// KZGCommitmentInclusionProofDepth is the depth of the KZG inclusion proof.
	KZGCommitmentInclusionProofDepth uint64 `mapstructure:"kzg-commitment-inclusion-proof-depth"`

	// Comet Values
	CometValues CometBFTConfigT `mapstructure:"comet-bft-config"`

	// Berachain Values
	//
	// ValidatorSetCap is the maximum number of validators that can be active
	// for a given epoch
	// Note: ValidatorSetCap must be smaller than ValidatorRegistryLimit.
	ValidatorSetCap uint64 `mapstructure:"validator-set-cap-size"`
	// EVMInflationAddress is the address on the EVM which will receive the
	// inflation amount of native EVM balance through a withdrawal every block.
	EVMInflationAddress common.ExecutionAddress `mapstructure:"evm-inflation-address"`
	// EVMInflationPerBlock is the amount of native EVM balance (in Gwei) to be
	// minted to the EVMInflationAddress via a withdrawal every block.
	EVMInflationPerBlock uint64 `mapstructure:"evm-inflation-per-block"`
}
