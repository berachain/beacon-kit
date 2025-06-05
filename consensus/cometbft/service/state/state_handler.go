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

package state

import (
	"context"
	"errors"

	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	"github.com/berachain/beacon-kit/log/phuslu"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ErrNilFinalizeBlockState = errors.New("finalizeBlockState is nil")

type Type uint8

const (
	Ephemeral Type = iota
	CandidateFinal

	InitialHeight int64 = 1
)

// For now, just embed finalize state and make a proper interface for it.
// We will take care later on to cache states from processProposals
type FinalizedStateHandler struct {
	manager            *Manager
	logger             *phuslu.Logger
	finalizeBlockState *State
}

func NewFinalizeStateHandler(manager *Manager, logger *phuslu.Logger) *FinalizedStateHandler {
	return &FinalizedStateHandler{
		manager: manager,
		logger:  logger,
	}
}

func (h *FinalizedStateHandler) ResetState(ctx context.Context, st Type) *State {
	var (
		log    = servercmtlog.WrapSDKLogger(h.logger)
		ms     = h.manager.GetCommitMultiStore().CacheMultiStore()
		newCtx = sdk.NewContext(ms, false, log).WithContext(ctx)
	)

	res := NewState(ms, newCtx)
	if st == CandidateFinal {
		h.finalizeBlockState = NewState(ms, newCtx)
	}
	return res
}

// GetContextForProposal returns the correct Context for PrepareProposal and
// ProcessProposal. We use finalizeBlockState on the first block to be able to
// access any state changes made in InitChain.
func (h *FinalizedStateHandler) GetContextForProposal(
	ctx sdk.Context,
	height int64,
) (sdk.Context, error) {
	if height != InitialHeight {
		return ctx, nil
	}

	if h.finalizeBlockState == nil {
		return sdk.Context{}, ErrNilFinalizeBlockState
	}
	newCtx, _ := h.finalizeBlockState.Context().CacheContext()
	return newCtx, nil
}

func (h *FinalizedStateHandler) GetSDKContext() (sdk.Context, error) {
	if h.finalizeBlockState == nil {
		return sdk.Context{}, ErrNilFinalizeBlockState
	}
	return h.finalizeBlockState.Context(), nil
}

func (h *FinalizedStateHandler) HasFinalizeState() bool {
	return h.finalizeBlockState != nil
}

func (h *FinalizedStateHandler) WipeState() {
	h.finalizeBlockState = nil
}

func (h *FinalizedStateHandler) WriteFinalizeState() ([]byte, error) {
	if h.finalizeBlockState == nil {
		return nil, ErrNilFinalizeBlockState
	}
	h.finalizeBlockState.Write()
	commitHash := h.manager.GetCommitMultiStore().WorkingHash()
	return commitHash, nil
}
