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

package beacon

import (
	"github.com/ethereum/go-ethereum/common"
)

// SetParentBlockRoot sets the parent block root in the BeaconStore.
// It panics if there is an error setting the parent block root.
func (s *Store) SetParentBlockRoot(parentRoot [32]byte) {
	if err := s.parentBlockRoot.Set(s.ctx, parentRoot[:]); err != nil {
		panic(err)
	}
}

// GetParentBlockRoot retrieves the parent block root from the BeaconStore.
// It returns an empty hash if there is an error retrieving the parent block
// root.
func (s *Store) GetParentBlockRoot() [32]byte {
	parentRoot, err := s.parentBlockRoot.Get(s.ctx)
	if err != nil {
		parentRoot = []byte{}
	}
	return common.BytesToHash(parentRoot)
}
