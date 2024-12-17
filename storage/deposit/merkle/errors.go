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

package merkle

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrFinalizedNodeCannotPushLeaf may occur when attempting to push a leaf
	// to a finalized node.
	// When a node is finalized, it cannot be modified or changed.
	ErrFinalizedNodeCannotPushLeaf = errors.New(
		"can't push a leaf to a finalized node",
	)

	// ErrLeafNodeCannotPushLeaf may occur when attempting to push a leaf to a
	// leaf node.
	ErrLeafNodeCannotPushLeaf = errors.New("can't push a leaf to a leaf node")

	// ErrZeroLevel occurs when the value of level is 0.
	ErrZeroLevel = errors.New("level should be greater than 0")

	// ErrZeroDepth occurs when the value of depth is 0.
	ErrZeroDepth = errors.New("depth should be greater than 0")
)

var (
	// ErrInvalidSnapshotRoot occurs when the snapshot root does not match the
	// calculated root.
	ErrInvalidSnapshotRoot = errors.New("snapshot root is invalid")

	// ErrInvalidDepositCount occurs when the value for mix in length is 0.
	ErrInvalidDepositCount = errors.New(
		"deposit count should be greater than 0",
	)

	// ErrInvalidIndex occurs when the index is less than the number of
	// finalized deposits.
	ErrInvalidIndex = errors.New(
		"index should be greater than finalizedDeposits - 1",
	)

	// ErrTooManyDeposits occurs when the number of deposits exceeds the
	// capacity of the tree.
	ErrTooManyDeposits = errors.New(
		"number of deposits should not be greater than the capacity of the tree",
	)
)
