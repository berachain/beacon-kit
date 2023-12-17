// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2023, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	prsymnetwork "github.com/prysmaticlabs/prysm/v4/network"
	"github.com/prysmaticlabs/prysm/v4/network/authorization"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/itsdevbear/bolaris/beacon/log"
)

const (
	defaultEndpointRetries    = 50
	defaultEndpointRetryDelay = 1 * time.Second
)

// NewAuthenticatedEthClient creates a new remote execution client.
func NewAuthenticatedEthClient(
	dialURL, jwtSecretPath string, logger log.Logger,
) (*ethclient.Client, error) {
	ctx := context.Background()
	var (
		jwtSecret []byte
		client    *rpc.Client
		ethClient *ethclient.Client
		chainID   *big.Int
		err       error
	)

	// Load the JWT secret from the provided path.
	if jwtSecret, err = loadJWTSecret(jwtSecretPath); err != nil {
		return nil, err
	}

	// Create a new Prysm endpoint with the dialURL and the loaded JWT secret.
	endpoint := newPrysmEndpoint(dialURL, jwtSecret)

	// Attempt to establish a new RPC client with authentication.

	if client, err = newRPCClientWithAuth(ctx, nil, endpoint); err != nil {
		return nil, err
	}

	// Attempt to connect to the execution layer and retrieve the chain ID.
	// Retry up to 100 times, with a 1-second delay between each attempt.
	for i := 0; i < defaultEndpointRetries; func() { i++; time.Sleep(defaultEndpointRetryDelay) }() {
		logger.Info("waiting for connection to execution layer", "dial-url", dialURL)
		ethClient = ethclient.NewClient(client)
		chainID, err = ethClient.ChainID(ctx)
		if err != nil {
			continue
		}
		// Log the successful connection and the chain ID.
		logger.Info("Successfully connected to execution layer", "ChainID", chainID)
		break
	}

	// If the connection still fails after 100 attempts, return an error.
	if client == nil || err != nil {
		return nil, fmt.Errorf("failed to establish connection to execution layer: %w", err)
	}

	return ethClient, nil
}

// loadJWTSecret reads the JWT secret from a file and returns it.
// It returns an error if the file cannot be read or if the JWT secret is not valid.
func loadJWTSecret(filepath string) ([]byte, error) {
	// Read the file.
	data, err := os.ReadFile(filepath)
	if err != nil {
		// Return an error if the file cannot be read.
		return nil, err
	}

	// Convert the data to a JWT secret.
	jwtSecret := common.FromHex(strings.TrimSpace(string(data)))
	// Check if the JWT secret is valid.
	if len(jwtSecret) == 32 { //nolint:gomnd // false positive.
		// Log that the JWT secret file has been loaded.
		// TODO: remove println.
		fmt.Println("Loaded JWT secret file", "path", filepath, "crc32")
		// ("%#x", crc32.ChecksumIEEE(jwtSecret))
		// Return the JWT secret.
		return jwtSecret, nil
	}
	// Return an error if the JWT secret is not valid.
	return nil, fmt.Errorf("failed to load JWT secret from %s", filepath)
}

// newRPCClientWithAuth initializes an RPC connection with authentication headers.
// It takes a context, a map of headers, and an endpoint as arguments.
// It returns an RPC client and an error.
func newRPCClientWithAuth(ctx context.Context, headersMap map[string]string,
	endpoint prsymnetwork.Endpoint) (*rpc.Client, error) {
	// Initialize an empty HTTP headers object.
	headers := http.Header{}
	// If the endpoint has an authorization method, add it to the headers.
	if endpoint.Auth.Method != authorization.None {
		// Convert the authorization data to a header value.
		header, err := endpoint.Auth.ToHeaderValue()
		// If there is an error, return it.
		if err != nil {
			return nil, err
		}
		// Add the authorization header to the headers.
		headers.Set("Authorization", header)
	}
	// Iterate over the headers map.
	for _, h := range headersMap {
		// If the header is empty, skip it.
		if h == "" {
			continue
		}
		// Split the header into a key and a value.
		keyValue := strings.Split(h, "=")
		// If the header does not have a key and a value, skip it.
		if len(keyValue) < 2 { //nolint:gomnd // false positive.
			// log.LoggerWarn("Incorrect HTTP header flag format. Skipping %v", keyValue[0])
			continue
		}
		// Add the header to the headers.
		headers.Set(keyValue[0], strings.Join(keyValue[1:], "="))
	}

	// Return a new RPC client with the endpoint and the headers.
	return prsymnetwork.NewExecutionRPCClient(ctx, endpoint, headers)
}

// newPrysmEndpoint creates a new Prysm network endpoint.
// If a secret is provided, it sets the authorization type to bearer and uses the secret as the value.
// If no secret is provided, it simply returns the HTTP endpoint.
func newPrysmEndpoint(endpointString string, secret []byte) prsymnetwork.Endpoint {
	// If no secret is provided, return the HTTP endpoint.
	if len(secret) == 0 {
		return prsymnetwork.HttpEndpoint(endpointString)
	}

	// If a secret is provided, overwrite the authorization type for all endpoints to be of a bearer type.
	hEndpoint := prsymnetwork.HttpEndpoint(endpointString)
	hEndpoint.Auth.Method = authorization.Bearer
	hEndpoint.Auth.Value = string(secret)

	// Return the modified endpoint with bearer authorization.
	return hEndpoint
}
