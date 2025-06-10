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

	"cosmossdk.io/store/types"
	servercmtlog "github.com/berachain/beacon-kit/consensus/cometbft/service/log"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/transition"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ErrUnknownStateMarkedAsFinal = errors.New("unknown state marked as final")
	ErrNilFinalizeBlockState     = errors.New("finalizeBlockState is nil")
)

type CacheDirective uint8

const (
	InitialHeight int64 = 1

	DoNotCache CacheDirective = iota
	Cache
)

type CacheElement struct {
	state      *State
	valUpdates transition.ValidatorUpdates
}

// For now, just embed finalize state and make a proper interface for it.
// We will take care later on to cache states from processProposals
type FinalizedStateHandler struct {
	manager *Manager
	logger  *phuslu.Logger

	candidates map[string]*CacheElement
	finalized  *CacheElement
}

func NewFinalizeStateHandler(manager *Manager, logger *phuslu.Logger) *FinalizedStateHandler {
	return &FinalizedStateHandler{
		manager:    manager,
		logger:     logger,
		candidates: make(map[string]*CacheElement),
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
		log = servercmtlog.WrapSDKLogger(h.logger)

		ms     types.CacheMultiStore
		newCtx sdk.Context
	)
	switch {
	case height != InitialHeight: // the normal case
		ms = h.manager.GetCommitMultiStore().CacheMultiStore()
		newCtx = sdk.NewContext(ms, false, log).WithContext(ctx)

	case h.finalized == nil:
		// here height is InitialHeight. Also
		// here we are processing genesis via init chain or replaying blocks.
		// In any case this state will be final, but we will mark it as such
		// in MarkStateAsFinal
		ms = h.manager.GetCommitMultiStore().CacheMultiStore()
		newCtx = sdk.NewContext(ms, false, log).WithContext(ctx)
		if cd != Cache {
			return sdk.Context{}, errors.New("finalize state not yet initialized but state not cachable")
		}

	default:
		// this is the special case of a block built at InitialHeight.
		// We provide a cached context to allow access to
		// genesis state and resuse the multistore in State
		ms = h.finalized.state.ms
		newCtx, _ = h.finalized.state.Context().CacheContext()
	}

	if cd == Cache {
		h.candidates[string(blkHash)] = &CacheElement{
			state: NewState(ms, newCtx),
		}
	}
	return newCtx, nil
}

func (h *FinalizedStateHandler) PurgeState(blkHash []byte) error {
	k := string(blkHash)
	s, found := h.candidates[k]
	if !found {
		return nil // nothing to do
	}
	if h.finalized == s { // comparing pointers here, not content
		return fmt.Errorf("attempt at purging finalized state, blkHash %x", blkHash)
	}
	delete(h.candidates, k)
	return nil
}

func (h *FinalizedStateHandler) CachePostProcessBlockData(blkHash []byte, valUpdate transition.ValidatorUpdates) error {
	s, found := h.candidates[string(blkHash)]
	if !found {
		return fmt.Errorf("attempt at caching post block data for unknown state, blkHash %x", blkHash)
	}
	s.valUpdates = valUpdate
	return nil
}

func (h *FinalizedStateHandler) GetPostProcessBlockData(blkHash []byte) (transition.ValidatorUpdates, error) {
	s, found := h.candidates[string(blkHash)]
	if !found {
		return nil, fmt.Errorf("unknown state %x", blkHash)
	}
	return s.valUpdates, nil
}

func (h *FinalizedStateHandler) MarkStateAsFinal(blkHash []byte) error {
	s, found := h.candidates[string(blkHash)]
	if !found {
		return fmt.Errorf("%w, blkHash %x", ErrUnknownStateMarkedAsFinal, blkHash)
	}

	// note: we dont' purge s from candidate states in case CachePostProcessBlockData
	// is called on it after it is marked as final
	h.finalized = s
	return nil
}

func (h *FinalizedStateHandler) GetFinalizeStateContext() (sdk.Context, error) {
	if h.finalized == nil {
		return sdk.Context{}, ErrNilFinalizeBlockState
	}
	return h.finalized.state.Context(), nil
}

func (h *FinalizedStateHandler) WriteFinalizeState() ([]byte, error) {
	if h.finalized == nil {
		return nil, ErrNilFinalizeBlockState
	}
	h.finalized.state.Write()
	commitHash := h.manager.GetCommitMultiStore().WorkingHash()
	return commitHash, nil
}

func (h *FinalizedStateHandler) WipeState() {
	h.finalized = nil
	h.candidates = make(map[string]*CacheElement)
}
