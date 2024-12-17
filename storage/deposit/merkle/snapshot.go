package merkle

import (
	"github.com/berachain/beacon-kit/primitives/math/pow"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/berachain/beacon-kit/primitives/merkle/zero"
)

// DepositTreeSnapshot represents the data used to create a deposit tree given a snapshot.
type DepositTreeSnapshot struct {
	finalized      [][32]byte
	depositRoot    [32]byte
	depositCount   uint64
	executionBlock executionBlock
	hasher         merkle.Hasher[[32]byte]
}

// CalculateRoot returns the root of a deposit tree snapshot.
func (ds *DepositTreeSnapshot) CalculateRoot() ([32]byte, error) {
	size := ds.depositCount
	index := len(ds.finalized)
	root := zero.Hashes[0]
	for i := uint64(0); i < DepositContractDepth; i++ {
		if (size & 1) == 1 {
			if index == 0 {
				break
			}
			index--
			root = ds.hasher.Combi(ds.finalized[index], root)
		} else {
			root = ds.hasher.Combi(root, zero.Hashes[i])
		}
		size >>= 1
	}
	return ds.hasher.MixIn(root, ds.depositCount), nil
}

// fromSnapshot returns a deposit tree from a deposit tree snapshot.
func fromSnapshot(hasher merkle.Hasher[[32]byte], snapshot DepositTreeSnapshot) (*DepositTree, error) {
	root, err := snapshot.CalculateRoot()
	if err != nil {
		return nil, err
	}
	if snapshot.depositRoot != root {
		return nil, ErrInvalidSnapshotRoot
	}
	if snapshot.depositCount >= pow.TwoToThePowerOf(DepositContractDepth) {
		return nil, ErrTooManyDeposits
	}
	tree, err := fromSnapshotParts(hasher, snapshot.finalized, snapshot.depositCount, DepositContractDepth)
	if err != nil {
		return nil, err
	}
	return &DepositTree{
		tree:                    tree,
		mixInLength:             snapshot.depositCount,
		finalizedExecutionBlock: snapshot.executionBlock,
	}, nil
}

// fromTreeParts constructs the deposit tree from pre-existing data.
func fromTreeParts(finalised [][32]byte, depositCount uint64, executionBlock executionBlock) (DepositTreeSnapshot, error) {
	snapshot := DepositTreeSnapshot{
		finalized:      finalised,
		depositRoot:    zero.Hashes[0],
		depositCount:   depositCount,
		executionBlock: executionBlock,
	}
	root, err := snapshot.CalculateRoot()
	if err != nil {
		return snapshot, ErrInvalidSnapshotRoot
	}
	snapshot.depositRoot = root
	return snapshot, nil
}
