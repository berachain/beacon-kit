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

// GetLastProcessedBlock retrieves the last processed block from the store.
func (s *Store) GetLastProcessedBlock(logSig common.Hash) (uint64, error) {
	return s.logLastProcessedBlock.Get(s.ctx, logSig)
}

// SetLastProcessedBlock sets the last processed block in the store.
func (s *Store) SetLastProcessedBlock(logSig common.Hash, blockNum uint64) error {
	return s.logLastProcessedBlock.Set(s.ctx, logSig, blockNum)
}

// GetLastProcessedIndex retrieves the last processed index from the store.
func (s *Store) GetLastProcessedIndex(logSig common.Hash) (uint64, error) {
	return s.logLastProcessedIndex.Get(s.ctx, logSig)
}

// SetLastProcessedIndex sets the last processed index in the store.
func (s *Store) SetLastProcessedIndex(logSig common.Hash, index uint64) error {
	return s.logLastProcessedIndex.Set(s.ctx, logSig, index)
}
