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

package merkle

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/errors"
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
