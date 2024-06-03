// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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

# Max operations per block constants.
#
# MaxDepositsPerBlock specifies the maximum number of deposit operations allowed per block.
max-deposits-per-block = "{{ .ChainSpec.MaxDepositsPerBlock }}"

# Rewards and penalties constants.
#
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
# CometValues is the consensus parameters for the Comet fork.
comet-bft-config = "{{ .ChainSpec.CometValues }}"
`
