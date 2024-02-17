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

package store_test

import (
	"testing"

	sdkcollections "cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	sdkruntime "github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/itsdevbear/bolaris/lib/store/collections"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper/store"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

func Test_DepositQueue(t *testing.T) {
	sk := storetypes.NewKVStoreKey("test")
	tsk := storetypes.NewTransientStoreKey("transient-test")
	kvs := sdkruntime.NewKVStoreService(sk)

	t.Run("should return correct deposits", func(t *testing.T) {
		ctx := testutil.DefaultContext(sk, tsk)
		dq := collections.NewQueue[*store.Deposit](
			sdkcollections.NewSchemaBuilder(kvs),
			"deposit_queue",
			store.DepositValue{})

		_, err := dq.Peek(ctx)
		require.Equal(t, sdkcollections.ErrNotFound, err)

		_, err = dq.Pop(ctx)
		require.Equal(t, sdkcollections.ErrNotFound, err)

		err = dq.Push(ctx, newDeposit([]byte("pubkey1"), 1))
		require.NoError(t, err)
		err = dq.Push(ctx, newDeposit([]byte("pubkey2"), 2))
		require.NoError(t, err)

		d, err := dq.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), d.GetAmount())
		require.Equal(t, []byte("pubkey1"), d.GetPubkey())

		d, err = dq.Pop(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), d.GetAmount())
		require.Equal(t, []byte("pubkey1"), d.GetPubkey())

		err = dq.Push(ctx, newDeposit([]byte("pubkey3"), 3))
		require.NoError(t, err)

		d, err = dq.Pop(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), d.GetAmount())
		require.Equal(t, []byte("pubkey2"), d.GetPubkey())

		d, err = dq.Peek(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), d.GetAmount())
		require.Equal(t, []byte("pubkey3"), d.GetPubkey())

		d, err = dq.Pop(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), d.GetAmount())
		require.Equal(t, []byte("pubkey3"), d.GetPubkey())

		_, err = dq.Peek(ctx)
		require.Equal(t, sdkcollections.ErrNotFound, err)

		_, err = dq.Pop(ctx)
		require.Equal(t, sdkcollections.ErrNotFound, err)
	})
}

func newDeposit(pubkey []byte, amount uint64) *store.Deposit {
	return &store.Deposit{
		Deposit: &consensusv1.Deposit{
			Data: &consensusv1.Deposit_Data{
				WithdrawalCredentials: common.BytesToAddress(pubkey).Bytes(),
				Pubkey:                pubkey,
				Amount:                amount,
			},
		},
	}
}
