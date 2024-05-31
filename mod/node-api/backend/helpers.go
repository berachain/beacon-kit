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

package backend

import (
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func resolveEpoch(stateDB StateDB, epoch string) (math.Epoch, error) {
	if epoch == "" {
		curEpoch, err := stateDB.GetEpoch()
		if err != nil {
			return math.Epoch(0), err
		}
		return curEpoch, nil
	}
	epochU64, errU64 := resolveUint64[math.Epoch](epoch)
	if errU64 != nil {
		return math.Epoch(0), errU64
	}
	return epochU64, nil
}

func resolveUint64[uint64T ~uint64](num string) (uint64T, error) {
	if num == "" {
		return uint64T(0), nil
	}
	epochU64, errU64 := strconv.ParseUint(num, 10, 64)
	if errU64 != nil {
		return uint64T(0), errU64
	}
	return uint64T(epochU64), nil
}
