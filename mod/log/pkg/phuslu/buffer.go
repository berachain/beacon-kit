package phuslu

import "sync"

// byteBuffer is a byte buffer.
type byteBuffer struct {
	Bytes []byte
}

// Write writes to the byte buffer.
func (b *byteBuffer) Write(bytes []byte) (int, error) {
	b.Bytes = append(b.Bytes, bytes...)
	return len(bytes), nil
}

// byteBufferPool is a pool of byte buffers.
//
//nolint:gochecknoglobals // buffer pool
var byteBufferPool = sync.Pool{
	New: func() any {
		return new(byteBuffer)
	},
}

func resetBuffer(b *byteBuffer) {
	if b.Bytes != nil {
		b.Bytes = b.Bytes[:0]
	} else {
		b.Bytes = make([]byte, 0)
	}
}
