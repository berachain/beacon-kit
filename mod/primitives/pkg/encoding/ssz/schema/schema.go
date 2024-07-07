package schema

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	uint8Size   = 1
	uint16Size  = 2
	uint32Size  = 4
	uint64Size  = 8
	uint128Size = 16
	chunkSize   = 32
)

type SSZType interface {
	Size() uint64
	Chunks() uint64

	position(p string) (uint64, uint8, error)
	child(p string) SSZType
}

// Basic Type

type basic struct {
	size uint64
}

func UInt8() SSZType   { return basic{size: uint8Size} }
func UInt16() SSZType  { return basic{size: uint16Size} }
func UInt32() SSZType  { return basic{size: uint32Size} }
func UInt64() SSZType  { return basic{size: uint64Size} }
func UInt128() SSZType { return basic{size: uint128Size} }
func UInt256() SSZType { return basic{size: chunkSize} }

func (b basic) Size() uint64 { return b.size }

func (b basic) Chunks() uint64 { return 1 }

func (b basic) child(_ string) SSZType { return b }

func (b basic) position(_ string) (uint64, uint8, error) {
	return 0, 0, errors.New("basic type has no children")
}

// Container Type

type container struct {
	Fields     []SSZType
	FieldIndex map[string]uint64
}

type field struct {
	name string
	typ  SSZType
}

func Field(name string, typ SSZType) field {
	return field{name: name, typ: typ}
}

func Container(fields ...field) SSZType {
	fieldIndex := make(map[string]uint64)
	types := make([]SSZType, len(fields))
	for i, f := range fields {
		fieldIndex[f.name] = uint64(i)
		types[i] = f.typ
	}
	return container{Fields: types, FieldIndex: fieldIndex}
}

func (c container) Size() uint64 { return chunkSize }

func (c container) Length() uint64 { return uint64(len(c.Fields)) }

func (c container) Chunks() uint64 { return uint64(len(c.Fields)) }

func (c container) child(p string) SSZType {
	return c.Fields[c.FieldIndex[p]]
}

func (c container) position(p string) (uint64, uint8, error) {
	pos, ok := c.FieldIndex[p]
	if !ok {
		return 0, 0, fmt.Errorf("field %s not found", p)
	}
	return pos, 0, nil
}

// Enumerable Type (vectors and lists)

func List(element SSZType, length uint64) SSZType {
	return enumerable{Element: element, maxLength: length}
}

func Vector(element SSZType, length uint64) SSZType {
	return enumerable{Element: element, length: length}
}

func Bytes(length uint64) SSZType {
	return Vector(UInt8(), length)
}

type enumerable struct {
	Element   SSZType
	length    uint64
	maxLength uint64
}

func (e enumerable) Size() uint64 { return chunkSize }

func (e enumerable) Chunks() uint64 {
	x := float64(e.Length()*e.Element.Size()) / chunkSize
	return uint64(math.Ceil(x))
}

func (e enumerable) child(_ string) SSZType {
	return e.Element
}

func (e enumerable) Length() uint64 {
	if e.length == 0 {
		return e.maxLength
	}
	return e.length
}

func (e enumerable) position(p string) (uint64, uint8, error) {
	i, err := strconv.ParseUint(p, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("expected index, got name %s", p)
	}
	start := i * e.Element.Size()
	return uint64(math.Floor(float64(start) / chunkSize)),
		uint8(start % chunkSize),
		nil
}

func (e enumerable) IsList() bool {
	return e.maxLength > 0
}

type Node struct {
	SSZType

	GIndex uint64
	Offset uint8
}

func Path(path string) []string {
	return strings.Split(path, "/")
}

// GetTreeNode locates a node in the SSZ merkle tree by its path and a root
// schema node to begin traversal from with gindex 1.
//
//nolint:mnd // binary math
func GetTreeNode(typ SSZType, path []string) (Node, error) {
	var (
		gindex = uint64(1)
		offset uint8
	)
	for _, p := range path {
		if p == "__len__" {
			if _, ok := typ.(enumerable); !ok {
				return Node{}, fmt.Errorf("type %T is not enumerable", typ)
			}
			gindex = 2*gindex + 1
			offset = 0
		} else {
			pos, off, err := typ.position(p)
			if err != nil {
				return Node{}, err
			}
			i := uint64(1)
			if e, ok := typ.(enumerable); ok && e.IsList() {
				i = 2
			}
			gindex = gindex*i*nextPowerOfTwo(typ.Chunks()) + pos
			typ = typ.child(p)
			offset = off
		}
	}
	return Node{SSZType: typ, GIndex: gindex, Offset: offset}, nil
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
