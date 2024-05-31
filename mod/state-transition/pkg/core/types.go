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

package core

import (
	"context"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT any, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool

	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// BeaconBlock represents a generic interface for a beacon block.
type BeaconBlock[
	DepositT any,
	BeaconBlockBodyT BeaconBlockBody[
		DepositT, ExecutionPayloadT, WithdrawalsT,
	],
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT, WithdrawalsT],
	WithdrawalsT any,
] interface {
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
	DepositT any,
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT, WithdrawalsT],
	WithdrawalsT any,
] interface {
	// GetRandaoReveal returns the RANDAO reveal signature.
	GetRandaoReveal() crypto.BLSSignature
	// GetExecutionPayload returns the execution payload.
	GetExecutionPayload() ExecutionPayloadT
	// GetDeposits returns the list of deposits.
	GetDeposits() []DepositT
	// HashTreeRoot returns the hash tree root of the block body.
	HashTreeRoot() ([32]byte, error)
	// GetBlobKzgCommitments returns the KZG commitments for the blobs.
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
}

// BlobSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	Len() int
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

	// Unwrap returns the underlying golang standard library context.
	Unwrap() context.Context
}

// Deposit is the interface for a deposit.
type Deposit[
	ForkDataT any,
	WithdrawlCredentialsT ~[32]byte,
] interface {
	// GetAmount returns the amount of the deposit.
	GetAmount() math.Gwei
	// GetIndex returns the index of the deposit.
	GetIndex() uint64
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetSignature returns the signature of the deposit.
	GetSignature() crypto.BLSSignature
	// GetWithdrawalCredentials returns the withdrawal credentials.
	GetWithdrawalCredentials() WithdrawlCredentialsT
	// VerifySignature verifies the deposit and creates a validator.
	VerifySignature(
		forkData ForkDataT,
		domainType common.DomainType,
		signatureVerificationFn func(
			pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
		) error,
	) error
}

type ExecutionPayload[ExecutionPayloadT, WithdrawalsT any] interface {
	Empty(uint32) ExecutionPayloadT
	Version() uint32
	GetTransactions() [][]byte
	GetParentHash() common.ExecutionHash
	GetBlockHash() common.ExecutionHash
	GetPrevRandao() bytes.B32
	GetWithdrawals() []WithdrawalsT
	GetFeeRecipient() common.ExecutionAddress
	GetStateRoot() bytes.B32
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
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine[
	ExecutionPayloadT ExecutionPayload[ExecutionPayloadT, WithdrawalsT],
	WithdrawalsT any,
] interface {
	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[ExecutionPayloadT],
	) error
}

// ForkData is the interface for the fork data.
type ForkData[ForkDataT any] interface {
	// New creates a new fork data object.
	New(primitives.Version, primitives.Root) ForkDataT
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor[BeaconBlockT, BeaconStateT any] interface {
	// ProcessRandao processes the RANDAO reveal and ensures it
	// matches the local state.
	ProcessRandao(BeaconStateT, BeaconBlockT, bool) error
	// ProcessRandaoMixesReset resets the RANDAO mixes as defined
	// in the Ethereum 2.0 specification.
	ProcessRandaoMixesReset(BeaconStateT) error
}

// Validator represents an interface for a validator with generic type
// ValidatorT.
type Validator[
	ValidatorT any,
	WithdrawalCredentialsT ~[32]byte,
] interface {
	ssz.Marshaler
	ssz.HashRoot
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
