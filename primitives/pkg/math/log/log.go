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

package log

import "math/bits"

// ILog2Ceil returns the ceiling of the base 2 logarithm of the input.
func ILog2Ceil[U64T ~uint64](u U64T) uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return uint8(bits.Len64(uint64(u - 1)))
}

// ILog2Floor returns the floor of the base 2 logarithm of the input.
func ILog2Floor[U64T ~uint64](u U64T) uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return uint8(bits.Len64(uint64(u))) - 1
}
