package node

import "github.com/berachain/beacon-kit/node-api/handlers"

// Version is a placeholder so that beacon API clients don't break.
// TODO: Implement with real data.
func (h *Handler) Version(handlers.Context) (any, error) {
	type VersionResponse struct {
		Data struct {
			Version string `json:"version"`
		} `json:"data"`
	}

	response := VersionResponse{}
	response.Data.Version = "1.1.0"

	return response, nil
}
