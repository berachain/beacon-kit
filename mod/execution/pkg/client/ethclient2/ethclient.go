package ethclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	json "github.com/goccy/go-json"
)

// EthRPC - Ethereum rpc client
type EthRPC struct {
	url    string
	client *http.Client
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *EthRPC)) *EthRPC {
	rpc := &EthRPC{
		url:    url,
		client: http.DefaultClient,
		// log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// Call calls the given method with the given parameters
func (rpc *EthRPC) Call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.CallRaw(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Call returns raw response of method call
func (rpc *EthRPC) CallRaw(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
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

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// Web3ClientVersion returns the current client version.
func (rpc *EthRPC) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.Call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Web3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *EthRPC) Web3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.Call("web3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// NetVersion returns the current network protocol version.
func (rpc *EthRPC) NetVersion() (string, error) {
	var version string

	err := rpc.Call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (rpc *EthRPC) NetListening() (bool, error) {
	var listening bool

	err := rpc.Call("net_listening", &listening)
	return listening, err
}
