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

package logs_test

import (
	"context"
	"errors"

	"github.com/itsdevbear/bolaris/beacon/staking/logs"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

var _ logs.StakingService = &mockStakingService{}

// mockStakingService is a mock implementation of the staking service.
type mockStakingService struct {
	depositQueue    []*consensusv1.Deposit
	withdrawalQueue []*enginev1.Withdrawal
}

func (m *mockStakingService) ProcessDeposit(_ context.Context, deposit *consensusv1.Deposit) error {
	m.depositQueue = append(m.depositQueue, deposit)
	return nil
}

func (m *mockStakingService) PersistDeposits(_ context.Context) error {
	m.depositQueue = nil
	return nil
}

func (m *mockStakingService) mostRecentDeposit() (*consensusv1.Deposit, error) {
	if len(m.depositQueue) == 0 {
		return nil, errors.New("no deposits")
	}
	return m.depositQueue[len(m.depositQueue)-1], nil
}

func (m *mockStakingService) ProcessWithdrawal(_ context.Context, withdrawal *enginev1.Withdrawal) error {
	m.withdrawalQueue = append(m.withdrawalQueue, withdrawal)
	return nil
}

func (m *mockStakingService) mostRecentWithdrawal() (*enginev1.Withdrawal, error) {
	if len(m.withdrawalQueue) == 0 {
		return nil, errors.New("no withdrawals")
	}
	return m.withdrawalQueue[len(m.withdrawalQueue)-1], nil
}
