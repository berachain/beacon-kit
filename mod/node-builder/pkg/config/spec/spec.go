package spec

import (
	"bytes"
	"os"
	"text/template"
)

const DefaultTemplate = `
###############################################################################
###                          Spec Data Configuration                        ###
###############################################################################

# Gwei value constants.
#
# MinDepositAmount is the minimum deposit amount per deposit transaction.
min-deposit-amount = "{{ .MinDepositAmount }}"
# MaxEffectiveBalance is the maximum effective balance allowed for a validator.
max-effective-balance = "{{ .MaxEffectiveBalance }}"
# EjectionBalance is the balance at which a validator is ejected.
ejection-balance = "{{ .EjectionBalance }}"
# EffectiveBalanceIncrement is the effective balance increment.
effective-balance-increment = "{{ .EffectiveBalanceIncrement }}"

# Time parameters constants.
#
# SlotsPerEpoch is the number of slots per epoch.
slots-per-epoch = "{{ .SlotsPerEpoch }}"
# SlotsPerHistoricalRoot is the number of slots per historical root.
slots-per-historical-root = "{{ .SlotsPerHistoricalRoot }}"

# Signature domains.
#
# DomainTypeProposer is the domain for beacon proposer signatures.
domain-type-beacon-proposer = "{{ .DomainTypeProposer }}"
# DomainTypeAttester is the domain for beacon attester signatures.
domain-type-beacon-attester = "{{ .DomainTypeAttester }}"
# DomainTypeRandao is the domain for RANDAO reveal signatures.
domain-type-randao = "{{ .DomainTypeRandao }}"
# DomainTypeDeposit is the domain for deposit contract signatures.
domain-type-deposit = "{{ .DomainTypeDeposit }}"
# DomainTypeVoluntaryExit is the domain for voluntary exit signatures.
domain-type-voluntary-exit = "{{ .DomainTypeVoluntaryExit }}"
# DomainTypeSelectionProof is the domain for selection proof signatures.
domain-type-selection-proof = "{{ .DomainTypeSelectionProof }}"
# DomainTypeAggregateAndProof is the domain for aggregate and proof signatures.
domain-type-aggregate-and-proof = "{{ .DomainTypeAggregateAndProof }}"
# DomainTypeApplicationMask is the domain for the application mask.
domain-type-application-mask = "{{ .DomainTypeApplicationMask }}"

# Eth1-related values.
#
# DepositContractAddress is the address of the deposit contract.
deposit-contract-address = "{{ .DepositContractAddress }}"
# DepositEth1ChainID is the chain ID of the execution client.
deposit-eth1-chain-id = "{{ .DepositEth1ChainID }}"
# Eth1FollowDistance is the distance between the eth1 chain and the beacon chain with respect to reading deposits.
eth1-follow-distance = "{{ .Eth1FollowDistance }}"
# TargetSecondsPerEth1Block is the target time between eth1 blocks.
target-seconds-per-eth1-block = "{{ .TargetSecondsPerEth1Block }}"

# Fork-related values.
#
# ElectraForkEpoch is the epoch at which the Electra fork is activated.
electra-fork-epoch = "{{ .ElectraForkEpoch }}"

# State list lengths
#
# EpochsPerHistoricalVector is the number of epochs in the historical vector.
epochs-per-historical-vector = "{{ .EpochsPerHistoricalVector }}"
# EpochsPerSlashingsVector is the number of epochs in the slashings vector.
epochs-per-slashings-vector = "{{ .EpochsPerSlashingsVector }}"
# HistoricalRootsLimit is the maximum number of historical roots.
historical-roots-limit = "{{ .HistoricalRootsLimit }}"
# ValidatorRegistryLimit is the maximum number of validators in the registry.
validator-registry-limit = "{{ .ValidatorRegistryLimit }}"

# Max operations per block constants.
#
# MaxDepositsPerBlock specifies the maximum number of deposit operations allowed per block.
max-deposits-per-block = "{{ .MaxDepositsPerBlock }}"

# Rewards and penalties constants.
#
# ProportionalSlashingMultiplier is the slashing multiplier relative to the base penalty.
proportional-slashing-multiplier = "{{ .ProportionalSlashingMultiplier }}"

# Capella Values
#
# MaxWithdrawalsPerPayload indicates the maximum number of withdrawal operations allowed in a single payload.
max-withdrawals-per-payload = "{{ .MaxWithdrawalsPerPayload }}"
# MaxValidatorsPerWithdrawalsSweep specifies the maximum number of validator withdrawals allowed per sweep.
max-validators-per-withdrawals-sweep = "{{ .MaxValidatorsPerWithdrawalsSweep }}"

# Deneb Values
#
# MinEpochsForBlobsSidecarsRequest is the minimum number of epochs the node will keep the blobs for.
min-epochs-for-blobs-sidecars-request = "{{ .MinEpochsForBlobsSidecarsRequest }}"
# MaxBlobCommitmentsPerBlock specifies the maximum number of blob commitments allowed per block.
max-blob-commitments-per-block = "{{ .MaxBlobCommitmentsPerBlock }}"
# MaxBlobsPerBlock specifies the maximum number of blobs allowed per block.
max-blobs-per-block = "{{ .MaxBlobsPerBlock }}"
# FieldElementsPerBlob specifies the number of field elements per blob.
field-elements-per-blob = "{{ .FieldElementsPerBlob }}"
# BytesPerBlob denotes the size of EIP-4844 blobs in bytes.
bytes-per-blob = "{{ .BytesPerBlob }}"
# KZGCommitmentInclusionProofDepth is the depth of the KZG inclusion proof.
kzg-commitment-inclusion-proof-depth = "{{ .KZGCommitmentInclusionProofDepth }}"`

func WriteSpecToFile(filepath string, spec any) error {
	tmpl, err := template.New("specTemplate").Parse(DefaultTemplate)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, spec); err != nil {
		return err
	}

	return os.WriteFile(filepath, buffer.Bytes(), 0o600)
}

// func MustReadSpecFromFile(filepath string) chain.Spec[
// 	common.DomainType,
// 	math.Epoch,
// 	common.ExecutionAddress,
// 	math.Slot,
// ] {
// 	spec, err := ReadSpecFromFile(filepath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return spec
// }

// func ReadSpecFromFile(filepath string) (chain.Spec[
// 	common.DomainType,
// 	math.Epoch,
// 	common.ExecutionAddress,
// 	math.Slot,
// ], error) {
// 	// spec, err := chain.ReadSpecFromFile(filepath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return spec, nil
// }
