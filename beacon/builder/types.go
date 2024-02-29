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

package builder

import (
	"context"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
	GetBestPayload(
		ctx context.Context,
		slot primitives.Slot,
		parentBlockRoot [32]byte,
		parentEth1Hash common.Hash,
	) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error)
}

// ExecutionService represents a service that is responsible for
// processing logs in an eth1 block.
type ExecutionService interface {
	// ProcessLogsInETH1Block processes logs in an eth1 block
	// by unmarshaling them into appropriate values.
	ProcessLogsInETH1Block(
		ctx context.Context, blkHash common.Hash,
	) ([]*reflect.Value, error)
}

// StakingService represents a service that is responsible for
// staking-related operations.
type StakingService interface {
	// AcceptDepositsIntoQueue accepts a list
	// of deposits into the deposit queue.
	AcceptDepositsIntoQueue(
		ctx context.Context,
		deposits []*consensusv1.Deposit,
	) error

	// DequeueDeposits dequeues deposits from the deposit queue.
	DequeueDeposits(
		ctx context.Context,
	) ([]*consensusv1.Deposit, error)
}
