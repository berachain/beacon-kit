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

package vm

import (
	"context"

	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"go.uber.org/zap"

	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
)

var _ snowman.Block = (*StatefulBlock)(nil)

type StatefulBlock struct {
	*block.StatelessBlock
	vm *VM
}

func (b *StatefulBlock) Verify(ctx context.Context) error {
	if err := b.vm.middleware.VerifyBlock(ctx, b.StatelessBlock); err != nil {
		return err
	}

	b.vm.verifiedBlocks[b.ID()] = b
	return nil
}

func (b *StatefulBlock) Accept(ctx context.Context) error {
	delete(b.vm.verifiedBlocks, b.ID())

	b.vm.state.SetLastAccepted(b.ID())
	b.vm.state.AddStatelessBlock(b.StatelessBlock)
	b.vm.preferredBlkID = b.ID()

	if err := b.vm.state.Commit(); err != nil {
		return err
	}

	b.vm.chainCtx.Log.Info(
		"accepted block",
		zap.Stringer("blkID", b.ID()),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	// TODO: handle dynamic validator set
	// At this stage of hooking stuff up, we consider a static validators set
	// where validators (and especially the mapping validator -> NodeID) is
	// setup in Genesis
	_, err := b.vm.middleware.AcceptBlock(ctx, b.StatelessBlock)
	return err
}

func (b *StatefulBlock) Reject(context.Context) error {
	delete(b.vm.verifiedBlocks, b.ID())
	return nil
}
