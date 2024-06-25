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

package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
type Transactions [][]byte

// HashTreeRoot returns the hash tree root of the Transactions list.
func (txs Transactions) HashTreeRoot() (common.Root, error) {
	var err error
	roots := make([][32]byte, len(txs))

	merkleizer := ssz.NewMerkleizer[
		common.ChainSpec, math.U64, math.U256L, [32]byte]()

	for i, tx := range txs {
		roots[i], err = merkleizer.MerkleizeByteSlice(tx)
		if err != nil {
			return common.Root{}, err
		}
	}

	roots2 := make([]ssz.Composite[
		common.ChainSpec, [32]byte], len(roots))
	for i, root := range roots {
		roots2[i] = common.Root(root)
	}

	return merkleizer.MerkleizeListComposite(
		roots2,
		constants.MaxTxsPerPayload,
	)
}
