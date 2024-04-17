package main

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/node-builder/config/spec"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/yap/bet"
)

func main() {
	chainSpec := spec.LocalnetChainSpec()
	deposit := &primitives.Deposit{}

	yapCave := bet.YapCave[primitives.ChainSpec]{
		Item2: []*primitives.Deposit{deposit},
		// Item3: 117,
	}

	yapCave2 := bet.YapCave2{
		Item2: []*primitives.Deposit{deposit},
		// Item3: 117,
	}

	root := yapCave.HashTreeRoot(chainSpec)
	fmt.Println(root)

	root2, err := (&yapCave2).HashTreeRoot()
	if err != nil {
		fmt.Println(err)
		return
	}

	// t, _ := yapCave2.GetTree()
	// t.Show(3)
	fmt.Println(primitives.Root(root2))

}
