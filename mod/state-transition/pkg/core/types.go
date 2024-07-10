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

package core

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlock represents a generic interface for a beacon block.
type BeaconBlock[
	DepositT any,
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalsT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	WithdrawalsT any,
] interface {
	IsNil() bool
	// GetProposerIndex returns the index of the proposer.
	GetProposerIndex() math.ValidatorIndex
	// GetSlot returns the slot number of the block.
	GetSlot() math.Slot
	// GetBody returns the body of the block.
	GetBody() BeaconBlockBodyT
	// GetParentBlockRoot returns the root of the parent block.
	GetParentBlockRoot() common.Root
	// GetStateRoot returns the state root of the block.
	GetStateRoot() common.Root
}

// BeaconBlockBody represents a generic interface for the body of a beacon
// block.
type BeaconBlockBody[
	BeaconBlockBodyT any,
	DepositT any,
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	WithdrawalT any,
] interface {
	constraints.EmptyWithVersion[BeaconBlockBodyT]
	// GetRandaoReveal returns the RANDAO reveal signature.
	GetRandaoReveal() crypto.BLSSignature
	// GetExecutionPayload returns the execution payload.
	GetExecutionPayload() ExecutionPayloadT
	// GetDeposits returns the list of deposits.
	GetDeposits() []DepositT
	// HashTreeRoot returns the hash tree root of the block body.
	HashTreeRoot() ([32]byte, error)
	// GetBlobKzgCommitments returns the KZG commitments for the blobs.
	GetBlobKzgCommitments() eip4844.KZGCommitments[gethprimitives.ExecutionHash]
}

// BeaconBlockHeader is the interface for a beacon block header.
type BeaconBlockHeader[BeaconBlockHeaderT any] interface {
	New(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		stateRoot common.Root,
		bodyRoot common.Root,
	) BeaconBlockHeaderT
	HashTreeRoot() ([32]byte, error)
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	SetStateRoot(common.Root)
}

// Context defines an interface for managing state transition context.
type Context interface {
	context.Context
	// GetOptimisticEngine returns whether to optimistically assume the
	// execution client has the correct state when certain errors are returned
	// by the execution engine.
	GetOptimisticEngine() bool
	// GetSkipPayloadVerification returns whether to skip verifying the payload
	// if
	// it already exists on the execution client.
	GetSkipPayloadVerification() bool
	// GetSkipValidateRandao returns whether to skip validating the RANDAO
	// reveal.
	GetSkipValidateRandao() bool
	// GetSkipValidateResult returns whether to validate the result of the state
	// transition.
	GetSkipValidateResult() bool
}

// Deposit is the interface for a deposit.
type Deposit[
	ForkDataT any,
	WithdrawlCredentialsT ~[32]byte,
] interface {
	// GetAmount returns the amount of the deposit.
	GetAmount() math.Gwei
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetWithdrawalCredentials returns the withdrawal credentials.
	GetWithdrawalCredentials() WithdrawlCredentialsT
	// VerifySignature verifies the deposit and creates a validator.
	VerifySignature(
		forkData ForkDataT,
		domainType common.DomainType,
		signatureVerificationFn func(
			pubkey crypto.BLSPubkey,
			message []byte, signature crypto.BLSSignature,
		) error,
	) error
}

type ExecutionPayload[
	ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT any,
] interface {
	constraints.EngineType[ExecutionPayloadT]
	GetTransactions() [][]byte
	GetParentHash() gethprimitives.ExecutionHash
	GetBlockHash() gethprimitives.ExecutionHash
	GetPrevRandao() common.Bytes32
	GetWithdrawals() []WithdrawalT
	GetFeeRecipient() gethprimitives.ExecutionAddress
	GetStateRoot() common.Bytes32
	GetReceiptsRoot() common.Root
	GetLogsBloom() []byte
	GetNumber() math.U64
	GetGasLimit() math.U64
	GetTimestamp() math.U64
	GetGasUsed() math.U64
	GetExtraData() []byte
	GetBaseFeePerGas() math.U256L
	GetBlobGasUsed() math.U64
	GetExcessBlobGas() math.U64
	ToHeader(
		txsMerkleizer *merkle.Merkleizer[[32]byte, common.Root],
		maxWithdrawalsPerPayload uint64,
	) (ExecutionPayloadHeaderT, error)
}

type ExecutionPayloadHeader interface {
	GetBlockHash() gethprimitives.ExecutionHash
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine[
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT],
	ExecutionPayloadHeaderT any,
	WithdrawalT Withdrawal[WithdrawalT],
] interface {
	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[ExecutionPayloadT, WithdrawalT],
	) error
}

// ForkData is the interface for the fork data.
type ForkData[ForkDataT any] interface {
	// New creates a new fork data object.
	New(common.Version, common.Root) ForkDataT
	// ComputeRandaoSigningRoot returns the signing root for the fork data.
	ComputeRandaoSigningRoot(
		domainType common.DomainType,
		epoch math.Epoch,
	) (common.Root, error)
}

// Validator represents an interface for a validator with generic type
// ValidatorT.
type Validator[
	ValidatorT any,
	WithdrawalCredentialsT ~[32]byte,
] interface {
	constraints.SSZMarshallable
	// New creates a new validator with the given parameters.
	New(
		pubkey crypto.BLSPubkey,
		withdrawalCredentials WithdrawalCredentialsT,
		amount math.Gwei,
		effectiveBalanceIncrement math.Gwei,
		maxEffectiveBalance math.Gwei,
	) ValidatorT
	// IsSlashed returns true if the validator is slashed.
	IsSlashed() bool
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator in
	// Gwei.
	GetEffectiveBalance() math.Gwei
	// SetEffectiveBalance sets the effective balance of the validator in Gwei.
	SetEffectiveBalance(math.Gwei)
	// GetWithdrawableEpoch returns the epoch when the validator can withdraw.
	GetWithdrawableEpoch() math.Epoch
}

// Withdrawal is the interface for a withdrawal.
type Withdrawal[WithdrawalT any] interface {
	// Equals returns true if the withdrawal is equal to the other.
	Equals(WithdrawalT) bool
	// GetAmount returns the amount of the withdrawal.
	GetAmount() math.Gwei
	// GetIndex returns the public key of the validator.
	GetIndex() math.U64
	// GetValidatorIndex returns the index of the validator.
	GetValidatorIndex() math.ValidatorIndex
	// GetAddress returns the address of the withdrawal.
	GetAddress() gethprimitives.ExecutionAddress
}
