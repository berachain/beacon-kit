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
	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// HTTP proxy in between beacon node and execution client. Supports three
// injection modes:
//   - activate(code, msg): return HTTP 200 with the given JSON-RPC error body
//   - activateHTTPStatus(status): return the given HTTP status with a generic body
//     (exercises the transport-level classification path, e.g. HTTP 4xx fatal)
//   - activateDropConn: hijack and close the TCP connection (unreachable EL)
type rpcErrorProxy struct {
	targetURL    string
	active       atomic.Bool
	dropConn     atomic.Bool
	httpStatus   atomic.Int32 // 0 = inactive; otherwise the status code to return
	errorCode    int
	errorMsg     string
	httpClient   *http.Client
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

// activateHTTPStatus makes the proxy respond to engine-API requests with the
// given HTTP status code and a short body. Use to exercise transport-level
// classification (e.g. 4xx fatal vs 5xx retryable).
func (p *rpcErrorProxy) activateHTTPStatus(statusCode int) {
	p.httpStatus.Store(int32(statusCode))
}

// activateDropConn simulates an unreachable EL by dropping the TCP
// connection on any engine-API request.
func (p *rpcErrorProxy) activateDropConn() {
	p.dropConn.Store(true)
}

func (p *rpcErrorProxy) deactivate() {
	p.active.Store(false)
	p.dropConn.Store(false)
	p.httpStatus.Store(0)
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

	if p.intercept(w, bodyBytes) {
		return
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

// intercept reports whether the request should be intercepted and writes
// the intercepted response to w. Returns true when the request was handled.
func (p *rpcErrorProxy) intercept(w http.ResponseWriter, bodyBytes []byte) bool {
	httpStatus := int(p.httpStatus.Load())
	if !p.active.Load() && !p.dropConn.Load() && httpStatus == 0 {
		return false
	}
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
	}
	if json.Unmarshal(bodyBytes, &req) != nil || !isTargetedEngineMethod(req.Method) {
		return false
	}
	if p.dropConn.Load() {
		dropTCPConn(w)
		return true
	}
	if httpStatus != 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(httpStatus)
		_, _ = fmt.Fprintf(w, "injected HTTP %d", httpStatus)
		return true
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(p.getErr(req.ID)))
	return true
}

// dropTCPConn hijacks and closes the TCP connection to simulate an
// unreachable EL.
func dropTCPConn(w http.ResponseWriter) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack not supported", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, "hijack failed", http.StatusInternalServerError)
		return
	}
	_ = conn.Close()
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
	configs, genesisValidatorsRoot := simulated.InitializeHomeDirs(s.T(), chainSpec, elGenesisPath, s.HomeDir)
	cometConfig := configs[0]
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

	s.InitializeChain(s.T(), 1)
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

// TestFinalizeBlock_FatalRPCError_Surfaces shows that a fatal engine-API
// response (e.g. -32700 parse error, HTTP 4xx) during FinalizeBlock is
// surfaced immediately rather than retried. Fatal errors encode "this request
// will never succeed against this EL" — retrying forever would turn a
// misconfigured JWT or wrong chain ID into a silent node hang. Transient
// outages take the IsNonFatalError path instead, which still retries
// indefinitely under PhaseFinalize (see TestFinalizeBlock_ConnectionDrop_Recovery).
func (s *RPCErrorProxySuite) TestFinalizeBlock_FatalRPCError_Surfaces() {
	pp := s.prepareForFinalize()

	// Inject -32700 parse errors on every engine-API call and leave them on.
	s.errProxy.activate(-32700, "Parse Error")
	defer s.errProxy.deactivate()

	_, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	})

	s.Require().Error(err, "FinalizeBlock must surface fatal engine-API errors instead of looping")

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "fatal error", "Should log the fatal error")
}

// TestFinalizeBlock_HTTP4xx_Surfaces is the integration-level twin of
// TestFinalizeBlock_FatalRPCError_Surfaces, covering the load-bearing
// transport-level path that PR #3109 was about: an HTTP 4xx response from
// the EL (e.g. 413 oversized payload, 401 bad JWT) must surface immediately
// rather than trap FinalizeBlock in an infinite retry. This protects against
// a misconfigured EL hanging a validator node forever.
func (s *RPCErrorProxySuite) TestFinalizeBlock_HTTP4xx_Surfaces() {
	pp := s.prepareForFinalize()

	// Inject HTTP 413 (the original PoC) on every engine-API call.
	s.errProxy.activateHTTPStatus(http.StatusRequestEntityTooLarge)
	defer s.errProxy.deactivate()

	_, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	})

	s.Require().Error(err, "FinalizeBlock must surface HTTP 4xx instead of looping")

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "fatal error", "Should log the fatal error")
}

// TestFinalizeBlock_HTTP429_Recovers pins the carve-out for RFC-retryable
// 4xx codes: a proxy/rate-limiter returning 429 must not surface as fatal,
// because retrying after backoff is the prescribed response. FinalizeBlock
// recovers once the 429 stops.
func (s *RPCErrorProxySuite) TestFinalizeBlock_HTTP429_Recovers() {
	pp := s.prepareForFinalize()

	s.errProxy.activateHTTPStatus(http.StatusTooManyRequests)

	// Stop returning 429 after a short delay so the retry can succeed.
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.errProxy.deactivate()
	}()

	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	})

	s.Require().NoError(err, "FinalizeBlock should recover after rate limit lifts")
	s.Require().NotNil(finalizeResp)

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "non fatal error", "429 must be classified as retryable, not fatal")
}

// TestFinalizeBlock_ConnectionDrop_Recovery shows that when the EL is
// unreachable (e.g. bera-reth restart) the engine keeps retrying and
// FinalizeBlock succeeds once the EL comes back.
func (s *RPCErrorProxySuite) TestFinalizeBlock_ConnectionDrop_Recovery() {
	pp := s.prepareForFinalize()

	// Simulate the EL going away: next engine-API requests have their TCP
	// connection dropped.
	s.errProxy.activateDropConn()

	// Bring the EL back after a short delay so the retry can succeed.
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.errProxy.deactivate()
	}()

	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	})

	s.Require().NoError(err, "FinalizeBlock should recover after EL comes back")
	s.Require().NotNil(finalizeResp)

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "non fatal error", "Should log non fatal retry attempts")
}
