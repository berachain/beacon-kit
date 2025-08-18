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

package cache

import (
	"sync"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type State struct {
	ms storetypes.CacheMultiStore

	mtx sync.RWMutex
	ctx sdk.Context
}

func NewState(ms storetypes.CacheMultiStore, ctx sdk.Context) *State {
	return &State{
		ms:  ms,
		ctx: ctx,
	}
}

// CacheMultiStore calls and returns a CacheMultiStore on the state's underling
// CacheMultiStore.
func (st *State) CacheMultiStore() storetypes.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

// SetContext updates the state's context to the context provided.
func (st *State) SetContext(ctx sdk.Context) {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	st.ctx = ctx
}

// Context returns the Context of the state.
func (st *State) Context() sdk.Context {
	st.mtx.RLock()
	defer st.mtx.RUnlock()
	return st.ctx
}

func (st *State) Write() {
	st.ms.Write()
}
