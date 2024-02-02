package proposal

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/assert"
)

func TestRemoveBeaconBlockFromTxs(t *testing.T) {
	tests := []struct {
		name            string
		inputTxs        [][]byte
		payloadPosition uint
		expectedTxs     [][]byte
	}{
		{
			name:            "remove first element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 0,
			expectedTxs:     [][]byte{{2}, {3}},
		},
		{
			name:            "remove last element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 2,
			expectedTxs:     [][]byte{{1}, {2}},
		},
		{
			name:            "remove middle element",
			inputTxs:        [][]byte{{1}, {2}, {3}},
			payloadPosition: 1,
			expectedTxs:     [][]byte{{1}, {3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				payloadPosition: tt.payloadPosition,
			}
			req := &abci.RequestProcessProposal{Txs: tt.inputTxs}
			result := handler.removeBeaconBlockFromTxs(req)
			assert.Equal(t, tt.expectedTxs, result.Txs)
		})
	}
}
