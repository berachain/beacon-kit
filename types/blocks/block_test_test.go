package blocks

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

func TestRoundtripMarshalSSZ(t *testing.T) {
	// Create a new ExecutionPayloadHeaderCapella object
	block := FunBlock{
		Henlo: make([]byte, 32),
		ExecutionPayload: &enginev1.ExecutionPayloadCapella{
			ParentHash:    common.Hash{1}.Bytes(),
			FeeRecipient:  common.Address{1}.Bytes(),
			StateRoot:     common.Hash{1}.Bytes(),
			ReceiptsRoot:  common.Hash{1}.Bytes(),
			LogsBloom:     make([]byte, 256),
			PrevRandao:    common.Hash{1}.Bytes(),
			BlockNumber:   0,
			GasLimit:      0,
			GasUsed:       0,
			Timestamp:     0,
			ExtraData:     common.Hash{1}.Bytes(),
			BaseFeePerGas: common.Hash{1}.Bytes(),
			BlockHash:     common.Hash{1}.Bytes(),
			Transactions:  [][]byte{},
			Withdrawals: []*enginev1.Withdrawal{
				{
					Amount:  65,
					Address: common.Address{1}.Bytes(),
				},
			},
		},
	}

	// Marshal the object to SSZ
	marshaled, err := block.MarshalSSZ()
	fmt.Println(len(marshaled))
	fmt.Println(uint64(len(marshaled)))
	require.NoError(t, err)

	// Unmarshal the SSZ back to the object
	unmarshaled := new(FunBlock)
	fmt.Println("NEW", unmarshaled)

	err = unmarshaled.UnmarshalSSZ(marshaled)
	require.NoError(t, err)

	// Assert that the original and unmarshaled objects are equal
	require.Equal(t, block.Henlo, unmarshaled.Henlo)
	require.Equal(t, block, *unmarshaled)
}
