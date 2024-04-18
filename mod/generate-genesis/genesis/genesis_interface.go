package genesis

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Genesis interface {
	AddAccount(address common.Address, balance *big.Int)
	AddPredeploy(address common.Address, code []byte, balance *big.Int, nonce uint64)
	ToJSON(filename string) error
}
