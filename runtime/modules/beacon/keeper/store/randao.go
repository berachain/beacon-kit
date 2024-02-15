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

package store

import (
	"github.com/itsdevbear/bolaris/beacon/core/randao"
)

func (s *BeaconStore) GetRandaoReveals() []randao.Reveal {
	ln := s.randaoReveals.Size()
	ret := make([]randao.Reveal, ln)
	it := s.randaoReveals.Iterator()

	idx := 0
	for it.Next() {
		ret[idx] = it.Value().(randao.Reveal)
		idx++
	}

	return ret
}

func (s *BeaconStore) GetLatestRandaoReveal() *randao.Reveal {
	it := s.randaoReveals.Iterator()
	if it.Last() {
		ret := it.Value().(randao.Reveal)
		return &ret
	}

	return nil
}
