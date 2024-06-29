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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestTransactions(t *testing.T) {
	txs := engineprimitives.TransactionsFromBytes(
		[][]byte{[]byte("transaction1"),
			[]byte("transaction2"),
			[]byte("transaction3")},
	)

	root, err := txs.HashTreeRoot()
	require.NoError(t, err)
	require.NotNil(t, root)

	require.NotEqual(t, common.Root{}, root)

	// Create two identical Transactions
	txs1 := engineprimitives.TransactionsFromBytes(
		[][]byte{[]byte("transaction1"),
			[]byte("transaction2"),
			[]byte("transaction3")},
	)
	txs2 := engineprimitives.Transactions2FromBytes(
		[][]byte{[]byte("transaction1"),
			[]byte("transaction2"),
			[]byte("transaction3")},
	)

	// Calculate HashTreeRoot for both
	root1, err1 := txs1.HashTreeRoot()
	require.NoError(t, err1)
	root2, err2 := txs2.HashTreeRoot()
	require.NoError(t, err2)

	// Check if the roots are equal
	require.Equal(t, root1, root2, "HashTreeRoot of identical Transactions should be equal")

}
