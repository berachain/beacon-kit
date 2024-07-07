package schema

import (
	"fmt"
	"math"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

// List Type
type list struct {
	Element SSZType
	limit   uint64
}

func List(element SSZType, limit uint64) SSZType {
	return list{Element: element, limit: limit}
}

func (l list) ID() types.Type { return types.List }

func (l list) ItemLength() uint64 { return l.Element.ItemLength() }

func (l list) Chunks() uint64 {
	totalBytes := l.N() * l.Element.ItemLength()
	chunks := (totalBytes + chunkSize - 1) / chunkSize
	return chunks
}

func (l list) child(_ string) SSZType {
	return l.Element
}

func (l list) N() uint64 {
	return l.limit
}

func (l list) position(p string) (uint64, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * l.Element.ItemLength()
	return uint64(math.Floor(float64(start) / chunkSize)),
		uint8(start % chunkSize),
		nil
}

func (l list) IsList() bool {
	return true
}
