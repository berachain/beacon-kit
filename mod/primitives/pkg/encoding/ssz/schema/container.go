package schema

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/tree/proof"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/types/types"
)

type container struct {
	Fields     []SSZType
	FieldIndex map[string]uint64
}

func Field(name string, typ SSZType) *proof.Field[SSZType] {
	return proof.NewField(name, typ)
}

func Container(fields ...*proof.Field[SSZType]) SSZType {
	fieldIndex := make(map[string]uint64)
	types := make([]SSZType, len(fields))
	for i, f := range fields {
		fieldIndex[f.GetName()] = uint64(i)
		types[i] = f.GetValue()
	}
	return container{Fields: types, FieldIndex: fieldIndex}
}

func (c container) ID() types.Type { return types.Container }

func (c container) ItemLength() uint64 { return chunkSize }

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
