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

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
)

// verifySufficientDepth ensures that the depth is sufficient to build a tree.
func verifySufficientDepth(numLeaves int, depth uint8) error {
	switch {
	case numLeaves == 0:
		return ErrEmptyLeaves
	case depth == 0:
		return ErrZeroDepth
	case depth > MaxTreeDepth:
		return ErrExceededDepth
	case numLeaves > (1 << depth):
		return errors.Wrap(
			ErrInsufficientDepthForLeaves,
			fmt.Sprintf(
				"attempted to build tree/root with %d leaves at depth %d",
				numLeaves, depth),
		)
	}
	return nil
}
