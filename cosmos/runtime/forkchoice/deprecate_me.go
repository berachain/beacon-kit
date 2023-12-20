// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package forkchoice

import (
	"context"
)

func (m *Service) SyncEl(ctx context.Context) error {
	// Trigger the execution client to begin building the block, and update
	// the proposers forkchoice state accordingly.

	// TODO `block` needs to come from the latest blocked stored IAVL tree
	// // on the consensus client.
	// // block, _ := m.EngineCaller.EarliestBlock(ctx)
	// // genesisHash := m.ek.RetrieveGenesis(ctx)
	// m.logger.Info("waiting for execution client to finish sync")
	// for {
	// 	fc, err := m.getForkchoiceFromExecutionClient(ctx)
	// 	if err != nil {
	// 		m.logger.Error("failed to get forkchoice state", "err", err)
	// 		return err
	// 	}
	// 	payloadID, lastValidHash, err := m.EngineCaller.ForkchoiceUpdated(ctx, fc, payloadattribute.EmptyWithVersion(3))
	// 	if err == nil {
	// 		break
	// 	}

	// 	m.logger.Info("waiting for execution client to sync", "error", err)
	// 	m.logger.Info("waiting for execution client to sync", "payloadID", payloadID, "lastValidHash", lastValidHash)
	// 	time.Sleep(1 * time.Second)
	// }

	return nil
}
