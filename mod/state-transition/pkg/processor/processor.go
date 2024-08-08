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

package processor

import (
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// StateProcessor is a basic Processor, which takes care of the
// main state transition for the beacon chain.
type StateProcessor[
	BeaconBlockT BeaconBlock[
		DepositT, BeaconBlockBodyT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT,
		BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT, WithdrawalCredentialsT,
	],
	ContextT Context,
	DepositT Deposit[ForkDataT, WithdrawalCredentialsT],
	Eth1DataT interface {
		New(common.Root, math.U64, gethprimitives.ExecutionHash) Eth1DataT
		GetDepositCount() math.U64
	},
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkT interface {
		New(common.Version, common.Version, math.Epoch) ForkT
	},
	ForkDataT ForkData[ForkDataT],
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT interface {
		~[]ValidatorT
		HashTreeRoot() common.Root
	},
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT interface {
		~[32]byte
		ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
	},
	StateTransitionT interface {
		Transition(
			ctx ContextT,
			st BeaconStateT,
			blk BeaconBlockT,
		) (transition.ValidatorUpdates, error)
		ProcessSlots(
			st BeaconStateT,
			slot math.U64,
		) (transition.ValidatorUpdates, error)
		ProcessBlock(
			ctx ContextT,
			st BeaconStateT,
			blk BeaconBlockT,
		) error
		ExpectedWithdrawals(st BeaconStateT) ([]WithdrawalT, error)
		InitializePreminedBeaconStateFromEth1(
			st BeaconStateT,
			deposits []DepositT,
			executionPayloadHeader ExecutionPayloadHeaderT,
			genesisVersion common.Version,
		) (transition.ValidatorUpdates, error)
	},
] struct {
	// cs is the chain specification for the beacon chain.
	cs common.ChainSpec
	// signer is the BLS signer used for cryptographic operations.
	signer crypto.BLSSigner
	// executionEngine is the engine responsible for executing transactions.
	executionEngine ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	]
	transition StateTransitionT
}

// NewStateProcessor creates a new state processor.
func NewStateProcessor[
	BeaconBlockT BeaconBlock[
		DepositT, BeaconBlockBodyT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockBodyT BeaconBlockBody[
		BeaconBlockBodyT, DepositT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT,
		BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, KVStoreT,
		ValidatorT, ValidatorsT, WithdrawalT, WithdrawalCredentialsT,
	],
	ContextT Context,
	DepositT Deposit[ForkDataT, WithdrawalCredentialsT],
	Eth1DataT interface {
		New(common.Root, math.U64, gethprimitives.ExecutionHash) Eth1DataT
		GetDepositCount() math.U64
	},
	ExecutionPayloadT ExecutionPayload[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	ExecutionPayloadHeaderT ExecutionPayloadHeader,
	ForkT interface {
		New(common.Version, common.Version, math.Epoch) ForkT
	},
	ForkDataT ForkData[ForkDataT],
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT interface {
		~[]ValidatorT
		HashTreeRoot() common.Root
	},
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT interface {
		~[32]byte
		ToExecutionAddress() (gethprimitives.ExecutionAddress, error)
	},
	StateTransitionT interface {
		Transition(
			ctx ContextT,
			st BeaconStateT,
			blk BeaconBlockT,
		) (transition.ValidatorUpdates, error)
		ProcessSlots(
			st BeaconStateT,
			slot math.U64,
		) (transition.ValidatorUpdates, error)
		ProcessBlock(
			ctx ContextT,
			st BeaconStateT,
			blk BeaconBlockT,
		) error
		ExpectedWithdrawals(st BeaconStateT) ([]WithdrawalT, error)
		InitializePreminedBeaconStateFromEth1(
			st BeaconStateT,
			deposits []DepositT,
			executionPayloadHeader ExecutionPayloadHeaderT,
			genesisVersion common.Version,
		) (transition.ValidatorUpdates, error)
	},
](
	cs common.ChainSpec,
	executionEngine ExecutionEngine[
		ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	],
	signer crypto.BLSSigner,
	transition StateTransitionT,

) *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, ContextT, DepositT, Eth1DataT, ExecutionPayloadT,
	ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
	ValidatorsT, WithdrawalT, WithdrawalCredentialsT, StateTransitionT,
] {
	return &StateProcessor[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, ContextT, DepositT, Eth1DataT, ExecutionPayloadT,
		ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT, ValidatorT,
		ValidatorsT, WithdrawalT, WithdrawalCredentialsT, StateTransitionT,
	]{
		cs:              cs,
		executionEngine: executionEngine,
		signer:          signer,
		transition:      transition,
	}
}

// InitializePreminedBeaconStateFromEth1 initializes the beacon state
// from the given eth1 data.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT, _, DepositT,
	Eth1DataT, _, ExecutionPayloadHeaderT, ForkT, _, _, ValidatorT, _, _, _, _,
]) InitializePreminedBeaconStateFromEth1(
	st BeaconStateT,
	deposits []DepositT,
	executionPayloadHeader ExecutionPayloadHeaderT,
	genesisVersion common.Version,
) (transition.ValidatorUpdates, error) {
	return sp.transition.InitializePreminedBeaconStateFromEth1(
		st, deposits, executionPayloadHeader, genesisVersion,
	)
}

// Transition is the main function for processing a state transition.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT,
	_, _, _, _, _, _, _, _, _, _, _, _,
]) Transition(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
) (transition.ValidatorUpdates, error) {
	return sp.transition.Transition(ctx, st, blk)
}

// ProcessSlots processes the given number of slots.
func (sp *StateProcessor[
	_, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) ProcessSlots(st BeaconStateT, slot math.U64,
) (transition.ValidatorUpdates, error) {
	return sp.transition.ProcessSlots(st, slot)
}

// ProcessBlock processes the block, it optionally verifies the state root.
func (sp *StateProcessor[
	BeaconBlockT, _, _, BeaconStateT, ContextT,
	_, _, _, _, _, _, _, _, _, _, _, _,
]) ProcessBlock(
	ctx ContextT,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	return sp.transition.ProcessBlock(ctx, st, blk)
}

// ExpectedWithdrawals returns the expected withdrawals for the given state.
func (sp *StateProcessor[
	_, BeaconBlockBodyT, _, BeaconStateT, _, _, _, _,
	_, _, _, _, ValidatorT, ValidatorsT, WithdrawalT, _, _,
]) ExpectedWithdrawals(st BeaconStateT) ([]WithdrawalT, error) {
	return sp.transition.ExpectedWithdrawals(st)
}
