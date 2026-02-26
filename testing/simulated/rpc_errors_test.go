//go:build simulated

// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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

package simulated_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/execution/client/ethclient"
	"github.com/berachain/beacon-kit/log/phuslu"
	jsonrpc "github.com/berachain/beacon-kit/primitives/net/json-rpc"
	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// HTTP proxy in between beacon node and execution client. When active,
// replaces responses with specified JSON-RPC error.
type rpcErrorProxy struct {
	targetURL  string
	active     atomic.Bool
	errorCode  int
	errorMsg   string
	httpClient *http.Client
}

func newRPCErrorProxy(targetURL string) *rpcErrorProxy {
	return &rpcErrorProxy{
		targetURL:  targetURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *rpcErrorProxy) activate(code int, msg string) {
	p.errorCode = code
	p.errorMsg = msg
	p.active.Store(true)
}

func (p *rpcErrorProxy) deactivate() {
	p.active.Store(false)
}

func (p *rpcErrorProxy) getErr(reqId json.RawMessage) string {
	return fmt.Sprintf(
		`{"jsonrpc":"2.0","id":%s,"error":{"code":%d,"message":"%s"}}`,
		string(reqId), p.errorCode, p.errorMsg,
	)
}

func (p *rpcErrorProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "proxy read error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Check if need interceptor.
	if p.active.Load() {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if json.Unmarshal(bodyBytes, &req) == nil {
			// Intercept targeted methods.
			if isTargetedEngineMethod(req.Method) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(p.getErr(req.ID)))
				return
			}
		}
	}

	// Forward original request.
	proxyReq, err := http.NewRequestWithContext(
		r.Context(), r.Method, p.targetURL,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		http.Error(w, "proxy forward error", http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header.Clone()

	resp, err := p.httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "proxy upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func isTargetedEngineMethod(method string) bool {
	switch method {
	case ethclient.NewPayloadMethodV3,
		ethclient.NewPayloadMethodV4,
		ethclient.NewPayloadMethodV4P11,
		ethclient.ForkchoiceUpdatedMethodV3,
		ethclient.ForkchoiceUpdatedMethodV3P11:
		return true
	}
	return false
}

type RPCErrorProxySuite struct {
	suite.Suite
	simulated.SharedAccessors
	errProxy       *rpcErrorProxy
	errProxyServer *httptest.Server
}

func TestRPCErrorProxySuite(t *testing.T) {
	suite.Run(t, new(RPCErrorProxySuite))
}

// SetupTest inserts a proxy in between the node and execution client,
// to enable injection and testing of JSON-RPC errors.
func (s *RPCErrorProxySuite) SetupTest() {
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())
	s.CtxComet = context.TODO()
	s.HomeDir = s.T().TempDir()

	const elGenesisPath = "./el-genesis-files/eth-genesis.json"
	chainSpecFunc := simulated.ProvideSimulationChainSpec
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), chainSpec, s.HomeDir, elGenesisPath)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start Geth.
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImage())
	elHandle, authRPC, elRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
	s.ElHandle = elHandle

	// Create the error proxy for AuthRPC.
	s.errProxy = newRPCErrorProxy(authRPC.String())
	s.errProxyServer = httptest.NewServer(s.errProxy)

	// Create a ConnectionURL pointing to the proxy instead of Geth.
	proxyURL, err := url.NewFromRaw(s.errProxyServer.URL)
	s.Require().NoError(err)

	s.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	// Use proxy connection URL as AuthRPC
	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     proxyURL,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	s.SimComet = s.TestNode.SimComet

	go func() {
		_ = s.TestNode.Start(s.CtxApp)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.EngineClient)
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

func (s *RPCErrorProxySuite) TearDownTest() {
	s.errProxyServer.Close()
	s.CleanupTest(s.T())
}

// preparedProposal holds the state needed to call FinalizeBlock.
type preparedProposal struct {
	txs             [][]byte
	height          int64
	proposerAddress []byte
	proposalTime    time.Time
}

// prepareForFinalize advances the chain and prepares a proposal, returning
// the data needed to call FinalizeBlock.
func (s *RPCErrorProxySuite) prepareForFinalize() preparedProposal {
	s.T().Helper()

	const blockHeight = 1
	const coreLoopIterations = 1

	s.InitializeChain(s.T())
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	startTime := time.Now()

	proposals, _, proposalTime := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, nodeAddress, startTime)
	s.Require().Len(proposals, coreLoopIterations)

	currentHeight := int64(blockHeight + coreLoopIterations)

	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          currentHeight,
		Time:            proposalTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	s.LogBuffer.Reset()

	return preparedProposal{
		txs:             proposal.Txs,
		height:          currentHeight,
		proposerAddress: nodeAddress,
		proposalTime:    proposalTime,
	}
}

// TestFinalizeBlock_FatalRPCError shows that when exec client returns a
// JSON-RPC error (e.g. -32700 parse error) during FinalizeBlock, the error is
// correctly identified and returned.
func (s *RPCErrorProxySuite) TestFinalizeBlock_HandleRPCError() {
	pp := s.prepareForFinalize()

	// Activate the error proxy with an RPC error code (-32700 parse error).
	s.errProxy.activate(-32700, "Parse Error")

	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	})

	s.Require().Error(err, "FinalizeBlock should fail on fatal RPC error")
	s.Require().Nil(finalizeResp)
	s.Require().ErrorIs(err, jsonrpc.ErrParse, "Error should be correctly classified")
}
