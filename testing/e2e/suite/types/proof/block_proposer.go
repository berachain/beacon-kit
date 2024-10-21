package proof

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// GetBlockProposerProof retrieves the block proposer index, block proposer pubkey,
// and validator pubkey proof for the specified timestamp.
func (c *Client) GetBlockProposerProof(
	ctx context.Context, timestamp uint64,
) (math.U64, bytes.B48, [][32]byte, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf(
		"%s/%s/t%d",
		c.baseURL, blockProposerEndpoint, timestamp,
	))
	if err != nil {
		return 0, bytes.B48{}, nil, err
	}
	var bppr blockProposerProofResponse
	if err = json.NewDecoder(resp.Body).Decode(&bppr); err != nil {
		return 0, bytes.B48{}, nil, err
	}

	proof := make([][32]byte, len(bppr.Proof))
	for i, root := range bppr.Proof {
		proof[i] = [32]byte(root)
	}
	return bppr.BeaconBlockHeader.ProposerIndex, bppr.Pubkey, proof, nil
}
