package sszdb

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/schema"
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
](db *Backend, monolith any) (*SchemaDb[ExecutionPayloadHeaderT], error) {
	schema, err := schema.CreateSchema(monolith)
	if err != nil {
		return nil, err
	}
	return &SchemaDb[ExecutionPayloadHeaderT]{Backend: db, schemaRoot: schema}, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) getLeafBytes(path schema.ObjectPath) ([]byte, error) {
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

	return d.getNodeBytes(node.GIndex, size)
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) SetLatestExecutionPayloadHeader(
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

	return d.save(treeNode, schemaNode.GIndex)
}

type executionPayloadHeaderDenebCodec []struct{}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetLatestExecutionPayloadHeader() (
	ExecutionPayloadHeaderT, error) {
	var e ExecutionPayloadHeaderT
	path := schema.Path("LatestExecutionPayloadHeader")

	// TODO fix with codec injection?
	switch header := any(e).(type) {
	case *types.ExecutionPayloadHeaderDeneb:
		bz, err := d.getLeafBytes(path.AppendName("ParentHash"))
		if err != nil {
			return e, err
		}
		copy(header.ParentHash[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("FeeRecipient"))
		if err != nil {
			return e, err
		}
		copy(header.FeeRecipient[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("StateRoot"))
		if err != nil {
			return e, err
		}
		copy(header.StateRoot[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("ReceiptsRoot"))
		if err != nil {
			return e, err
		}
		copy(header.ReceiptsRoot[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("LogsBloom"))
		if err != nil {
			return e, err
		}
		header.LogsBloom = bz

		bz, err = d.getLeafBytes(path.AppendName("Random"))
		if err != nil {
			return e, err
		}
		copy(header.Random[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("Number"))
		if err != nil {
			return e, err
		}
		header.Number = math.U64(fastssz.UnmarshallUint64(bz))

		bz, err = d.getLeafBytes(path.AppendName("GasLimit"))
		if err != nil {
			return e, err
		}
		header.GasLimit = math.U64(fastssz.UnmarshallUint64(bz))

		bz, err = d.getLeafBytes(path.AppendName("GasUsed"))
		if err != nil {
			return e, err
		}
		header.GasUsed = math.U64(fastssz.UnmarshallUint64(bz))

		bz, err = d.getLeafBytes(path.AppendName("Timestamp"))
		if err != nil {
			return e, err
		}
		header.Timestamp = math.U64(fastssz.UnmarshallUint64(bz))

		bz, err = d.getLeafBytes(path.AppendName("ExtractData"))
		if err != nil {
			return e, err
		}
		header.ExtraData = bz

		bz, err = d.getLeafBytes(path.AppendName("BaseFeePerGas"))
		if err != nil {
			return e, err
		}
		copy(header.BaseFeePerGas[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("BlockHash"))
		if err != nil {
			return e, err
		}
		copy(header.BlockHash[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("TransactionsRoot"))
		if err != nil {
			return e, err
		}
		copy(header.TransactionsRoot[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("WithdrawalsRoot"))
		if err != nil {
			return e, err
		}
		copy(header.WithdrawalsRoot[:], bz)

		bz, err = d.getLeafBytes(path.AppendName("BlobGasUsed"))
		if err != nil {
			return e, err
		}
		header.BlobGasUsed = math.U64(fastssz.UnmarshallUint64(bz))

		bz, err = d.getLeafBytes(path.AppendName("ExcessBlobGas"))
		if err != nil {
			return e, err
		}
		header.ExcessBlobGas = math.U64(fastssz.UnmarshallUint64(bz))
	default:
		return e, errors.New("unsupported payload header type")
	}

	return e, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetGenesisValidatorsRoot() (common.Root, error) {
	path := schema.Path("GenesisValidatorsRoot")
	bz, err := d.getLeafBytes(path)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetSlot() (math.Slot, error) {
	path := schema.Path("Slot")
	n, err := d.getLeafBytes(path)
	if err != nil {
		return 0, err
	}
	slot := fastssz.UnmarshallUint64(n)
	return math.Slot(slot), nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetFork() (*types.Fork, error) {
	f := &types.Fork{}
	forkPath := schema.Path("Fork")
	bz, err := d.getLeafBytes(forkPath.AppendName("PreviousVersion"))
	if err != nil {
		return nil, err
	}
	copy(f.PreviousVersion[:], bz)

	bz, err = d.getLeafBytes(forkPath.AppendName("CurrentVersion"))
	if err != nil {
		return nil, err
	}
	copy(f.CurrentVersion[:], bz)

	bz, err = d.getLeafBytes(forkPath.AppendName("Epoch"))
	if err != nil {
		return nil, err
	}
	f.Epoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	return f, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetLatestBlockHeader() (*types.BeaconBlockHeader, error) {
	bh := &types.BeaconBlockHeader{}
	path := schema.Path("LatestBlockHeader")
	bz, err := d.getLeafBytes(path.AppendName("Slot"))
	if err != nil {
		return nil, err
	}
	bh.Slot = fastssz.UnmarshallUint64(bz)

	bz, err = d.getLeafBytes(path.AppendName("ProposerIndex"))
	if err != nil {
		return nil, err
	}
	bh.ProposerIndex = fastssz.UnmarshallUint64(bz)

	bz, err = d.getLeafBytes(path.AppendName("ParentBlockRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.ParentBlockRoot[:], bz)

	bz, err = d.getLeafBytes(path.AppendName("StateRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.StateRoot[:], bz)

	bz, err = d.getLeafBytes(path.AppendName("BodyRoot"))
	if err != nil {
		return nil, err
	}
	copy(bh.BodyRoot[:], bz)

	return bh, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetBlockRoots() ([]common.Root, error) {
	path := schema.Path("BlockRoots", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(node.GIndex, node.Size())
	if err != nil {
		return nil, err
	}

	length := fastssz.UnmarshallUint64(bz)
	roots := make([]common.Root, length)
	for i := uint64(0); i < length; i++ {
		path = schema.Path("BlockRoots").AppendIndex(i)
		bz, err = d.getLeafBytes(path)
		if err != nil {
			return nil, err
		}
		roots[i] = common.Root(bz)
	}

	return roots, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetValidatorAtIndex(index uint64) (*types.Validator, error) {
	path := schema.Path("Validators").AppendIndex(index)
	val := &types.Validator{}

	bz, err := d.getLeafBytes(path.AppendName("Pubkey"))
	if err != nil {
		return nil, err
	}
	copy(val.Pubkey[:], bz)

	bz, err = d.getLeafBytes(path.AppendName("WithdrawalCredentials"))
	if err != nil {
		return nil, err
	}
	copy(val.WithdrawalCredentials[:], bz)

	bz, err = d.getLeafBytes(path.AppendName("EffectiveBalance"))
	if err != nil {
		return nil, err
	}
	val.EffectiveBalance = math.U64(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("Slashed"))
	if err != nil {
		return nil, err
	}
	val.Slashed = fastssz.UnmarshalBool(bz)

	bz, err = d.getLeafBytes(path.AppendName("ActivationEligibilityEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEligibilityEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("ActivationEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("ExitEpoch"))
	if err != nil {
		return nil, err
	}
	val.ExitEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("WithdrawableEpoch"))
	if err != nil {
		return nil, err
	}
	val.WithdrawableEpoch = math.Epoch(fastssz.UnmarshallUint64(bz))

	return val, nil
}

func (d *SchemaDb[ExecutionPayloadHeaderT]) GetValidators() ([]*types.Validator, error) {
	path := schema.Path("Validators", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(node.GIndex, node.Size())
	if err != nil {
		return nil, err
	}

	length := fastssz.UnmarshallUint64(bz)
	validators := make([]*types.Validator, length)
	for i := uint64(0); i < length; i++ {
		val, err := d.GetValidatorAtIndex(i)
		if err != nil {
			return nil, err
		}
		validators[i] = val
	}

	return validators, nil
}
