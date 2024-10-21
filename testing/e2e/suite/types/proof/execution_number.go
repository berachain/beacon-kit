package proof

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetExecutionNumberProof gets the execution number proof for the given
// timestamp.
func (c *Client) GetExecutionNumberProof(
	ctx context.Context, timestamp uint64,
) ([][32]byte, error) {
	resp, err := c.httpClient.Get(
		fmt.Sprintf(
			"%s/%s/t%d", c.baseURL, executionNumberEndpoint, timestamp,
		),
	)
	if err != nil {
		return nil, err
	}

	var enr executionNumberResponse
	if err = json.NewDecoder(resp.Body).Decode(&enr); err != nil {
		return nil, err
	}

	proof := make([][32]byte, len(enr.Proof))
	for i, root := range enr.Proof {
		proof[i] = [32]byte(root)
	}
	return proof, nil
}
