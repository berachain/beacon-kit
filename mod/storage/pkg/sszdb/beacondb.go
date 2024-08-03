package sszdb

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	fastssz "github.com/ferranbt/fastssz"
)

// bootstrappedKey is a temporary key to check if the database has been
// bootstrapped.  it will be made obsolete by versioning.
const bootstrappedKey = "bootstrapped"

type objectPath = merkle.ObjectPath[uint64, [32]byte]

func isByteVector(typ schema.SSZType) bool {
	return typ.ID() == schema.Vector && typ.ElementType("") == schema.U8()
}

type BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT interface {
		fastssz.HashRoot
		constraints.SSZMarshallable
	},
	// ForkT,
	// ValidatorT beacondb.Validator,
	// ValidatorsT ~[]ValidatorT,
] struct {
	*Backend
	schemaRoot schema.SSZType
}

func NewBeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT interface {
		fastssz.HashRoot
		constraints.SSZMarshallable
	},
	// ForkT,
	// ValidatorT beacondb.Validator,
	// ValidatorsT ~[]ValidatorT,
](
	backend *Backend,
	monolith fastssz.HashRoot,
) (*BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,

// ForkT,
// ValidatorT,
// ValidatorsT,
], error) {
	schemaRoot, err := CreateSchema(monolith)
	if err != nil {
		return nil, err
	}
	db := &BeaconStateDB[
		BeaconBlockHeaderT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
	// ForkT,
	// ValidatorT,
	// ValidatorsT,
	]{
		Backend:    backend,
		schemaRoot: schemaRoot,
	}
	return db, db.bootstrap(monolith)
}

func (db *BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,

// ForkT,
// ValidatorT,
// ValidatorsT,
]) bootstrap(
	monolith fastssz.HashRoot,
) error {
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

func (db *BeaconStateDB[_, _, _ /*_, _, _*/]) getLeafBytes(
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
		typ, gindex, offset, err = path.Append("0").
			GetGeneralizedIndex(db.schemaRoot)
		if err != nil {
			return nil, err
		}
		// set size to the length of byte vector
		size = typ.Length()
	}

	return db.getNodeBytes(ctx, gindex, size, offset)
}

func (db *BeaconStateDB[_, _, _ /*_, _, _*/]) getSSZBytes(
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
	case schema.Vector:
	case schema.List:
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

func (db *BeaconStateDB[
	BeaconBlockHeaderT,
	Eth1DataT,
	ExecutionPayloadHeaderT,

// ForkT,
// ValidatorT,
// ValidatorsT,
]) SetLatestExecutionPayloadHeader(
	ctx context.Context,
	header ExecutionPayloadHeaderT,
) error {
	path := objectPath("LatestExecutionPayloadHeader")
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

func (db *BeaconStateDB[
	_, _, ExecutionPayloadHeaderT, /*_, _, _,*/
]) GetLatestExecutionPayloadHeader(
	ctx context.Context,
) (ExecutionPayloadHeaderT, error) {
	var e ExecutionPayloadHeaderT
	path := objectPath("LatestExecutionPayloadHeader")
	_, _, bz, err := db.getSSZBytes(ctx, path)
	if err != nil {
		return e, err
	}
	return e, e.UnmarshalSSZ(bz)
}
