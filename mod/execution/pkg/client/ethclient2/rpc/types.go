package rpc

import (
	"fmt"

	json "github.com/goccy/go-json"
)

// Request represents an Ethereum JSON-RPC request.
type Request struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// Response represents an Ethereum JSON-RPC response.
type Response struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error"`
}

// Error represents an Ethereum JSON-RPC error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns a formatted error string.
func (err Error) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}
