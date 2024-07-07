package schema

import (
	"fmt"
	"math"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

// B4 creates a Vector of 4 bytes (32 bits).
func B4() SSZType {
	return Vector(U8(), 4)
}

// B8 creates a Vector of 8 bytes (64 bits).
func B8() SSZType {
	return Vector(U8(), 8)
}

// B16 creates a Vector of 16 bytes (128 bits).
func B16() SSZType {
	return Vector(U8(), 16)
}

// B20 creates a Vector of 20 bytes (160 bits).
func B20() SSZType {
	return Vector(U8(), 20)
}

// B32 creates a Vector of 32 bytes (256 bits).
func B32() SSZType {
	return Vector(U8(), 32)
}

// B48 creates a Vector of 48 bytes (384 bits).
func B48() SSZType {
	return Vector(U8(), 48)
}

// B64 creates a Vector of 64 bytes (512 bits).
func B64() SSZType {
	return Vector(U8(), 64)
}

// B96 creates a Vector of 96 bytes (768 bits).
func B96() SSZType {
	return Vector(U8(), 96)
}

// B256 creates a Vector of 256 bytes (2048 bits).
func B256() SSZType {
	return Vector(U8(), 256)
}

type vector struct {
	Element SSZType
	length  uint64
}

func Vector(element SSZType, length uint64) SSZType {
	return vector{Element: element, length: length}
}

func Bytes(length uint64) SSZType {
	return Vector(U8(), length)
}

func ByteList(length uint64) SSZType {
	return List(U8(), length)
}

func (v vector) ID() types.Type { return types.Vector }

func (v vector) ItemLength() uint64 { return chunkSize }

func (v vector) Chunks() uint64 {
	totalBytes := v.Length() * v.Element.ItemLength()
	chunks := (totalBytes + chunkSize - 1) / chunkSize
	return chunks
}

func (v vector) child(_ string) SSZType {
	return v.Element
}

func (v vector) Length() uint64 {
	return v.length
}

func (v vector) position(p string) (uint64, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * v.Element.ItemLength()
	return uint64(math.Floor(float64(start) / chunkSize)),
		uint8(start % chunkSize),
		nil
}

func (v vector) IsList() bool {
	return false
}
