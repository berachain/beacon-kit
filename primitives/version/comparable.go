// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package version

import (
	"github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
)

/* -------------------------------------------------------------------------- */
/*                                Comparable                                  */
/* -------------------------------------------------------------------------- */

// NOTE: IsBefore and Equals are implemented to set a canonical sorting
// algorithm for the common.Version type. A standard cmp function would require
// IsBefore and Equals and could be implemented as:
//
//	func cmp(a, b common.Version) int {
//      // a is before b
//		if version.IsBefore(a, b) {
//			return -1
//		}
//      // a is the same version as b
//		if version.Equals(a, b) {
//			return 0
//		}
//      // a is after b
//		return 1
//	}

// IsBefore returns true if a is before b. This compares bytes from most significant
// to least significant in "little-endian" order.
func IsBefore(a, b common.Version) bool {
	// Iterate in order of significance.
	for i := range bytes.B4Size {
		// We short-circuit if a[i] != b[i] since we are iterating in order of significance.
		if a[i] < b[i] {
			return true
		} else if a[i] > b[i] {
			return false
		}
	}

	// If we reach this point, a and b are the same version.
	return false
}

// IsBeforeOrEquals returns true if a is before or at the same version as b.
func IsBeforeOrEquals(a, b common.Version) bool {
	return IsBefore(a, b) || Equals(a, b)
}

// Equals returns true if a and b are equal (each byte in the 4-byte vector is the same).
func Equals(a, b common.Version) bool {
	return a == b
}

// IsAfter returns true if a is after b. This compares bytes from most significant
// to least significant in "little-endian" order.
func IsAfter(a, b common.Version) bool {
	return !IsBefore(a, b) && !Equals(a, b)
}

// EqualsOrIsAfter returns true if a is the same version as b or after.
func EqualsOrIsAfter(a, b common.Version) bool {
	return !IsBefore(a, b)
}
