package sszdb

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
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
	monolith Treeable,
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

func (db *SchemaDB) bootstrap(monolith Treeable) error {
	bootstrapped, err := db.get([]byte(bootstrappedKey))
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
	header Treeable,
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

func (db *SchemaDB) GetObject(
	ctx context.Context,
	path objectPath,
	obj constraints.SSZUnmarshaler,
) error {
	bz, err := db.GetPath(ctx, path)
	if err != nil {
		return err
	}
	return obj.UnmarshalSSZ(bz)
}

func (db *SchemaDB) SetObject(
	ctx context.Context,
	path objectPath,
	obj Treeable,
) error {
	treeNode, err := NewTreeFromFastSSZ(obj)
	if err != nil {
		return err
	}
	_, gidx, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	return db.stage(ctx, treeNode, gidx)
}

func (db *SchemaDB) GetSlot(ctx context.Context) (math.U64, error) {
	path := objectPath("slot")
	_, _, bz, err := db.getSSZBytes(ctx, path)
	if err != nil {
		return 0, err
	}
	return math.U64(fastssz.UnmarshallUint64(bz)), nil
}

func (db *SchemaDB) getListLength(
	ctx context.Context,
	path string,
) (uint64, error) {
	op := objectPath(path + "/__len__")
	bz, err := db.GetPath(ctx, op)
	if err != nil {
		return 0, err
	}
	return fastssz.UnmarshallUint64(bz), nil
}

func (db *SchemaDB) setListLength(
	ctx context.Context,
	path string,
	length uint64,
) error {
	op := objectPath(path + "/__len__")
	_, gindex, _, err := op.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	val := make([]byte, 32)
	binary.LittleEndian.PutUint64(val, length)
	return db.stage(
		ctx,
		&Node{Value: val},
		gindex,
	)
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

func (db *SchemaDB) SetBlockRoots(
	ctx context.Context,
	roots []common.Root,
) error {
	path := objectPath("block_roots")
	typ, gindex, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	if typ.ID() != schema.List {
		return fmt.Errorf("expected list type, got %d", typ.ID())
	}
	if uint64(len(roots)) > typ.Length() {
		return fmt.Errorf(
			"expected max %d roots, got %d",
			typ.Length(),
			len(roots),
		)
	}
	// use fastssz to produce a tree
	hh := &fastssz.Wrapper{}
	for _, root := range roots {
		hh.Append(root[:])
	}
	hh.MerkleizeWithMixin(0, uint64(len(roots)), typ.Length())
	node := copyTree(hh.Node())
	node.CachedHash()
	return db.stage(ctx, node, gindex)
}

func (db *SchemaDB) GetBlockRootAtIndex(
	ctx context.Context,
	index uint64,
) (common.Root, error) {
	path := objectPath(fmt.Sprintf("block_roots/%d", index))
	bz, err := db.GetPath(ctx, path)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

func (db *SchemaDB) SetBlockRootAtIndex(
	ctx context.Context,
	index uint64,
	root common.Root,
) error {
	length, err := db.getListLength(ctx, "block_roots")
	if err != nil {
		return err
	}
	if index > length {
		return fmt.Errorf("index %d out of bounds; len=%d", index, length)
	}
	path := objectPath(fmt.Sprintf("block_roots/%d", index))
	_, gidx, _, err := path.GetGeneralizedIndex(db.schemaRoot)
	if err != nil {
		return err
	}
	err = db.stage(ctx, &Node{Value: root[:]}, gidx)
	if err != nil {
		return err
	}

	// when the index is at the end of the list, we need to update the length
	// and potentially add some zero hashes
	if index == length {
		gindex := gidx
		depth := 0
		for gindex > 1 {
			if gindex%2 == 0 {
				sibling, err := db.getNode(ctx, gindex+1)
				if err != nil {
					return err
				}
				if sibling != nil {
					// exit condition: once pre-existing sibling is found
					// upward traversal can be stopped
					break
				}
				db.stages[gindex+1] = db.zeroHashes[depth]
			}
			depth += 1
			gindex /= 2
		}
		return db.setListLength(ctx, "block_roots", index+1)
	}
	return nil
}
