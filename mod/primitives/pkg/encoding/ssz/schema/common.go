package schema

import "github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"

const (
	B4Size  = 4
	B8Size  = 8
	B16Size = 16
	B20Size = 20
	B32Size = 32
	B48Size = 48
	B64Size = 64
	B96Size = 96
)

// Basic SSZ types.
//
//nolint:gochecknoglobals // reduce allocs.
var (
	boolType = basic(constants.BoolSize)
	u8Type   = basic(constants.U8Size)
	u16Type  = basic(constants.U16Size)
	u32Type  = basic(constants.U32Size)
	u64Type  = basic(constants.U64Size)
	u128Type = basic(constants.U128Size)
	u256Type = basic(constants.U256Size)
)

// Bool returns an SSZType representing a boolean.
func Bool() SSZType { return boolType }

// U8 returns an SSZType representing an 8-bit unsigned integer.
func U8() SSZType { return u8Type }

// U16 returns an SSZType representing a 16-bit unsigned integer.
func U16() SSZType { return u16Type }

// U32 returns an SSZType representing a 32-bit unsigned integer.
func U32() SSZType { return u32Type }

// U64 returns an SSZType representing a 64-bit unsigned integer.
func U64() SSZType { return u64Type }

// U128 returns an SSZType representing a 128-bit unsigned integer.
func U128() SSZType { return u128Type }

// U256 returns an SSZType representing a 256-bit unsigned integer.
func U256() SSZType { return u256Type }

// B4 creates a Vector of 4 bytes (32 bits).
func B4() SSZType { return Vector(U8(), B4Size) }

// B8 creates a Vector of 8 bytes (64 bits).
func B8() SSZType { return Vector(U8(), B8Size) }

// B16 creates a Vector of 16 bytes (128 bits).
func B16() SSZType { return Vector(U8(), B16Size) }

// B20 creates a Vector of 20 bytes (160 bits).
func B20() SSZType { return Vector(U8(), B20Size) }

// B32 creates a Vector of 32 bytes (256 bits).
func B32() SSZType { return Vector(U8(), B32Size) }

// B48 creates a Vector of 48 bytes (384 bits).
func B48() SSZType { return Vector(U8(), B48Size) }

// B64 creates a Vector of 64 bytes (512 bits).
func B64() SSZType { return Vector(U8(), B64Size) }

// B96 creates a Vector of 96 bytes (768 bits).
func B96() SSZType { return Vector(U8(), B96Size) }
