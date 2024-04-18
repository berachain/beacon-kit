package genesis

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Genesis is an interface that defines the methods that a genesis configuration struct should implement.
type Genesis interface {
	AddAccount(address common.Address, balance *big.Int)
	AddPredeploy(address common.Address, code []byte, balance *big.Int, nonce uint64)
	ToJSON(filename string) error
}
