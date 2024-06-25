package ssz

import (
	"sync"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
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
	if cap(b.Bytes) < size {
		b.Bytes = make([]common.Root, size)
	}
	b.Bytes = b.Bytes[:size]
	return b
}

// Reset resets the byte buffer.
func (b *byteBuffer) Put() {
	byteBufferPool.Put(b)
	b.Bytes = b.Bytes[:0]
}
