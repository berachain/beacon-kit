package types

import "reflect"

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
)

var (
	hashT = reflect.TypeOf(Hash{})
	// addressT = reflect.TypeOf(Address{})

	// MaxAddress represents the maximum possible address value.
	// MaxAddress = HexToAddress("0xffffffffffffffffffffffffffffffffffffffff")

	// MaxHash represents the maximum possible hash value.
	MaxHash = HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
)
