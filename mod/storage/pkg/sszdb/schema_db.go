package sszdb

import (
	"context"
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
)

type SchemaDb[
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		ssz.SSZTreeable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
] struct {
	*Backend
	schemaRoot schema.SSZType
}

func NewSchemaDb[
	ExecutionPayloadHeaderT interface {
		constraints.SSZMarshallable
		ssz.SSZTreeable
		NewFromSSZ([]byte, uint32) (ExecutionPayloadHeaderT, error)
		Version() uint32
	},
](db *Backend, monolith ssz.SSZTreeable) (
	*SchemaDb[ExecutionPayloadHeaderT], error,
) {
	schema, err := schema.CreateSchema(monolith)
	if err != nil {
		return nil, err
	}
	schemaDB := &SchemaDb[ExecutionPayloadHeaderT]{Backend: db, schemaRoot: schema}
	if err = schemaDB.bootstrap(monolith); err != nil {
		return nil, err
	}
	return schemaDB, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) getLeafBytes(
	ctx context.Context,
	path schema.ObjectPath,
) ([]byte, error) {
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	size := node.Size()

	// if the path was to a _byte vector_ unmarshal all of its leaves
	if en, ok := node.SSZType.(schema.Enumerable); ok && en.IsByteVector() {
		node, err = schema.GetTreeNode(d.schemaRoot, path.AppendIndex(0))
		if err != nil {
			return nil, err
		}
		// set size to the length of byte vector
		size = en.Length()
	}

	return d.getNodeBytes(ctx, node.GIndex, size, node.Offset)
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) bootstrap(
	monolith ssz.SSZTreeable,
) error {
	bootstrapped, err := d.Get([]byte("bootstrapped"))
	if err != nil {
		return err
	}
	if bootstrapped != nil {
		return nil
	}
	err = d.SaveMonolith(monolith)
	if err != nil {
		return err
	}
	return d.Set([]byte("bootstrapped"), []byte{1})
}

type offsetBytes struct {
	bz  []byte
	idx uint32
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) getSSZBytes(
	ctx context.Context,
	root schema.ObjectPath,
) (uint32, *offsetBytes, []byte, error) {
	var (
		offsets  []*offsetBytes
		n        uint32
		sszBytes []byte
		bz       []byte
	)
	rootNode, err := schema.GetTreeNode(d.schemaRoot, root)
	if err != nil {
		return 0, nil, nil, err
	}

	switch typ := rootNode.SSZType.(type) {
	case schema.Basic:
		bz, err = d.getLeafBytes(ctx, root)
		if err != nil {
			return 0, nil, nil, err
		}
		sszBytes = append(sszBytes, bz...)
		// TODO remove type cast with refactor
		n += uint32(typ.Size())
		return n, nil, sszBytes, nil
	case schema.Enumerable:
		if typ.IsByteVector() {
			bz, err = d.getLeafBytes(ctx, root)
			if err != nil {
				return 0, nil, nil, err
			}
			sszBytes = append(sszBytes, bz...)
			n += uint32(typ.Length())
			return n, nil, sszBytes, nil
		} else if typ.IsList() {
			bz, err = d.getLeafBytes(ctx, root.AppendName("__len__"))
			if err != nil {
				return 0, nil, nil, err
			}
			length := fastssz.UnmarshallUint64(bz)

			var offsetBz []byte
			for i := range length {
				// list of dynamic elements not yet supported
				_, _, bz, err = d.getSSZBytes(ctx, root.AppendIndex(i))
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
		paths := make([]schema.ObjectPath, len(typ.Fields))
		for p, i := range typ.FieldIndex {
			paths[i] = root.AppendName(p)
		}
		for _, p := range paths {
			size, off, bz, err := d.getSSZBytes(ctx, p)
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

func (d *SchemaDb[ExecutionPayloadHeaderT]) SetLatestExecutionPayloadHeader(
	ctx context.Context,
	header ExecutionPayloadHeaderT,
) error {
	schemaNode, err := schema.GetTreeNode(
		d.schemaRoot,
		schema.Path("LatestExecutionPayloadHeader"))
	if err != nil {
		return err
	}
	treeNode, err := header.GetRootNode()
	if err != nil {
		return err
	}

	return d.stage(ctx, treeNode, schemaNode.GIndex)
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetLatestExecutionPayloadHeader(
	ctx context.Context,
) (
	ExecutionPayloadHeaderT, error,
) {
	var e ExecutionPayloadHeaderT
	path := schema.Path("LatestExecutionPayloadHeader")
	_, _, bz, err := d.getSSZBytes(ctx, path)
	if err != nil {
		return e, err
	}
	return e.NewFromSSZ(bz, version.Deneb)
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetGenesisValidatorsRoot(
	ctx context.Context,
) (
	common.Root, error,
) {
	path := schema.Path("GenesisValidatorsRoot")
	bz, err := d.getLeafBytes(ctx, path)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetSlot(
	ctx context.Context,
) (math.Slot, error) {
	path := schema.Path("Slot")
	n, err := d.getLeafBytes(ctx, path)
	if err != nil {
		return 0, err
	}
	slot := fastssz.UnmarshallUint64(n)
	return math.Slot(slot), nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetFork(
	ctx context.Context,
) (*types.Fork, error) {
	f := &types.Fork{}
	forkPath := schema.Path("Fork")
	bz, err := d.getLeafBytes(ctx, forkPath.AppendName("PreviousVersion"))
	if err != nil {
		return nil, err
	}
	copy(f.PreviousVersion[:], bz)

	bz, err = d.getLeafBytes(ctx, forkPath.AppendName("CurrentVersion"))
	if err != nil {
		return nil, err
	}
	copy(f.CurrentVersion[:], bz)

	bz, err = d.getLeafBytes(ctx, forkPath.AppendName("Epoch"))
	if err != nil {
		return nil, err
	}
	f.Epoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	return f, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetLatestBlockHeader(
	ctx context.Context,
) (
	*types.BeaconBlockHeader, error,
) {
	bh := &types.BeaconBlockHeader{}
	path := schema.Path("LatestBlockHeader")
	bz, err := d.getLeafBytes(ctx, path.AppendName("Slot"))
	if err != nil {
		return nil, err
	}
	bh.Slot = fastssz.UnmarshallUint64(bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("ProposerIndex"))
	if err != nil {
		return nil, err
	}
	bh.ProposerIndex = fastssz.UnmarshallUint64(bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("ParentBlockRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.ParentBlockRoot[:], bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("StateRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.StateRoot[:], bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("BodyRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.BodyRoot[:], bz)

	return bh, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetBlockRoots(
	ctx context.Context,
) (
	[]common.Root, error,
) {
	path := schema.Path("BlockRoots", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(ctx, node.GIndex, node.Size(), node.Offset)
	if err != nil {
		return nil, err
	}

	length := fastssz.UnmarshallUint64(bz)
	roots := make([]common.Root, length)
	for i := range length {
		path = schema.Path("BlockRoots").AppendIndex(i)
		bz, err = d.getLeafBytes(ctx, path)
		if err != nil {
			return nil, err
		}
		roots[i] = common.Root(bz)
	}

	return roots, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetValidatorAtIndex(
	ctx context.Context,
	index uint64,
) (*types.Validator, error) {
	path := schema.Path("Validators").AppendIndex(index)
	val := &types.Validator{}

	bz, err := d.getLeafBytes(ctx, path.AppendName("Pubkey"))
	if err != nil {
		return nil, err
	}
	copy(val.Pubkey[:], bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("WithdrawalCredentials"))
	if err != nil {
		return nil, err
	}
	copy(val.WithdrawalCredentials[:], bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("EffectiveBalance"))
	if err != nil {
		return nil, err
	}
	val.EffectiveBalance = math.U64(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(ctx, path.AppendName("Slashed"))
	if err != nil {
		return nil, err
	}
	val.Slashed = fastssz.UnmarshalBool(bz)

	bz, err = d.getLeafBytes(ctx, path.AppendName("ActivationEligibilityEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEligibilityEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(ctx, path.AppendName("ActivationEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(ctx, path.AppendName("ExitEpoch"))
	if err != nil {
		return nil, err
	}
	val.ExitEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(ctx, path.AppendName("WithdrawableEpoch"))
	if err != nil {
		return nil, err
	}
	val.WithdrawableEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	return val, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetValidators(
	ctx context.Context,
) (
	[]*types.Validator, error,
) {
	path := schema.Path("Validators", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(ctx, node.GIndex, node.Size(), node.Offset)
	if err != nil {
		return nil, err
	}

	length := fastssz.UnmarshallUint64(bz)
	validators := make([]*types.Validator, length)
	for i := range length {
		validators[i], err = d.GetValidatorAtIndex(ctx, i)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
