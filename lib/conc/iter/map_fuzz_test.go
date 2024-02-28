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

package iter_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/lib/conc/iter"
	"github.com/stretchr/testify/require"
)

func FuzzMapErrNoNils(f *testing.F) {
	f.Add(uint64(0))
	f.Fuzz(func(t *testing.T, size uint64) {
		input := make([]int, size)
		for i := range input {
			input[i] = i
		}
		f := func(i *int) (*int, error) {
			if *i%2 == 0 {
				return nil, nil
			}
			return i, nil
		}
		expected := make([]*int, 0, size/2)
		for i := range input {
			if i%2 != 0 {
				expected = append(expected, &input[i])
			}
		}
		actual, err := iter.MapErrNoNils(input, f)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
