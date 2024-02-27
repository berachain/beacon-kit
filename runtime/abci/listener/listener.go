// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package listener

import (
	"context"
	"errors"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	builder "github.com/itsdevbear/bolaris/beacon/builder/local"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// BeaconListener is the implementation of the ABCIListener interface.
type BeaconListener struct {
	chainService *blockchain.Service
	logger       log.Logger
}

// NewBeaconListener returns a new BeaconListener.
func NewBeaconListener(
	logger log.Logger,
	chainService *blockchain.Service,
) *BeaconListener {
	return &BeaconListener{
		logger:       logger,
		chainService: chainService,
	}
}

// ListenFinalizeBlock updates the streaming service with the latest
// FinalizeBlock messages
//
// TODO: Optimistic Execution chains can trigger their builder to start
// building earlier.
func (l *BeaconListener) ListenFinalizeBlock(
	ctx context.Context,
	req abci.RequestFinalizeBlock,
	_ abci.ResponseFinalizeBlock,
) error {
	// Technically speaking, there is a chance FinalizeBlock fails after this
	// call. While seemingly impossible in practice, IN THEORY the execution
	// client could end up in a bad state.
	//
	// TODO: figure out if this is a real concern or not.
	// TODO: we really should try to fork choice as soon as we have an
	// BlockHash, which is before here. This moved earlier forkchoice call
	// should 100% not be finalizing
	// the block on the EL.
	if err := l.chainService.PostFinalizeBeaconBlock(
		ctx, primitives.Slot(req.Height),
	); err != nil && !errors.Is(err, builder.ErrLocalBuildingDisabled) {
		return err
	}

	return nil
}

// ListenCommit updates the steaming service with the latest Commit messages
// and state changes.
func (*BeaconListener) ListenCommit(
	context.Context, abci.ResponseCommit, []*storetypes.StoreKVPair,
) error {
	// TODO: we can perform our block finalization here, opposed to PreBlocker()
	// if we
	// desired to.
	return nil
}
