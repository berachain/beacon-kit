package proof

import (
	"fmt"
	"strings"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// ObjectPath represents a path to an object in a Merkle tree.
type ObjectPath string

// Split returns the path split by "/"
func (p ObjectPath) Split() []string {
	return strings.Split(string(p), "/")
}

// GetElemType returns the type of the element of an object of the given type with the given index
// or member variable name (e.g. 7 for x[7], "foo" for x.foo)
func GetElemType(element SSZType, indexOrVariableName interface{}) (SSZType, error) {
	switch t := element.Type(); t {
	case types.Container:
		variableName, ok := indexOrVariableName.(string)
		if !ok {
			return nil, fmt.Errorf("expected string variable name for Container type, got %T", indexOrVariableName)
		}
		return element.(Container[SSZType]).GetFieldByName(variableName), nil
	case types.Elements:
		return element.(Elements), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", element)
	}
}

// GetItemPosition returns the position of an item within an SSZ type.
// It returns three values:
// 1. The index of the chunk in which the given element of the item is represented
// 2. The starting byte position within the chunk
// 3. The ending byte position within the chunk
// For example: for a 6-item list of uint64 values, index=2 will return (0, 16, 24), index=5 will return (1, 8, 16)
func GetItemPosition(typ SSZType, indexOrVariableName interface{}) (uint64, uint64, uint64, error) {
	switch typ.Type() {
	case types.Elements:
		index, ok := indexOrVariableName.(uint64)
		if !ok {
			return 0, 0, 0, fmt.Errorf("expected uint64 index for Elements type, got %T", indexOrVariableName)
		}
		start := index * ItemLength(typ)
		return start / 32, start % 32, start%32 + ItemLength(typ), nil
	case types.Container:
		variableName, ok := indexOrVariableName.(string)
		if !ok {
			return 0, 0, 0, fmt.Errorf("expected string variable name for Container type, got %T", indexOrVariableName)
		}

		t := typ.(Container[SSZType])
		return t.GetFieldIndex(variableName), 0, ItemLength(typ), nil
	default:
		return 0, 0, 0, fmt.Errorf("only lists/vectors/containers supported, got %T", typ)
	}
}

// GetGeneralizedIndex converts a path (e.g. [7, "foo", 3] for x[7].foo[3], [12, "bar", "__len__"] for
// len(x[12].bar)) into the generalized index representing its position in the Merkle tree.
func GetGeneralizedIndex(typ SSZType, path ...interface{}) (uint64, error) {
	root := uint64(1)
	for _, p := range path {
		if typ.Type() == types.Basic {
			return 0, fmt.Errorf("cannot descend further from a basic type")
		}

		if p == "__len__" {
			if typ.Type() != types.Elements {
				return 0, fmt.Errorf("__len__ is only valid for List or ByteList types")
			}
			typ = math.U64(0)
			root = root*2 + 1
		} else {
			pos, _, _, err := GetItemPosition(typ, p)
			if err != nil {
				return 0, err
			}

			baseIndex := uint64(1)
			if typ.Type() == types.Elements {
				baseIndex = 2
			}

			chunkCount := typ.(types.SSZEnumerable[SSZType]).ChunkCount()
			root = root*baseIndex*nextPowerOfTwo(chunkCount) + pos

			elemType, err := GetElemType(typ, p)
			if err != nil {
				return 0, err
			}
			typ = elemType
		}
	}
	return root, nil
}

//nolint:mnd // binary math
func nextPowerOfTwo(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
