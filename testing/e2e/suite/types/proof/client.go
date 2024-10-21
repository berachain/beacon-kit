package proof

import "net/http"

// HTTPClient is a http client that can make GET requests.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// Client is a client for the proof namespace of the Beacon node API.
type Client struct {
	httpClient HTTPClient
	baseURL    string
}

// NewClient creates a new client for the proof namespace of the Beacon node
// API.
func NewClient(httpClient HTTPClient, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}
