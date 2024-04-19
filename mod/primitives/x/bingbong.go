package main

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/x/bet"
)

func main() {

	// U64T U64[U64T], RootT ~[32]byte,
	// SpecT any, C Container[RootT]](value C) ([]byte, error) {
	x := bet.ItemA{
		Val1: math.U64(1),
		Val2: []math.U64{math.U64(2)},
	}
	y := bet.ItemB{
		Val1: math.U64(1),
		Val2: []uint64{2},
	}
	bz, err := ssz.SerializeContainer[math.U64, primitives.Root, any](x)
	if err != nil {
		panic(err)
	}

	bz2, err := y.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	fmt.Println(bz, bz2)
}
