package sszdb

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/schema"
	ssz "github.com/ferranbt/fastssz"
)

type SchemaDb struct {
	*DB
	schemaRoot schema.SSZType
}

func NewSchemaDb(db *DB, monolith any) (*SchemaDb, error) {
	schema, err := schema.CreateSchema(monolith)
	if err != nil {
		return nil, err
	}
	return &SchemaDb{DB: db, schemaRoot: schema}, nil
}

func (d *SchemaDb) getLeafBytes(path schema.ObjectPath) ([]byte, error) {
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

func (d *SchemaDb) GetGenesisValidatorsRoot() (common.Root, error) {
	path := schema.Path("GenesisValidatorsRoot")
	bz, err := d.getLeafBytes(path)
	if err != nil {
		return common.Root{}, err
	}
	return common.Root(bz), nil
}

func (d *SchemaDb) GetSlot() (math.Slot, error) {
	path := schema.Path("Slot")
	n, err := d.getLeafBytes(path)
	if err != nil {
		return 0, err
	}
	slot := ssz.UnmarshallUint64(n)
	return math.Slot(slot), nil
}

func (d *SchemaDb) GetFork() (*types.Fork, error) {
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
	f.Epoch = math.Epoch(ssz.UnmarshallUint64(bz))

	return f, nil
}

func (d *SchemaDb) GetLatestBlockHeader() (*types.BeaconBlockHeader, error) {
	bh := &types.BeaconBlockHeader{}
	path := schema.Path("LatestBlockHeader")
	bz, err := d.getLeafBytes(path.AppendName("Slot"))
	if err != nil {
		return nil, err
	}
	bh.Slot = ssz.UnmarshallUint64(bz)

	bz, err = d.getLeafBytes(path.AppendName("ProposerIndex"))
	if err != nil {
		return nil, err
	}
	bh.ProposerIndex = ssz.UnmarshallUint64(bz)

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

func (d *SchemaDb) GetBlockRoots() ([]common.Root, error) {
	path := schema.Path("BlockRoots", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(node.GIndex, node.Size())
	if err != nil {
		return nil, err
	}

	length := ssz.UnmarshallUint64(bz)
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

func (d *SchemaDb) GetValidatorAtIndex(index uint64) (*types.Validator, error) {
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
	val.EffectiveBalance = math.U64(ssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("Slashed"))
	if err != nil {
		return nil, err
	}
	val.Slashed = ssz.UnmarshalBool(bz)

	bz, err = d.getLeafBytes(path.AppendName("ActivationEligibilityEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEligibilityEpoch = math.Epoch(ssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("ActivationEpoch"))
	if err != nil {
		return nil, err
	}
	val.ActivationEpoch = math.Epoch(ssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("ExitEpoch"))
	if err != nil {
		return nil, err
	}
	val.ExitEpoch = math.Epoch(ssz.UnmarshallUint64(bz))

	bz, err = d.getLeafBytes(path.AppendName("WithdrawableEpoch"))
	if err != nil {
		return nil, err
	}
	val.WithdrawableEpoch = math.Epoch(ssz.UnmarshallUint64(bz))

	return val, nil
}

func (d *SchemaDb) GetValidators() ([]*types.Validator, error) {
	path := schema.Path("Validators", "__len__")
	node, err := schema.GetTreeNode(d.schemaRoot, path)
	if err != nil {
		return nil, err
	}
	bz, err := d.getNodeBytes(node.GIndex, node.Size())
	if err != nil {
		return nil, err
	}

	length := ssz.UnmarshallUint64(bz)
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
