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

package blockchain

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// RetrieveDepositsFromBlock gets logs in the Eth1 block
// received from the execution client and processes them to
// convert them into appropriate objects that can be consumed
// by other services.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) retrieveDepositsFromBlock(
	ctx context.Context,
	blockNumber math.U64,
) error {
	deposits, err := s.dc.GetDeposits(ctx, blockNumber.Unwrap())
	if err != nil {
		return err
	}

	return s.sb.DepositStore(ctx).EnqueueDeposits(deposits)
}

// PruneDepositEvents prunes deposit events.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) PruneDepositEvents(
	ctx context.Context,
	idx uint64,
) error {
	return s.sb.DepositStore(ctx).PruneToIndex(idx)
}
