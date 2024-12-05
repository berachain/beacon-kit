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
	// ErrNegativeIndex indicates that a negative index was provided.
	ErrNegativeIndex = errors.New("negative index provided")

	// ErrEmptyLeaves indicates that no items were provided to generate a Merkle
	// tree.
	ErrEmptyLeaves = errors.New("no items provided to generate Merkle tree")

	// ErrInsufficientDepthForLeaves indicates that the depth provided for the
	// Merkle tree is insufficient to store the provided leaves.
	ErrInsufficientDepthForLeaves = errors.New(
		"insufficient depth to store leaves",
	)

	// ErrZeroDepth indicates that the depth provided for the Merkle tree is
	// zero, which is invalid.
	ErrZeroDepth = errors.New("depth must be greater than 0")

	// ErrDepthExceedsLimitDepth indicates that the depth provided for the
	// Merkle
	// tree exceeds the specified limit depth.
	ErrDepthExceedsLimitDepth = errors.New(
		"depth exceeds the specified limit depth",
	)

	// ErrExceededDepth indicates that the provided depth exceeds the supported
	// maximum depth for a Merkle tree.
	ErrExceededDepth = errors.New("supported merkle tree depth exceeded")

	// ErrOddLengthTreeRoots is returned when the input list length must be
	// even.
	ErrOddLengthTreeRoots = errors.New("input list length must be even")

	// ErrMaxRootsExceeded is returned when the number of roots exceeds the
	// maximum allowed.
	ErrMaxRootsExceeded = errors.New(
		"number of roots exceeds the maximum allowed",
	)

	// ErrLeavesExceedsLimit is returned when the number of leaves exceeds the
	// maximum allowed.
	ErrLeavesExceedsLimit = errors.New(
		"number of leaves exceeds the maximum allowed",
	)
)
