package engineprimitives

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/merkle"
	"github.com/karalabe/ssz"
)

type BartioTx []byte

func (tx *BartioTx) SizeSSZ() uint32 {
	return ssz.SizeDynamicBytes(*tx)
}

func (tx *BartioTx) DefineSSZ(codec *ssz.Codec) {
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineDynamicBytesOffset(
			codec,
			(*[]byte)(tx),
			constants.MaxBytesPerTx,
		)
	})
}

type Roots []common.Root

func (roots Roots) SizeSSZ() uint32 {
	return ssz.SizeSliceOfStaticBytes(roots)
}

func (roots Roots) DefineSSZ(codec *ssz.Codec) {
	codec.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticBytesContent(
			codec,
			(*[]common.Root)(&roots),
			constants.MaxTxsPerPayload,
		)
	})
	codec.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticBytesContent(
			codec,
			(*[]common.Root)(&roots),
			constants.MaxTxsPerPayload,
		)
	})
	codec.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticBytesOffset(
			codec, (*[]common.Root)(&roots), constants.MaxTxsPerPayload,
		)
	})
}

// Transactions is a typealias for [][]byte, which is how transactions are
// received in the execution payload.
//
// TODO: Remove and deprecate this type once migrated to ProperTransactions.
type BartioTransactions [][]byte

// HashTreeRoot returns the hash tree root of the Transactions list.
//
// NOTE: Uses a new merkleizer for each call.
func (txs BartioTransactions) HashTreeRoot() common.Root {
	return txs.HashTreeRootWith(
		merkle.NewMerkleizer[[32]byte, common.Root](),
	)
}

func (txs BartioTransactions) HashTreeRoot2() common.Root {
	roots := make(Roots, len(txs))
	merkleizer := merkle.NewMerkleizer[[32]byte, common.Root]()
	for i, tx := range txs {
		var err error
		roots[i], err = merkleizer.MerkleizeByteSlice(tx)
		if err != nil {
			panic(err)
		}
	}
	return ssz.HashConcurrent(roots)
}

// HashTreeRootWith returns the hash tree root of the Transactions list
// using the given merkle.
func (txs BartioTransactions) HashTreeRootWith(
	merkleizer *merkle.Merkleizer[[32]byte, common.Root],
) common.Root {
	var (
		err   error
		root  common.Root
		roots = make([]common.Root, len(txs))
	)

	for i, tx := range txs {
		roots[i], err = merkleizer.MerkleizeByteSlice(tx)
		if err != nil {
			panic(err)
		}
	}

	root, err = merkleizer.MerkleizeListComposite(
		roots,
		constants.MaxTxsPerPayload,
	)
	if err != nil {
		panic(err)
	}
	return root
}
