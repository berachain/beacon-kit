package sszdb

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
)

// bootstrappedKey is a temporary key to check if the database has been
// bootstrapped.  it will be made obsolete by versioning.
const bootstrappedKey = "bootstrapped"

type objectPath = merkle.ObjectPath[uint64, [32]byte]

func isByteVector(typ schema.SSZType) bool {
	return typ.ID() == schema.Vector && typ.ElementType("") == schema.U8()
}

type SchemaDB struct {
	*Backend
	schemaRoot schema.SSZType
}

func NewSchemaDB(
	backend *Backend,
	monolith treeable,
) (*SchemaDB, error) {
	schemaRoot, err := schema.Build(monolith)
	if err != nil {
		return nil, err
	}
	db := &SchemaDB{
		Backend:    backend,
		schemaRoot: schemaRoot,
	}
	return db, db.bootstrap(monolith)
}

func (db *SchemaDB) bootstrap(monolith treeable) error {
	bootstrapped, err := db.Get([]byte(bootstrappedKey))
	if err != nil {
		return err
	}
	if bootstrapped != nil {
		return nil
	}
	err = db.SaveMonolith(monolith)
	if err != nil {
		return err
	}
	return db.Set([]byte(bootstrappedKey), []byte{1})
}

type offsetBytes struct {
	bz  []byte
	idx uint32
}

func (db *SchemaDB) getLeafBytes(
	ctx context.Context,
	path objectPath,
) ([]byte, error) {
	typ, gindex, offset, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return nil, err
	}
	size := typ.ItemLength()

	// if the path was to a _byte vector_ unmarshal all of its leaves
	if isByteVector(typ) {
		_, gindex, offset, err = path.Append("0").
			GetGeneralizedIndex(db.schemaRoot)
		if err != nil {
			return nil, err
		}
		// set size to the length of byte vector
		size = typ.Length()
	}

	return db.getNodeBytes(ctx, gindex, size, offset)
}

func (db *SchemaDB) getSSZBytes(
	ctx context.Context,
	root objectPath,
) (uint32, *offsetBytes, []byte, error) {
	var (
		offsets  []*offsetBytes
		n        uint32
		sszBytes []byte
		bz       []byte
	)
	typ, _, _, err := root.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return 0, nil, nil, err
	}

	switch typ.ID() {
	case schema.Basic:
		bz, err = db.getLeafBytes(ctx, root)
		if err != nil {
			return 0, nil, nil, err
		}
		sszBytes = append(sszBytes, bz...)
		n += uint32(typ.ItemLength())
		return n, nil, sszBytes, nil
	case schema.Vector, schema.List:
		if isByteVector(typ) {
			bz, err = db.getLeafBytes(ctx, root)
			if err != nil {
				return 0, nil, nil, err
			}
			sszBytes = append(sszBytes, bz...)
			n += uint32(typ.Length())
			return n, nil, sszBytes, nil
		} else if typ.ID().IsList() {
			bz, err = db.getLeafBytes(ctx, root.Append("__len__"))
			if err != nil {
				return 0, nil, nil, err
			}
			length := fastssz.UnmarshallUint64(bz)

			var offsetBz []byte
			for i := range length {
				// TODO: list of dynamic elements not yet supported
				// result not guaranteed to be correct
				_, _, bz, err = db.getSSZBytes(ctx, root.Append(fmt.Sprintf("%d", i)))
				if err != nil {
					return 0, nil, nil, err
				}
				offsetBz = append(offsetBz, bz...)
			}
			// write empty offset address
			sszBytes = append(sszBytes, make([]byte, 4)...)
			n += 4
			return n, &offsetBytes{bz: offsetBz}, sszBytes, nil
		}
	case schema.Container:
		// TODO assumes fixed size container
		paths := make([]objectPath, typ.Length())
		for i, p := range schema.ContainerFields(typ) {
			paths[i] = root.Append(p)
		}
		for _, p := range paths {
			size, off, bz, err := db.getSSZBytes(ctx, p)
			if err != nil {
				return 0, nil, nil, err
			}
			sszBytes = append(sszBytes, bz...)
			if off != nil {
				off.idx = n
				offsets = append(offsets, off)
			}
			n += size
		}
	}

	for _, o := range offsets {
		// write offset address
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, n)
		copy(sszBytes[o.idx:], buf)
		sszBytes = append(sszBytes, o.bz...)
	}

	return n, nil, sszBytes, nil
}

// TODO clean up abstraction lines
// SchemaDB should handle object paths and bytes
// BeaconDB should handle types

func (db *SchemaDB) SetLatestExecutionPayloadHeader(
	ctx context.Context,
	header treeable,
) error {
	path := objectPath("latest_execution_payload_header")
	_, gindex, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	treeNode, err := NewTreeFromFastSSZ(header)
	if err != nil {
		return err
	}

	return db.stage(ctx, treeNode, gindex)
}

func (db *SchemaDB) GetPath(
	ctx context.Context,
	path objectPath,
) ([]byte, error) {
	_, offsetBz, bz, err := db.getSSZBytes(ctx, path)
	if offsetBz != nil {
		return offsetBz.bz, err
	}
	return bz, err
}

func (db *SchemaDB) GetSlot(ctx context.Context) (math.U64, error) {
	path := objectPath("slot")
	_, _, bz, err := db.getSSZBytes(ctx, path)
	if err != nil {
		return 0, err
	}
	return math.U64(fastssz.UnmarshallUint64(bz)), nil
}

func (db *SchemaDB) GetBlockRoots(ctx context.Context) ([]common.Root, error) {
	path := objectPath("block_roots/__len__")
	typ, gindex, offset, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return nil, err
	}
	bz, err := db.getNodeBytes(ctx, gindex, typ.ItemLength(), offset)
	if err != nil {
		return nil, err
	}

	length := fastssz.UnmarshallUint64(bz)
	sszBytes, err := db.GetPath(ctx, "block_roots")
	if err != nil {
		return nil, err
	}
	roots := make([]common.Root, length)
	n := 0
	for i := range length {
		roots[i] = common.Root(sszBytes[n : n+32])
		n += 32
	}

	return roots, nil
}
