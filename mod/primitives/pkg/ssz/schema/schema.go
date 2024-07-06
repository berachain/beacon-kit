package schema

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/tree"
)

const (
	uint8Size  = 1
	uint16Size = 2
	uint32Size = 4
	uint64Size = 8
	chunkSize  = 32
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

func Basic(size uint64) basic {
	return basic{size: size}
}

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
	name  string
	typ   SSZType
	order uint64
}

func Field(order uint64, name string, typ SSZType) field {
	return field{name: name, typ: typ, order: order}
}

func Container(fields ...field) container {
	fieldIndex := make(map[string]uint64)
	var types []SSZType
	for _, f := range fields {
		fieldIndex[f.name] = uint64(f.order)
		types = append(types, f.typ)
	}
	return container{Fields: types, FieldIndex: fieldIndex}
}

func (c container) Size() uint64 { return chunkSize }

func (c container) Length() uint64 { return uint64(len(c.Fields)) }

func (c container) Chunks() uint64 { return uint64(len(c.Fields)) }

func (c container) child(
	p string,
) SSZType {
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

func List(element SSZType, length uint64) enumerable {
	return enumerable{Element: element, maxLength: length}
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

func (e enumerable) IsByteVector() bool {
	return e.Element.Size() == 1 && e.length > 0
}

func (e enumerable) IsList() bool {
	return e.maxLength > 0
}

func (e enumerable) IsFixed() bool {
	// TODO fill out cases, abstract
	_, ok := e.Element.(basic)
	return ok
}

type Node struct {
	SSZType

	GIndex uint64
	Offset uint8
}

// GetTreeNode locates a node in the SSZ merkle tree by its path and a root
// schema node to begin traversal from with gindex 1.
//
//nolint:mnd // binary math
func GetTreeNode(typ SSZType, path tree.ObjectPath) (Node, error) {
	var (
		gindex = uint64(1)
		offset uint8
	)
	for head, rest := path.Head(); head != ""; head, rest = rest.Head() {
		if head == "__len__" {
			if _, ok := typ.(enumerable); !ok {
				return Node{}, fmt.Errorf("type %T is not enumerable", typ)
			}
			gindex = 2*gindex + 1
			offset = 0
		} else {
			pos, off, err := typ.position(head)
			if err != nil {
				return Node{}, err
			}
			i := uint64(1)
			if e, ok := typ.(enumerable); ok && e.maxLength > 0 {
				// list case
				i = 2
			}
			gindex = gindex*i*nextPowerOfTwo(typ.Chunks()) + pos
			typ = typ.child(head)
			offset = off
		}
		if rest.Empty() {
			break
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
