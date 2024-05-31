package genesis

import (
	"math/big"
)

// Genesis is an interface that defines the methods that a genesis configuration struct should implement.
type Genesis interface {
	AddAccount(address string, balance *big.Int) error
	AddPredeploy(address string, code []byte, balance *big.Int, nonce uint64) error
	WriteJSON(filename string) error
}
