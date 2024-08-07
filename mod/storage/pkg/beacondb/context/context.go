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

package context

import (
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/server/v2/stf/branch"
	storev2 "cosmossdk.io/store/v2"
)

type Context struct {
	context.Context

	actor     []byte
	WriterMap store.WriterMap
}

func Wrap(rawCtx context.Context, actor []byte) *Context {
	if ctx, ok := rawCtx.(*Context); ok {
		return ctx
	}
	return &Context{
		Context:   rawCtx,
		WriterMap: nil,
		actor:     actor,
	}
}

func (c *Context) CacheCopy(store storev2.RootStore) (*Context, error) {
	_, state, err := store.StateLatest()
	if err != nil {
		return nil, err
	}
	writerMap := branch.DefaultNewWriterMap(state)
	changes, err := c.WriterMap.GetStateChanges()
	if err != nil {
		return nil, err
	}
	if err = writerMap.ApplyStateChanges(changes); err != nil {
		return nil, err
	}
	return &Context{
		Context:   c.Context,
		actor:     c.actor,
		WriterMap: writerMap,
	}, nil
}

func (c *Context) AttachStore(store storev2.RootStore) error {
	if c.WriterMap != nil {
		return nil
	}
	_, state, err := store.StateLatest()
	if err != nil {
		return err
	}
	c.WriterMap = branch.DefaultNewWriterMap(state)
	return nil
}
