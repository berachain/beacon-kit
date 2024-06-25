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

package engineprimitives

import (
	"sync"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
)

// byteBuffer is a byte buffer.
type byteBuffer struct {
	Bytes []common.Root
}

// byteBufferPool is a pool of byte buffers.
//
//nolint:gochecknoglobals // buffer pool
var byteBufferPool = sync.Pool{
	New: func() any {
		return &byteBuffer{
			//nolint:mnd // reasonable number of bytes
			Bytes: make([]common.Root, 0, 256),
		}
	},
}

// getBytes retrieves a byte buffer from the pool.
func getBytes(size int) *byteBuffer {
	//nolint:errcheck // its okay.
	b := byteBufferPool.Get().(*byteBuffer)
	if b == nil {
		b = &byteBuffer{
			Bytes: make([]common.Root, size),
		}
	} else {
		if b.Bytes == nil || cap(b.Bytes) < size {
			b.Bytes = make([]common.Root, size)
		}
		b.Bytes = b.Bytes[:size]
	}
	return b
}

// Reset resets the byte buffer.
func (b *byteBuffer) Reset() {
	b.Bytes = b.Bytes[:0]
}

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
type Transactions [][]byte

// HashTreeRoot returns the hash tree root of the Transactions list.
func (txs Transactions) HashTreeRoot() (common.Root, error) {
	var err error
	roots := getBytes(len(txs))
	defer byteBufferPool.Put(roots)

	// Ensure roots.Bytes is not nil
	if roots.Bytes == nil {
		return common.Root{}, errors.New("failed to allocate byte buffer")
	}

	for i, tx := range txs {
		roots.Bytes[i], err = ssz.MerkleizeByteSlice[math.U64, common.Root](tx)
		if err != nil {
			return common.Root{}, err
		}
	}

	return ssz.MerkleizeListComposite[any, math.U64](
		roots.Bytes, constants.MaxTxsPerPayload,
	)
}
