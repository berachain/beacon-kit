package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/karalabe/ssz"
)

// Deposits is a typealias for a list of Deposits.
type Deposits []*Deposit

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the SSZ encoded size in bytes for the Deposits.
func (ds Deposits) SizeSSZ(bool) uint32 {
	return ssz.SizeSliceOfStaticObjects(([]*Deposit)(ds))
}

// DefineSSZ defines the SSZ encoding for the Deposits object.
func (ds Deposits) DefineSSZ(c *ssz.Codec) {
	c.DefineDecoder(func(*ssz.Decoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
	c.DefineEncoder(func(*ssz.Encoder) {
		ssz.DefineSliceOfStaticObjectsContent(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
	c.DefineHasher(func(*ssz.Hasher) {
		ssz.DefineSliceOfStaticObjectsOffset(
			c, (*[]*Deposit)(&ds), constants.MaxDepositsPerBlock)
	})
}

// HashTreeRoot returns the hash tree root of the Deposits.
func (ds Deposits) HashTreeRoot() common.Root {
	return ssz.HashSequential(ds)
}
