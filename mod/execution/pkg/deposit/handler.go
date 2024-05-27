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

package deposit

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
)

// TODO: eth1FollowDistance should be done actually properly
const eth1FollowDistance = 1

func (s *Service[DepositStoreT]) handleDepositEvent(e types.BlockEvent) error {
	slot := e.Block().GetSlot() - eth1FollowDistance
	s.logger.Info("💵 processing deposit logs 💵", "slot", slot)
	deposits, err := s.dc.GetDeposits(e.Context(), slot.Unwrap())
	if err != nil {
		return err
	}

	if err = s.sb.DepositStore(e.Context()).EnqueueDeposits(
		deposits,
	); err != nil {
		return err
	}
	return nil
}
