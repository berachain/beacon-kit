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
//

package cache

import (
	"errors"

	"github.com/berachain/beacon-kit/consensus/cometbft/service/delay"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
)

const minActivationHeight = 2

func IsStateCachingActive(cfg delay.ConfigGetter, height math.Slot) bool {
	// We choose to activate caching only after stable block time is activated
	// so to avoid rushes of blocks. In any case however caching cannot be activated
	// for genesis and first block, which has special handling.
	activationHeight := math.Slot(max(cfg.SbtConsensusEnableHeight(), minActivationHeight)) //#nosec: G115
	return height >= activationHeight
}

var (
	_ States = (*candidateStates)(nil)

	ErrStateNotFound          = errors.New("state not found")
	ErrNoFinalState           = errors.New("no state marked as final")
	ErrFinalStateIsNil        = errors.New("state marked as final is nil")
	ErrFinalizingUnknownState = errors.New("attempt at finalizing unknown state")
)

// States is the interface controlling the caching of state and metadata once a block
// has been succeesfully executed once. States assumes that consensus allows no forks,
// and that no blocks at height H+1 is ever processed if height H has not been finalized.
type States interface {
	// SetCached allows caching state and metadata once a block has been executed.
	// It should be called for blocks which successfully pass PrepareProposal.
	// toCache is simply cached without checks, since CometBFT should guarantee that
	// blocks are verified at most once.
	SetCached(hash string, toCache *Element)

	// GetCached allows retrieval of cached state and metadata.
	// It should be called by FinalizedBlock to verify whether the block
	// has already been executed
	GetCached(hash string) (*Element, error)

	// MarkAsFinal allows marking one of the cached states as the one to be finalized and committed.
	// It should be called upon FinalizeBlock to then allow Commit to store the state.
	// It errs if we attempt to finalize a state which was not previously cached
	MarkAsFinal(hash string) error

	// GetFinal allows to retrieve state and stateHash marked as final.
	// Is should be call by Commit to retrieve the state finalized in FinalizeBlock
	GetFinal() (string, *State, error)

	// Reset allows wiping out all the verified states once one of the has been committed
	Reset()
}

type Element struct {
	State      *State
	ValUpdates transition.ValidatorUpdates
}

type candidateStates struct {
	states         map[string]*Element
	finalStateHash *string
}

func New() States {
	return &candidateStates{
		states:         make(map[string]*Element),
		finalStateHash: nil,
	}
}

func (cs *candidateStates) SetCached(hash string, toCache *Element) {
	cs.states[hash] = toCache
}

func (cs *candidateStates) GetCached(hash string) (*Element, error) {
	cached, found := cs.states[hash]
	if !found {
		return nil, ErrStateNotFound
	}
	return cached, nil
}

func (cs *candidateStates) MarkAsFinal(hash string) error {
	_, found := cs.states[hash]
	if !found {
		return ErrFinalizingUnknownState
	}
	cs.finalStateHash = &hash
	return nil
}

func (cs *candidateStates) GetFinal() (string, *State, error) {
	if cs.finalStateHash == nil {
		return "", nil, ErrNoFinalState
	}
	cached, found := cs.states[*cs.finalStateHash]
	if !found {
		return "", nil, ErrFinalStateIsNil
	}
	return *cs.finalStateHash, cached.State, nil
}

func (cs *candidateStates) Reset() {
	cs.states = make(map[string]*Element)
	cs.finalStateHash = nil
}
