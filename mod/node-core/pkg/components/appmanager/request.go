package appmanager

import (
	"time"

	appmanager "cosmossdk.io/core/app"
	"cosmossdk.io/core/transaction"
)

type abciRequest struct {
	height int64
	time   time.Time
	txs    [][]byte
}

func blockToABCIRequest[T transaction.Tx](
	block *appmanager.BlockRequest[T],
) *abciRequest {
	txs := make([][]byte, len(block.Txs))
	for i, tx := range block.Txs {
		txs[i] = tx.Bytes()
	}
	return &abciRequest{
		height: int64(block.Height),
		time:   block.Time,
		txs:    txs,
	}
}

func (req *abciRequest) GetHeight() int64 {
	return req.height
}

func (req *abciRequest) GetTime() time.Time {
	return req.time
}

func (req *abciRequest) GetTxs() [][]byte {
	return req.txs
}

func (req *abciRequest) SetTxs(txs [][]byte) {
	req.txs = txs
}
