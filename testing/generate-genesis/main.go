package main

import (
	"log"

	"github.com/berachain/beacon-kit/testing/generate-genesis/cmd"
)

func main() {
	if err := cmd.CreateEthGenesisCmd().Execute(); err != nil {
		log.Println("Error:", err)
	}
}
