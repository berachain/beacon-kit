package ssz

import (
	"errors"
	"math"
	"math/big"
)

// https://github.com/prysmaticlabs/prysm/blob/feb16ae4aaa41d9bcd066b54b779dcd38fc928d2/tools/specs-checker/data/ssz/merkle-proofs.md

// type SSZType struct {
// 	Basic
// 	Composite
// }

type GIndex big.Int

type GIndexMap map[string]big.Int
type FieldEntry struct {
	fieldName string
	fieldType string
	jsonKey   string
	gindex    big.Int
}

func CreateGIndexMap(fieldEntries []FieldEntry, depth int64) (any, error) {
	fieldGIndex := make(map[string]big.Int)
	for i := 0; i < len(fieldEntries); i++ {
		gindex, err := toGIndex(depth, *big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		fieldGIndex[fieldEntries[i].fieldName] = *gindex
	}
	return fieldGIndex, nil
}

// Return three variables:
// (i) the index of the chunk in which the given element of the item is represented;
// (ii) the starting byte position within the chunk;
// (iii) the ending byte position within the chunk.
// For example: for a 6-item list of uint64 values, index=2 will return (0, 16, 24), index=5 will return (1, 8, 16)
func GetItemPosition(typ any, indexOrVarName any) (int, int, int) {

}

func GetItemPositionContainer[B Container[RootT], RootT ~[32]byte](typ Container[RootT], indexOrVarName any) (int, int, int) {
	typ.getFieldNames
}

// Converts a path (eg. `[7, "foo", 3]` for `x[7].foo[3]`, `[12, "bar", "__len__"]` for
// `len(x[12].bar)`) into the generalized index representing its position in the Merkle tree.
func getGeneralizedIndex[B Basic[RootT], RootT ~[32]byte](typ any, path []any) (GIndex, error) {
	root := big.NewInt(1)
	for i := 0; i < len(path); i++ {
		p := path[i]
		switch typ.(type) {
		case Basic[RootT]:
			continue // If we descend to a basic type, the path cannot continue further
		default:
			pos, _, _ := GetItemPosition(typ, p)
			// baseIndex := // todo
		}
	}
}

func ConcatGIndices(gindices []big.Int) *big.Int {
	o := big.NewInt(1)
	t := new(big.Int)
	for i := 0; i < len(gindices); i++ {
		cur_gi := big.NewInt(int64(gindices[i].Int64()))
		o = t.Mul(o,
			t.Add(getPowerOfTwoFloor(cur_gi), t.Sub(cur_gi, getPowerOfTwoFloor(cur_gi))))
	}
	return o
}

// from https://github.com/ChainSafe/ssz/blob/3cc1529541990d7ac63725baa354f20a0f36f670/packages/persistent-merkle-tree/src/gindex.ts#L8
func toGIndex(depth int64, index big.Int) (*big.Int, error) {
	anchor := big.NewInt(1).Lsh(big.NewInt(1), uint(big.NewInt(depth).Uint64()))
	if anchor.Cmp(&index) == -1 {
		return big.NewInt(1), errors.New("Index too large for depth in toGIndex")
	}
	return anchor.Or(anchor, &index), nil
}

// Get the power of 2 for given input, or the closest lower power of 2 if the input is not a power of 2.
// The zero case is a placeholder and not used for math with generalized indices.
// Commonly used for "what power of two makes up the root bit of the generalized index?"
// Example: 0->1, 1->1, 2->2, 3->2, 4->4, 5->4, 6->4, 7->4, 8->8, 9->8
func getPowerOfTwoFloor(x *big.Int) *big.Int {
	one := big.NewInt(1)
	two := big.NewInt(2)
	t := new(big.Int)
	if x.Cmp(one) <= 0 {
		return one
	} else if x.Cmp(two) <= 0 {
		return x
	} else {
		half, _ := t.Div(x, two).Float64()
		floorDivHalf, _ := big.NewFloat(math.Floor(half)).Int64()
		return t.Mul(two, getPowerOfTwoFloor(big.NewInt(floorDivHalf)))
	}
}

// Get the power of 2 for given input, or the closest higher power of 2 if the input is not a power of 2.
// Commonly used for "how many nodes do I need for a bottom tree layer fitting x elements?"
// Example: 0->1, 1->1, 2->2, 3->4, 4->4, 5->8, 6->8, 7->8, 8->8, 9->16.
func getPowerOfTwoCeil(x int) int {
	if x <= 1 {
		return 1
	} else if x == 2 {
		return 2
	} else {
		return 2 * getPowerOfTwoCeil(int(math.Floor(float64(x+1)/2)))
	}
}

// Return the number of bytes in a basic type, or 32 (a full hash) for compound types.
func ItemLengthBasic[B Basic[RootT], RootT ~[32]byte](b []B) uint64 {
	return SizeOfBasic[RootT, B](b[0])
}

func ItemLength[B Composite[RootT], RootT ~[32]byte](b []B) uint64 {
	return 32
}
