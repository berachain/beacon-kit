package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	json "github.com/goccy/go-json"
)

// EthRPC - Ethereum rpc client
type Client struct {
	url       string
	client    *http.Client
	reqPool   *sync.Pool
	JwtSecret *jwt.Secret
	header    http.Header
}

// New create new rpc client with given url
func NewClient(url string, options ...func(rpc *Client)) *Client {
	rpc := &Client{
		url:    url,
		client: http.DefaultClient,
		reqPool: &sync.Pool{
			New: func() interface{} {
				return &Request{
					ID:      1,
					JSONRPC: "2.0",
				}
			},
		},
		header: make(http.Header),
	}

	rpc.header.Set("Content-Type", "application/json")

	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// Start starts the rpc client
func (rpc *Client) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	rpc.updateHeader()
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			rpc.updateHeader()
		}
	}
}

// Close closes the RPC client
func (rpc *Client) Close() error {
	rpc.client.CloseIdleConnections()
	return nil
}

// updateHeader builds an http.Header that has the JWT token
// attached for authorization.
func (rpc *Client) updateHeader() error {
	// Build the JWT token.
	token, err := rpc.JwtSecret.BuildSignedToken()
	if err != nil {
		return err
	}

	// Add the JWT token to the headers.
	rpc.header.Set("Authorization", "Bearer "+token)
	return nil
}

// Call calls the given method with the given parameters
func (rpc *Client) Call(
	ctx context.Context, target interface{}, method string, params ...interface{},
) error {
	result, err := rpc.CallRaw(ctx, method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Call returns raw response of method call
func (rpc *Client) CallRaw(
	ctx context.Context, method string, params ...interface{},
) (json.RawMessage, error) {
	// Pull a request from the pool, we know that it already has the correct
	// JSONRPC version and ID set.
	request := rpc.reqPool.Get().(*Request)
	defer rpc.reqPool.Put(request)

	// Update the request with the method and params.
	request.Method = method
	request.Params = params

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rpc.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = rpc.header

	response, err := rpc.client.Do(req)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}
