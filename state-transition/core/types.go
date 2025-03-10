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

package core

import (
	stdbytes "bytes"
	"context"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
)

// ReadOnlyContext defines an interface for managing state transition context.
type ReadOnlyContext interface {
	ConsensusCtx() context.Context
	ConsensusTime() math.U64
	ProposerAddress() []byte
	VerifyPayload() bool
	VerifyRandao() bool
	VerifyResult() bool
	MeterGas() bool
}

// Withdrawals defines the interface for managing withdrawal operations.
type Withdrawals interface {
	Len() int
	EncodeIndex(int, *stdbytes.Buffer)
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// NotifyNewPayload notifies the execution client of the new payload.
	NotifyNewPayload(
		ctx context.Context,
		req *ctypes.NewPayloadRequest,
		retryOnSyncingStatus bool,
	) error
}

// Validator represents an interface for a validator with generic type
// ValidatorT.
type Validator[
	ValidatorT any,
] interface {
	constraints.SSZMarshallableRootable
	SizeSSZ(*ssz.Sizer) uint32
	// New creates a new validator with the given parameters.
	New(
		pubkey crypto.BLSPubkey,
		withdrawalCredentials ctypes.WithdrawalCredentials,
		amount math.Gwei,
		effectiveBalanceIncrement math.Gwei,
		maxEffectiveBalance math.Gwei,
	) ValidatorT
	// IsSlashed returns true if the validator is slashed.
	IsSlashed() bool

	IsEligibleForActivationQueue(threshold math.Gwei) bool
	IsEligibleForActivation(finalizedEpoch math.Epoch) bool
	IsActive(epoch math.Epoch) bool

	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator in
	// Gwei.
	GetEffectiveBalance() math.Gwei
	// SetEffectiveBalance sets the effective balance of the validator in Gwei.
	SetEffectiveBalance(math.Gwei)

	GetActivationEligibilityEpoch() math.Epoch
	SetActivationEligibilityEpoch(math.Epoch)

	GetActivationEpoch() math.Epoch
	SetActivationEpoch(math.Epoch)

	GetExitEpoch() math.Epoch
	SetExitEpoch(e math.Epoch)

	GetWithdrawableEpoch() math.Epoch
	SetWithdrawableEpoch(math.Epoch)
}

type Validators interface {
	HashTreeRoot() common.Root
}

// Withdrawal is the interface for a withdrawal.
type Withdrawal interface {
	// Equals returns true if the withdrawal is equal to the other.
	Equals(*engineprimitives.Withdrawal) bool
	// GetAmount returns the amount of the withdrawal.
	GetAmount() math.Gwei
	// GetIndex returns the public key of the validator.
	GetIndex() math.U64
	// GetValidatorIndex returns the index of the validator.
	GetValidatorIndex() math.ValidatorIndex
	// GetAddress returns the address of the withdrawal.
	GetAddress() common.ExecutionAddress
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	SetGauge(key string, value int64, args ...string)
	// IncrementCounter increments the counter identified by
	// the provided key.
	IncrementCounter(key string, args ...string)
}
