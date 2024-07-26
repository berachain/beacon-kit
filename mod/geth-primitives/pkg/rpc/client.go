package rpc

import (
	"context"

	jsonrpc "github.com/ybbus/jsonrpc/v3"
)

type Client2 struct {
	jsonrpc.RPCClient
}

func NewClient2(url string) *Client2 {
	return &Client2{
		jsonrpc.NewClient(url),
	}
}

func (c *Client2) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	res, err := c.RPCClient.Call(ctx, method, args)
	if err != nil {
		return err
	}
	return res.GetObject(result)

}

func (c *Client2) Call(method string, result interface{}, args ...interface{}) error {
	ctx := context.Background()
	return c.CallContext(ctx, result, method, args...)
}
