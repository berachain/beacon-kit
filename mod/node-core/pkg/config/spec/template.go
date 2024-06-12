// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

//nolint:lll // template lines get long
const Template = `
###############################################################################
###                               ChainSpec                                 ###
###############################################################################
[chain-spec]
# Gwei value constants.
#
# MinDepositAmount is the minimum deposit amount per deposit transaction.
min-deposit-amount = "{{ .ChainSpec.MinDepositAmount }}"
# MaxEffectiveBalance is the maximum effective balance allowed for a validator.
max-effective-balance = "{{ .ChainSpec.MaxEffectiveBalance }}"
# EjectionBalance is the balance at which a validator is ejected.
ejection-balance = "{{ .ChainSpec.EjectionBalance }}"
# EffectiveBalanceIncrement is the effective balance increment.
effective-balance-increment = "{{ .ChainSpec.EffectiveBalanceIncrement }}"

# Time parameters constants.
#
# SlotsPerEpoch is the number of slots per epoch.
slots-per-epoch = "{{ .ChainSpec.SlotsPerEpoch }}"
# SlotsPerHistoricalRoot is the number of slots per historical root.
slots-per-historical-root = "{{ .ChainSpec.SlotsPerHistoricalRoot }}"
# MinEpochsToInactivityPenalty is the minimum number of epochs before a validator is penalized for inactivity.
min-epochs-to-inactivity-penalty = "{{ .ChainSpec.MinEpochsToInactivityPenalty }}"

# Signature domains.
#
# DomainTypeProposer is the domain for beacon proposer signatures.
domain-type-beacon-proposer = "{{ .ChainSpec.DomainTypeProposer }}"
# DomainTypeAttester is the domain for beacon attester signatures.
domain-type-beacon-attester = "{{ .ChainSpec.DomainTypeAttester }}"
# DomainTypeRandao is the domain for RANDAO reveal signatures.
domain-type-randao = "{{ .ChainSpec.DomainTypeRandao }}"
# DomainTypeDeposit is the domain for deposit contract signatures.
domain-type-deposit = "{{ .ChainSpec.DomainTypeDeposit }}"
# DomainTypeVoluntaryExit is the domain for voluntary exit signatures.
domain-type-voluntary-exit = "{{ .ChainSpec.DomainTypeVoluntaryExit }}"
# DomainTypeSelectionProof is the domain for selection proof signatures.
domain-type-selection-proof = "{{ .ChainSpec.DomainTypeSelectionProof }}"
# DomainTypeAggregateAndProof is the domain for aggregate and proof signatures.
domain-type-aggregate-and-proof = "{{ .ChainSpec.DomainTypeAggregateAndProof }}"
# DomainTypeApplicationMask is the domain for the application mask.
domain-type-application-mask = "{{ .ChainSpec.DomainTypeApplicationMask }}"

# Eth1-related values.
#
# DepositContractAddress is the address of the deposit contract.
deposit-contract-address = "{{ .ChainSpec.DepositContractAddress }}"
# MaxDepositsPerBlock specifies the maximum number of deposit operations allowed per block.
max-deposits-per-block = "{{ .ChainSpec.MaxDepositsPerBlock }}"
# DepositEth1ChainID is the chain ID of the execution client.
deposit-eth1-chain-id = "{{ .ChainSpec.DepositEth1ChainID }}"
# Eth1FollowDistance is the distance between the eth1 chain and the beacon chain with respect to reading deposits.
eth1-follow-distance = "{{ .ChainSpec.Eth1FollowDistance }}"
# TargetSecondsPerEth1Block is the target time between eth1 blocks.
target-seconds-per-eth1-block = "{{ .ChainSpec.TargetSecondsPerEth1Block }}"

# Fork-related values.
#
# ElectraForkEpoch is the epoch at which the Electra fork is activated.
electra-fork-epoch = "{{ .ChainSpec.ElectraForkEpoch }}"

# State list lengths
#
# EpochsPerHistoricalVector is the number of epochs in the historical vector.
epochs-per-historical-vector = "{{ .ChainSpec.EpochsPerHistoricalVector }}"
# EpochsPerSlashingsVector is the number of epochs in the slashings vector.
epochs-per-slashings-vector = "{{ .ChainSpec.EpochsPerSlashingsVector }}"
# HistoricalRootsLimit is the maximum number of historical roots.
historical-roots-limit = "{{ .ChainSpec.HistoricalRootsLimit }}"
# ValidatorRegistryLimit is the maximum number of validators in the registry.
validator-registry-limit = "{{ .ChainSpec.ValidatorRegistryLimit }}"

# Rewards and penalties constants.
#
# InactivityPenaltyQuotient is the penalty quotient for inactivity.
inactivity-penalty-quotient = "{{ .ChainSpec.InactivityPenaltyQuotient }}"
# ProportionalSlashingMultiplier is the slashing multiplier relative to the base penalty.
proportional-slashing-multiplier = "{{ .ChainSpec.ProportionalSlashingMultiplier }}"

# Capella Values
#
# MaxWithdrawalsPerPayload indicates the maximum number of withdrawal operations allowed in a single payload.
max-withdrawals-per-payload = "{{ .ChainSpec.MaxWithdrawalsPerPayload }}"
# MaxValidatorsPerWithdrawalsSweep specifies the maximum number of validator withdrawals allowed per sweep.
max-validators-per-withdrawals-sweep = "{{ .ChainSpec.MaxValidatorsPerWithdrawalsSweep }}"

# Deneb Values
#
# MinEpochsForBlobsSidecarsRequest is the minimum number of epochs the node will keep the blobs for.
min-epochs-for-blobs-sidecars-request = "{{ .ChainSpec.MinEpochsForBlobsSidecarsRequest }}"
# MaxBlobCommitmentsPerBlock specifies the maximum number of blob commitments allowed per block.
max-blob-commitments-per-block = "{{ .ChainSpec.MaxBlobCommitmentsPerBlock }}"
# MaxBlobsPerBlock specifies the maximum number of blobs allowed per block.
max-blobs-per-block = "{{ .ChainSpec.MaxBlobsPerBlock }}"
# FieldElementsPerBlob specifies the number of field elements per blob.
field-elements-per-blob = "{{ .ChainSpec.FieldElementsPerBlob }}"
# BytesPerBlob denotes the size of EIP-4844 blobs in bytes.
bytes-per-blob = "{{ .ChainSpec.BytesPerBlob }}"
# KZGCommitmentInclusionProofDepth is the depth of the KZG inclusion proof.
kzg-commitment-inclusion-proof-depth = "{{ .ChainSpec.KZGCommitmentInclusionProofDepth }}"

# Comet Values
#
# CometBFTValues is the configuration for the CometBFT consensus engine.
comet-bft-config = "{{ .ChainSpec.CometBFTValues }}"
`
