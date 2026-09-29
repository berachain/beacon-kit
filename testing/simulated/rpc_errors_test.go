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

// maxBoundedCall is the longest a bounded engine call may take with default
// config. It is the retry budget (1.5s) plus one RPC timeout (2s).
const maxBoundedCall = 3500 * time.Millisecond

// HTTP proxy in between beacon node and execution client. Supports five
// injection modes:
//   - activate(code, msg): return HTTP 200 with the given JSON-RPC error body
//   - activateHTTPStatus(status): return the given HTTP status with a generic body
//     (exercises the transport-level classification path, e.g. HTTP 4xx fatal)
//   - activateDropConn: hijack and close the TCP connection (unreachable EL)
//   - activateSyncingNewPayload: forward newPayload to the EL, answer SYNCING
//   - activateHang: never answer, the caller gives up at its RPC timeout
type rpcErrorProxy struct {
	targetURL         string
	active            atomic.Bool
	dropConn          atomic.Bool
	httpStatus        atomic.Int32 // 0 = inactive; otherwise the status code to return
	syncingNewPayload atomic.Bool
	hang              atomic.Bool
	fcuCalls          atomic.Int32 // forkchoiceUpdated requests seen by the proxy
	newPayloadCalls   atomic.Int32 // newPayload requests seen by the proxy
	errorCode         int
	errorMsg          string
	httpClient        *http.Client
}

// rpcRequest holds the JSON-RPC request fields the proxy cares about.
type rpcRequest struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
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

// activateSyncingNewPayload makes the proxy answer newPayload with a SYNCING
// status. The request is still forwarded, so the EL imports the payload.
func (p *rpcErrorProxy) activateSyncingNewPayload() {
	p.syncingNewPayload.Store(true)
}

// activateHang makes the proxy hold engine-API requests without answering,
// which simulates an EL that accepts the connection and then stalls.
func (p *rpcErrorProxy) activateHang() {
	p.hang.Store(true)
}

func (p *rpcErrorProxy) deactivate() {
	p.active.Store(false)
	p.dropConn.Store(false)
	p.httpStatus.Store(0)
	p.syncingNewPayload.Store(false)
	p.hang.Store(false)
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

	// Requests that fail to parse have no method and are forwarded untouched.
	var req rpcRequest
	_ = json.Unmarshal(bodyBytes, &req)
	if isForkchoiceUpdatedMethod(req.Method) {
		p.fcuCalls.Add(1)
	}
	if isNewPayloadMethod(req.Method) {
		p.newPayloadCalls.Add(1)
	}

	if p.hang.Load() && isTargetedEngineMethod(req.Method) {
		// The request context is done once the caller drops the connection.
		<-r.Context().Done()
		return
	}

	if p.intercept(w, req) {
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

	if p.syncingNewPayload.Load() && isNewPayloadMethod(req.Method) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w,
			`{"jsonrpc":"2.0","id":%s,"result":{"status":"SYNCING","latestValidHash":null,"validationError":null}}`,
			string(req.ID),
		)
		return
	}

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
func (p *rpcErrorProxy) intercept(w http.ResponseWriter, req rpcRequest) bool {
	// Snapshot all flags once so a concurrent deactivate() can't change them mid-request.
	active, dropConn, httpStatus := p.active.Load(), p.dropConn.Load(), int(p.httpStatus.Load())
	if !active && !dropConn && httpStatus == 0 {
		return false
	}
	if !isTargetedEngineMethod(req.Method) {
		return false
	}
	if dropConn {
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
	return isNewPayloadMethod(method) || isForkchoiceUpdatedMethod(method)
}

func isNewPayloadMethod(method string) bool {
	switch method {
	case ethclient.NewPayloadMethodV3,
		ethclient.NewPayloadMethodV4,
		ethclient.NewPayloadMethodV4P11:
		return true
	}
	return false
}

func isForkchoiceUpdatedMethod(method string) bool {
	switch method {
	case ethclient.ForkchoiceUpdatedMethodV3,
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

	// Start Reth.
	elNode := execution.NewRethNode(s.HomeDir, execution.ValidRethImage())
	elHandle, authRPC, elRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
	s.ElHandle = elHandle

	// Create the error proxy for AuthRPC.
	s.errProxy = newRPCErrorProxy(authRPC.String())
	s.errProxyServer = httptest.NewServer(s.errProxy)

	// Create a ConnectionURL pointing to the proxy instead of execution client.
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

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.ContractBackend)
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

func (s *RPCErrorProxySuite) TearDownTest() {
	s.errProxyServer.Close()
	s.CleanupTest(s.T())
}

// preparedProposal holds the state needed to call ProcessProposal and
// FinalizeBlock.
type preparedProposal struct {
	txs             [][]byte
	height          int64
	proposerAddress []byte
	proposalTime    time.Time
}

func (pp preparedProposal) processRequest() *types.ProcessProposalRequest {
	return &types.ProcessProposalRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	}
}

func (pp preparedProposal) finalizeRequest() *types.FinalizeBlockRequest {
	return &types.FinalizeBlockRequest{
		Txs:             pp.txs,
		Height:          pp.height,
		ProposerAddress: pp.proposerAddress,
		Time:            pp.proposalTime,
	}
}

// initChain initializes the chain and returns the address of the node.
func (s *RPCErrorProxySuite) initChain() []byte {
	s.T().Helper()

	s.InitializeChain(s.T(), 1)
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)
	return nodeAddress
}

// prepareProposal initializes the chain, finalizes every block below height
// and prepares a proposal for height.
func (s *RPCErrorProxySuite) prepareProposal(height int64) preparedProposal {
	s.T().Helper()

	const firstHeight = 1
	nodeAddress := s.initChain()

	finalized := height - firstHeight
	proposals, _, proposalTime := s.MoveChainToHeight(s.T(), firstHeight, finalized, nodeAddress, time.Now())
	s.Require().Len(proposals, int(finalized))

	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
		Height:          height,
		Time:            proposalTime,
		ProposerAddress: nodeAddress,
	})
	s.Require().NoError(err)
	s.Require().Len(proposal.Txs, 2)

	s.LogBuffer.Reset()

	return preparedProposal{
		txs:             proposal.Txs,
		height:          height,
		proposerAddress: nodeAddress,
		proposalTime:    proposalTime,
	}
}

// prepareForFinalize prepares a proposal on top of one finalized block, so
// the startup sync has already run.
func (s *RPCErrorProxySuite) prepareForFinalize() preparedProposal {
	s.T().Helper()

	const secondHeight = 2
	return s.prepareProposal(secondHeight)
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

	_, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, pp.finalizeRequest())

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

	_, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, pp.finalizeRequest())

	s.Require().Error(err, "FinalizeBlock must surface HTTP 4xx instead of looping")

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "fatal error", "Should log the fatal error")
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

	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, pp.finalizeRequest())

	s.Require().NoError(err, "FinalizeBlock should recover after EL comes back")
	s.Require().NotNil(finalizeResp)

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "non fatal error", "Should log non fatal retry attempts")
}

// TestFinalizeBlock_NewPayloadSyncing_SendsFCU shows that FinalizeBlock accepts
// a SYNCING status from NewPayload and sends the post block FCU, even if the
// optimistic build already sent the same one. The block is verified under one
// CometBFT hash and finalized under another, as happens when a later round
// re-proposes the same payload.
func (s *RPCErrorProxySuite) TestFinalizeBlock_NewPayloadSyncing_SendsFCU() {
	pp := s.prepareForFinalize()

	// Triggers an optimistic build, which sends the FCU for this block.
	processReq := pp.processRequest()
	processReq.NextProposerAddress = pp.proposerAddress
	processReq.Hash = []byte("round-0")
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, processReq)
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.

	s.errProxy.activateSyncingNewPayload()
	defer s.errProxy.deactivate()
	fcuCallsBefore := s.errProxy.fcuCalls.Load()

	finalizeReq := pp.finalizeRequest()
	finalizeReq.Hash = []byte("round-1")
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, finalizeReq)

	s.Require().NoError(err, "FinalizeBlock should accept a SYNCING payload status")
	s.Require().NotNil(finalizeResp)
	s.Require().Contains(s.LogBuffer.String(), "pushed new payload to SYNCING/ACCEPTED node")
	s.Require().Greater(s.errProxy.fcuCalls.Load(), fcuCallsBefore, "FCU must not be skipped after a SYNCING payload")
}

// TestProcessProposal_StartupSync_WaitsForEL shows that the startup FCU sent
// by the first ProcessProposal is not bounded by the validate budget. It
// retries while the EL is unreachable and the proposal is accepted once the
// EL is back.
func (s *RPCErrorProxySuite) TestProcessProposal_StartupSync_WaitsForEL() {
	const (
		firstHeight = 1
		// Longer than the validate budget, so a bounded FCU would give up.
		outage = 2500 * time.Millisecond
	)
	pp := s.prepareProposal(firstHeight)

	s.errProxy.activateDropConn()

	// Bring the EL back after the outage so the startup FCU can succeed.
	go func() {
		time.Sleep(outage)
		s.errProxy.deactivate()
	}()

	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, pp.processRequest())
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	logs := s.LogBuffer.String()
	s.Require().Contains(logs, "Sending startup forkchoice update to execution client")
	s.Require().NotContains(logs, "failed to send force head FCU", "startup FCU must not give up")
}

// TestProcessProposal_ConnectionDrop_Rejects shows that ProcessProposal does
// not wait for an unreachable EL. It retries within the validate budget and
// then rejects the proposal.
func (s *RPCErrorProxySuite) TestProcessProposal_ConnectionDrop_Rejects() {
	pp := s.prepareForFinalize()

	s.errProxy.activateDropConn()
	defer s.errProxy.deactivate()
	callsBefore := s.errProxy.newPayloadCalls.Load()

	start := time.Now()
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, pp.processRequest())
	elapsed := time.Since(start)

	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().Less(elapsed, maxBoundedCall, "ProcessProposal must give up once the budget is spent")
	s.Require().Greater(s.errProxy.newPayloadCalls.Load()-callsBefore, int32(1), "fast failures must be retried")
}

// TestProcessProposal_HungEL_Rejects shows that ProcessProposal does not wait
// for an EL that stalls. With default config the validate budget is below the
// RPC timeout, so the call that timed out is not retried.
func (s *RPCErrorProxySuite) TestProcessProposal_HungEL_Rejects() {
	pp := s.prepareForFinalize()

	s.errProxy.activateHang()
	defer s.errProxy.deactivate()
	callsBefore := s.errProxy.newPayloadCalls.Load()

	start := time.Now()
	processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, pp.processRequest())
	elapsed := time.Since(start)

	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
	s.Require().GreaterOrEqual(elapsed, s.TestNode.EngineClient.GetRPCTimeout(), "the call in flight is not interrupted")
	s.Require().Less(elapsed, maxBoundedCall, "ProcessProposal must give up after the RPC timeout")
	s.Require().Equal(int32(1), s.errProxy.newPayloadCalls.Load()-callsBefore, "a timed out call must not be retried")
}

// TestPrepareProposal_ConnectionDrop_SkipsProposal shows that PrepareProposal
// does not wait for an unreachable EL. It returns an empty proposal once the
// build budget is spent, and proposes again once the EL is back.
func (s *RPCErrorProxySuite) TestPrepareProposal_ConnectionDrop_SkipsProposal() {
	const firstHeight = 1
	nodeAddress := s.initChain()

	// No payload was built for this height yet, so the proposer must send an FCU.
	prepareReq := &types.PrepareProposalRequest{
		Height:          firstHeight,
		Time:            time.Now(),
		ProposerAddress: nodeAddress,
	}
	s.errProxy.activateDropConn()
	fcuCallsBefore := s.errProxy.fcuCalls.Load()

	start := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, prepareReq)
	elapsed := time.Since(start)

	s.Require().NoError(err)
	s.Require().Empty(proposal.Txs, "an unreachable EL must result in no proposal")
	s.Require().Less(elapsed, maxBoundedCall, "PrepareProposal must give up once the budget is spent")
	s.Require().Greater(s.errProxy.fcuCalls.Load()-fcuCallsBefore, int32(1), "fast failures must be retried")

	s.errProxy.deactivate()
	proposal, err = s.SimComet.Comet.PrepareProposal(s.CtxComet, prepareReq)
	s.Require().NoError(err)
	s.Require().Len(proposal.Txs, 2, "the node must propose once the EL is back")
}
