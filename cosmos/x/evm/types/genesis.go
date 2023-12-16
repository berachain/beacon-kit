package types

import "github.com/ethereum/go-ethereum/common"

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Eth1GenesisHash: common.Hex2Bytes("0x1de2227088b69dbf4f8a77c5c6721a4c221aa030212418a01702e23ba4934588"),
	}
}
