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

package cometbft

import (
	"errors"

	"github.com/berachain/beacon-kit/primitives/transition"
)

var (
	ErrStateNotFound          = errors.New("state not found")
	ErrNoFinalState           = errors.New("no state marked as final")
	ErrFinalStateIsNil        = errors.New("state marked as final is nil")
	ErrFinalizingUnknownState = errors.New("attempt at finalizing unknown state")
)

type CacheElement struct {
	state      *state
	valUpdates transition.ValidatorUpdates
}

type candidateStates struct {
	states         map[string]*CacheElement
	finalStateHash *string
}

func newCandidateStates() *candidateStates {
	return &candidateStates{
		states:         make(map[string]*CacheElement),
		finalStateHash: nil,
	}
}

func (cs *candidateStates) cache(hash string, toCache *CacheElement) {
	cs.states[hash] = toCache
}

func (cs *candidateStates) getState(hash string) (*CacheElement, error) {
	cached, found := cs.states[hash]
	if !found {
		return nil, ErrStateNotFound
	}
	return cached, nil
}

func (cs *candidateStates) markAsFinal(hash string) error {
	_, found := cs.states[hash]
	if !found {
		return ErrFinalizingUnknownState
	}
	cs.finalStateHash = &hash
	return nil
}

func (cs *candidateStates) getFinalState() (string, *state, error) {
	if cs.finalStateHash == nil {
		return "", nil, ErrNoFinalState
	}
	cached, found := cs.states[*cs.finalStateHash]
	if !found {
		return "", nil, ErrFinalStateIsNil
	}
	return *cs.finalStateHash, cached.state, nil
}

func (cs *candidateStates) reset() {
	cs.states = make(map[string]*CacheElement)
	cs.finalStateHash = nil
}
