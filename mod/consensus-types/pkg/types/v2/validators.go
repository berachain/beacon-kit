package v2

import (
	"github.com/karalabe/ssz"
)

// Validators represents a list of Validator objects.
type Validators []*Validator

// SizeSSZ returns the SSZ encoded size of the Validators.
func (v Validators) SizeSSZ() uint32 {
	return uint32(len(v)) * v.SizeSSZ()
}

// DefineSSZ defines the SSZ encoding for the Validators.
func (v Validators) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineSliceOfStaticObjectsContent(codec, (*[]*Validator)(&v), 1099511627776)
}
