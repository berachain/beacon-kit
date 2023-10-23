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
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	prsymnetwork "github.com/prysmaticlabs/prysm/v4/network"
	"github.com/prysmaticlabs/prysm/v4/network/authorization"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/itsdevbear/bolaris/beacon/log"
	"github.com/itsdevbear/bolaris/beacon/prysm"
)

type (
	// TxPool represents the `TxPool` that exists on the backend of the execution layer.
	TxPoolAPI interface {
		Add([]*coretypes.Transaction, bool, bool) []error
		Stats() (int, int)
		SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	}

	EngineAPI interface {
		prysm.EngineCaller
		prysm.ExecutionBlockCaller
	}
)

// ExecutionClient represents the execution layer client.
type ExecutionClient struct {
	TxPoolAPI
	EngineAPI
}

// NewRemoteExecutionClient creates a new remote execution client.
func NewRemoteExecutionClient(
	dialURL, jwtSecretPath string,
	logger log.Logger) (*ExecutionClient, error) {
	ctx := context.Background()
	var (
		client  *rpc.Client
		chainID *big.Int
		err     error
	)

	jwtSecret, err := loadJWTSecret(jwtSecretPath)
	if err != nil {
		return nil, err
	}

	endpoint := NewPrysmEndpoint(dialURL, jwtSecret)
	client, err = newRPCClientWithAuth(ctx, nil, endpoint)
	if err != nil {
		return nil, err
	}

	var ethClient *ethclient.Client
	for i := 0; i < 100; func() { i++; time.Sleep(time.Second) }() {
		logger.Info("waiting for connection to execution layer", "dial-url", dialURL)
		ethClient = ethclient.NewClient(client)
		chainID, err = ethClient.ChainID(ctx)
		if err != nil {
			continue
		}
		logger.Info("Successfully connected to execution layer", "ChainID", chainID)
		break
	}
	if client == nil || err != nil {
		return nil, fmt.Errorf("failed to establish connection to execution layer: %w", err)
	}

	prsymClient := prysm.NewEngineClientService(ethClient)

	return &ExecutionClient{
		TxPoolAPI: &txPoolAPI{Client: ethClient},
		EngineAPI: prsymClient,
	}, nil
}

func loadJWTSecret(filepath string) ([]byte, error) {
	if data, err := os.ReadFile(filepath); err == nil {
		jwtSecret := common.FromHex(strings.TrimSpace(string(data)))
		if len(jwtSecret) == 32 { //nolint:gomnd // false positive.
			// log.Info("Loaded JWT secret file", "path", filepath, "crc32",
			// ("%#x", crc32.ChecksumIEEE(jwtSecret))
			return jwtSecret, nil
		}
	}
	// log.Error("Invalid JWT secret", "path", filepath, "length", len(jwtSecret))
	return nil, errors.New("invalid JWT secret")
}

// Initializes an RPC connection with authentication headers.
func newRPCClientWithAuth(ctx context.Context, headersMap map[string]string,
	endpoint prsymnetwork.Endpoint) (*rpc.Client, error) {
	headers := http.Header{}
	if endpoint.Auth.Method != authorization.None {
		header, err := endpoint.Auth.ToHeaderValue()
		if err != nil {
			return nil, err
		}
		headers.Set("Authorization", header)
	}
	for _, h := range headersMap {
		if h == "" {
			continue
		}
		keyValue := strings.Split(h, "=")
		if len(keyValue) < 2 { //nolint:gomnd // false positive.
			// log.LoggerWarn("Incorrect HTTP header flag format. Skipping %v", keyValue[0])
			continue
		}
		headers.Set(keyValue[0], strings.Join(keyValue[1:], "="))
	}

	return prsymnetwork.NewExecutionRPCClient(ctx, endpoint, headers)
}

func NewPrysmEndpoint(endpointString string, secret []byte) prsymnetwork.Endpoint {
	if len(secret) == 0 {
		return prsymnetwork.HttpEndpoint(endpointString)
	}
	// Overwrite authorization type for all endpoints to be of a bearer type.
	hEndpoint := prsymnetwork.HttpEndpoint(endpointString)
	hEndpoint.Auth.Method = authorization.Bearer
	hEndpoint.Auth.Value = string(secret)

	return hEndpoint
}
