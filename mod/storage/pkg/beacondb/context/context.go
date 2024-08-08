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
	storev2 "cosmossdk.io/store/v2"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/context/cache"
)

type Context struct {
	context.Context

	actor []byte
	Cache store.Writer
}

func Wrap(
	rawCtx context.Context,
	actor []byte,
	store storev2.RootStore,
) (*Context, error) {
	writer, err := latestWriterFromStore(store, actor)
	if err != nil {
		return nil, err
	}
	return &Context{
		Context: rawCtx,
		Cache:   writer,
		actor:   actor,
	}, nil
}

func (c *Context) CacheCopy(s storev2.RootStore) (*Context, error) {
	writer, err := latestWriterFromStore(s, c.actor)
	if err != nil {
		return nil, err
	}

	changes, err := c.Cache.ChangeSets()
	if err != nil {
		return nil, err
	}
	if err := writer.ApplyChangeSets(changes); err != nil {
		return nil, err
	}
	return &Context{
		Context: c.Context,
		actor:   c.actor,
		Cache:   writer,
	}, nil
}

func (c *Context) StateChanges() ([]store.StateChanges, error) {
	changes, err := c.Cache.ChangeSets()
	if err != nil {
		return nil, err
	}
	return []store.StateChanges{
		{
			Actor:        c.actor,
			StateChanges: changes,
		},
	}, nil
}

func latestWriterFromStore(
	store storev2.RootStore,
	actor []byte,
) (store.Writer, error) {
	_, state, err := store.StateLatest()
	if err != nil {
		return nil, err
	}
	reader, err := state.GetReader(actor)
	if err != nil {
		return nil, err
	}
	return cache.NewStore(reader), nil
}
