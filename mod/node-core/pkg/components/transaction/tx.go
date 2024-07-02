package transaction

import (
	"cosmossdk.io/core/transaction"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/cosmos/gogoproto/proto"
)

var _ transaction.Tx = (*SSZTx)(nil)

type SSZTx struct {
	constraints.SSZMarshallable
}

func NewTxFromSSZ[T transaction.Tx](bz []byte) T {
	tx := &SSZTx{}
	if err := tx.UnmarshalSSZ(bz); err != nil {
		panic(err)
	}

	return any(tx).(T)
}

// func (tx *SSZTx) Unwrap() *SSZTx {
// 	return tx
// }

func (tx *SSZTx) Bytes() []byte {
	bz, err := tx.MarshalSSZ()
	if err != nil {
		panic(err)
	}
	return bz
}

func (tx *SSZTx) Hash() [32]byte {
	bz, err := tx.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	return bz
}

func (tx *SSZTx) GetGasLimit() (uint64, error) { return 0, nil }

func (tx *SSZTx) GetMessages() ([]proto.Message, error) { return nil, nil }

func (tx *SSZTx) GetSenders() ([][]byte, error) { return nil, nil }
