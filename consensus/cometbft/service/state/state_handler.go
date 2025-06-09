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
	"fmt"

	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	"github.com/berachain/beacon-kit/log/phuslu"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var ErrNilFinalizeBlockState = errors.New("finalizeBlockState is nil")

type CacheDirective uint8

const (
	InitialHeight int64 = 1

	DoNotCache CacheDirective = iota
	Cache
)

// For now, just embed finalize state and make a proper interface for it.
// We will take care later on to cache states from processProposals
type FinalizedStateHandler struct {
	manager *Manager
	logger  *phuslu.Logger

	candidateStates    map[string]*State
	finalizeBlockState *State
}

func NewFinalizeStateHandler(manager *Manager, logger *phuslu.Logger) *FinalizedStateHandler {
	return &FinalizedStateHandler{
		manager:         manager,
		logger:          logger,
		candidateStates: make(map[string]*State),
	}
}

// NewStateCtx returns the correct Context for relevant CometBFT callbacks.
// We use finalizeBlockState on the first block to be able to
// access any state changes made in InitChain. Also we properly cache a state
// indexed by the blkHash.
func (h *FinalizedStateHandler) NewStateCtx(
	ctx context.Context,
	height int64,
	blkHash []byte,
	cd CacheDirective,
) (sdk.Context, error) {
	var (
		log    = servercmtlog.WrapSDKLogger(h.logger)
		ms     = h.manager.GetCommitMultiStore().CacheMultiStore()
		newCtx = sdk.NewContext(ms, false, log).WithContext(ctx)
	)

	if height != InitialHeight {
		if cd == Cache {
			h.candidateStates[string(blkHash)] = NewState(ms, newCtx)
		}
		return newCtx, nil
	}

	// hereinafter height is InitialHeight
	if h.finalizeBlockState == nil {
		// here we are processing genesis via init chain or replaying blocks.
		// In any case this state will be final, but we will mark it as such
		// in MarkStateAsFinal
		if cd != Cache {
			return sdk.Context{}, errors.New("finalize state not yet initialized but state not cachable")
		}
		h.candidateStates[string(blkHash)] = NewState(ms, newCtx)
		return newCtx, nil
	}

	// this is the special case of a block built at InitialHeight. We provide a cached context
	// to allow access to genesis state and resuse the multistore in State
	newCtx, _ = h.finalizeBlockState.Context().CacheContext()
	if cd == Cache {
		h.candidateStates[string(blkHash)] = NewState(h.finalizeBlockState.ms, newCtx)
	}
	return newCtx, nil
}

func (h *FinalizedStateHandler) MarkStateAsFinal(blkHash []byte) error {
	s, found := h.candidateStates[string(blkHash)]
	if !found {
		return fmt.Errorf("attempt to make unknown state as final, blkHash %x", blkHash)
	}
	h.finalizeBlockState = s
	return nil
}

func (h *FinalizedStateHandler) GetFinalizeStateContext() (sdk.Context, error) {
	if h.finalizeBlockState == nil {
		return sdk.Context{}, ErrNilFinalizeBlockState
	}
	return h.finalizeBlockState.Context(), nil
}

func (h *FinalizedStateHandler) WriteFinalizeState() ([]byte, error) {
	if h.finalizeBlockState == nil {
		return nil, ErrNilFinalizeBlockState
	}
	h.finalizeBlockState.Write()
	commitHash := h.manager.GetCommitMultiStore().WorkingHash()
	return commitHash, nil
}

func (h *FinalizedStateHandler) WipeState() {
	h.finalizeBlockState = nil
	h.candidateStates = make(map[string]*State)
}
