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

package cometbft

import (
	"fmt"
	"sync"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type state struct {
	ms storetypes.CacheMultiStore

	mtx sync.RWMutex
	ctx sdk.Context
}

// CacheMultiStore calls and returns a CacheMultiStore on the state's underling
// CacheMultiStore.
func (st *state) CacheMultiStore() storetypes.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

// SetContext updates the state's context to the context provided.
func (st *state) SetContext(ctx sdk.Context) {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	st.ctx = ctx
}

// Context returns the Context of the state.
func (st *state) Context() sdk.Context {
	st.mtx.RLock()
	defer st.mtx.RUnlock()
	return st.ctx
}

// ephemeralStates collects all the states related to a block
// which has been verified successfully. They are candidate states
// to be finalized
type ephemeralStates struct {
	states map[int64]map[string]*state // blkHeight -> blkHash -> state
}

func newEphemeralStates() *ephemeralStates {
	return &ephemeralStates{
		states: make(map[int64]map[string]*state),
	}
}

func (es *ephemeralStates) CacheVerified(blkHeight int64, blkHash []byte, blkState *state) error {
	states, found := es.states[blkHeight]
	if !found {
		states = make(map[string]*state)
		es.states[blkHeight] = states
	}
	if _, found = states[string(blkHash)]; found {
		return fmt.Errorf("overwriting state height %d, blk hash %s", blkHeight, blkHash)
	}
	states[string(blkHash)] = blkState
	return nil
}

func (es *ephemeralStates) Get(blkHeight int64, blkHash []byte) (*state, bool) {
	states, found := es.states[blkHeight]
	if !found {
		return nil, found
	}
	res, found := states[string(blkHash)]
	return res, found
}

func (es *ephemeralStates) Clear(blkHeight int64) error {
	if _, found := es.states[blkHeight]; !found {
		return fmt.Errorf("attempt at deleting unknown blkHeight %d", blkHeight)
	}
	delete(es.states, blkHeight)
	return nil
}
