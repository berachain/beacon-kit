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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package preconf

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	payloadbuilder "github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	gjwt "github.com/golang-jwt/jwt/v5"
)

const (
	// jwtValidityWindow is the time window for JWT validity (iat claim).
	jwtValidityWindow = 5 * time.Minute

	// serverShutdownTimeout is the timeout for graceful server shutdown.
	serverShutdownTimeout = 5 * time.Second

	// serverReadHeaderTimeout is the timeout for reading request headers.
	serverReadHeaderTimeout = 10 * time.Second

	// authHeaderParts is the expected number of parts in the Authorization header.
	authHeaderParts = 2
)

// PayloadProvider is an interface for retrieving payloads by slot and parent block root.
type PayloadProvider interface {
	// GetPayloadBySlot returns the payload for the given slot and parent block root if available.
	GetPayloadBySlot(ctx context.Context, slot math.Slot, parentBlockRoot common.Root) (ctypes.BuiltExecutionPayloadEnv, error)
}

// TLSPaths holds the filesystem paths to a TLS certificate and its private key.
// An empty value means TLS is disabled (plaintext HTTP).
type TLSPaths struct {
	Cert string
	Key  string
}

// Enabled reports whether both a cert and key path are set.
func (p TLSPaths) Enabled() bool {
	return p.Cert != "" && p.Key != ""
}

// SyncChecker exposes the node's sync status for health checks.
// Keeps the original signatures from the cometbft service interface.
type SyncChecker interface {
	// IsAppReady returns nil if the chain is ready (at least one block has been committed).
	// In case of error we set the server as not available.
	IsAppReady() error
	// GetSyncData returns the latest committed height and the target height being synced to.
	GetSyncData() (latestHeight int64, syncToHeight int64)
}

// ELChecker exposes the execution-layer client's connectivity status.
type ELChecker interface {
	// IsConnected returns true if the execution client is reachable.
	IsConnected() bool
}

// Server is the preconf API server that serves GetPayload requests from validators.
type Server struct {
	logger                 log.Logger
	validatorJWTs          ValidatorJWTs
	whitelist              Whitelist
	preconfProposerTracker ProposerTracker
	payloadProvider        PayloadProvider
	syncChecker            SyncChecker
	elChecker              ELChecker
	port                   int
	tlsPaths               TLSPaths
	metrics                *serverMetrics

	// tlsCert holds the currently-served TLS certificate. It is swapped
	// atomically on SIGHUP so cert rotation needs no restart.
	tlsCert atomic.Pointer[tls.Certificate]

	mu         sync.RWMutex
	httpServer *http.Server
}

// NewServer creates a new preconf API server.
func NewServer(
	logger log.Logger,
	validatorJWTs ValidatorJWTs,
	whitelist Whitelist,
	preconfProposerTracker ProposerTracker,
	payloadProvider PayloadProvider,
	syncChecker SyncChecker,
	elChecker ELChecker,
	port int,
	tlsPaths TLSPaths,
	sink TelemetrySink,
) *Server {
	return &Server{
		logger:                 logger,
		validatorJWTs:          validatorJWTs,
		whitelist:              whitelist,
		preconfProposerTracker: preconfProposerTracker,
		payloadProvider:        payloadProvider,
		syncChecker:            syncChecker,
		elChecker:              elChecker,
		port:                   port,
		tlsPaths:               tlsPaths,
		metrics:                newServerMetrics(sink),
	}
}

// Name returns the name of the service.
func (s *Server) Name() string {
	return "preconf-server"
}

// Start starts the preconf API server.
func (s *Server) Start(_ context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(HealthEndpoint, s.handleHealth)
	mux.HandleFunc(PayloadEndpoint, s.handleGetPayload)

	addr := fmt.Sprintf(":%d", s.port)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	tlsEnabled := s.tlsPaths.Enabled()
	if tlsEnabled {
		// Load the cert once up front so a malformed cert/key fails Start
		// synchronously rather than only logging inside the serve goroutine.
		if err := s.loadTLSCert(); err != nil {
			return errors.Wrap(err, "failed to load TLS certificate")
		}
		// GetCertificate is consulted on every handshake, so swapping the
		// stored cert on SIGHUP takes effect without a restart.
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
				return s.tlsCert.Load(), nil
			},
		}
	}

	// Bind synchronously so a failure (port in use, permission denied) fails
	// Start and propagates to the service registry, rather than only logging
	// inside the serve goroutine with a dead endpoint left behind.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrapf(err, "failed to bind preconf API server to %s", addr)
	}

	s.mu.Lock()
	s.httpServer = server
	s.mu.Unlock()

	s.logger.Info("Starting preconf API server",
		"address", addr,
		"tls_enabled", tlsEnabled,
		"num_validator_jwts", len(s.validatorJWTs),
	)

	// Log the registered validator pubkeys for debugging
	for pubkey := range s.validatorJWTs {
		s.logger.Info("Registered validator JWT", "pubkey", pubkey.String())
	}

	go func() {
		var serveErr error
		if tlsEnabled {
			// Cert and key are supplied via TLSConfig.GetCertificate.
			serveErr = server.ServeTLS(ln, "", "")
		} else {
			serveErr = server.Serve(ln)
		}
		// A graceful Stop() makes Serve return http.ErrServerClosed, which is
		// expected. Anything else is a genuine serve failure worth logging.
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			s.logger.Error("Preconf API server error", "error", serveErr)
		}
	}()

	return nil
}

// OnSIGHUP implements SIGHUPHandler. It hot-reloads the whitelist and, when TLS
// is configured, the server certificate from disk.
func (s *Server) OnSIGHUP() {
	if err := s.whitelist.Reload(); err != nil {
		s.logger.Error("Failed to reload preconf whitelist", "error", err)
	} else {
		s.logger.Info("Preconf whitelist reloaded", "whitelist_count", s.whitelist.Len())
	}

	if !s.tlsPaths.Enabled() {
		return
	}
	// On a parse failure, keep serving the existing cert rather than break the
	// listener (same policy as the whitelist reload above).
	if err := s.loadTLSCert(); err != nil {
		s.logger.Error("Failed to reload TLS certificate", "error", err)
		return
	}
	s.logger.Info("TLS certificate reloaded", "cert", s.tlsPaths.Cert)
}

// loadTLSCert reads the cert/key from disk and atomically stores them as the
// certificate served on subsequent handshakes.
func (s *Server) loadTLSCert() error {
	cert, err := tls.LoadX509KeyPair(s.tlsPaths.Cert, s.tlsPaths.Key)
	if err != nil {
		return err
	}
	s.tlsCert.Store(&cert)
	return nil
}

// Stop stops the preconf API server.
func (s *Server) Stop() error {
	s.mu.Lock()
	server := s.httpServer
	s.httpServer = nil
	s.mu.Unlock()

	if server == nil {
		return nil
	}

	s.logger.Info("Stopping preconf API server")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	return server.Shutdown(ctx)
}

// handleHealth checks sync status and returns 200 when the sequencer is synced
// and ready to produce blocks, or 503 when it is not.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	resp := s.buildHealthResponse()

	w.Header().Set("Content-Type", "application/json")
	healthy := resp.IsReady && !resp.IsSyncing && resp.ELConnected
	if !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Error("Failed to encode health response", "error", err)
	}
}

// buildHealthResponse inspects the node's sync state and EL connectivity
// and produces a HealthResponse.
func (s *Server) buildHealthResponse() *HealthResponse {
	resp := new(HealthResponse)

	if s.syncChecker != nil {
		resp.IsReady = s.syncChecker.IsAppReady() == nil
		latestHeight, syncToHeight := s.syncChecker.GetSyncData()
		resp.IsSyncing = syncToHeight > latestHeight
	}

	if s.elChecker != nil {
		resp.ELConnected = s.elChecker.IsConnected()
	}

	return resp
}

// handleGetPayload handles the GetPayload endpoint.
func (s *Server) handleGetPayload(w http.ResponseWriter, r *http.Request) {
	result := ServerResultOK
	defer func() { s.metrics.markPayloadRequest(result) }()

	// Only accept POST requests
	if r.Method != http.MethodPost {
		result = ServerResultMethodNotAllowed
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Validate JWT and extract validator pubkey
	pubkey, err := s.validateJWT(r)
	if err != nil {
		result = ServerResultUnauthorized
		s.logger.Warn("JWT validation failed", "error", err)
		s.writeError(w, http.StatusUnauthorized, "unauthorized: "+err.Error())
		return
	}

	// Check if validator is whitelisted
	if s.whitelist != nil && !s.whitelist.IsWhitelisted(pubkey) {
		result = ServerResultNotWhitelisted
		s.logger.Warn("Validator not whitelisted", "pubkey", pubkey)
		s.writeError(w, http.StatusForbidden, "validator not whitelisted")
		return
	}

	// Parse request body
	var req GetPayloadRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		result = ServerResultBadRequest
		s.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Verify the requesting validator is the expected proposer for this slot.
	matched := s.preconfProposerTracker.IsExpectedProposer(req.Slot, pubkey)
	s.metrics.markProposerCheck(matched)
	if !matched {
		result = ServerResultWrongProposer
		s.logger.Warn("Validator is not the expected proposer for slot",
			"pubkey", pubkey.String(),
			"slot", req.Slot,
		)
		s.writeError(w, http.StatusForbidden, "not the expected proposer for this slot")
		return
	}

	s.logger.Info("Preconf server received payload request",
		"slot", req.Slot,
		"validator_pubkey", pubkey.String(),
	)

	// Fetch the payload via payloadProvider and write the HTTP response to the caller.
	result = s.fetchAndWritePayload(r.Context(), w, req)
}

func (s *Server) fetchAndWritePayload(ctx context.Context, w http.ResponseWriter, req GetPayloadRequest) ServerResult {
	startTime := time.Now()
	envelope, err := s.payloadProvider.GetPayloadBySlot(ctx, req.Slot, req.ParentBlockRoot)
	elapsed := time.Since(startTime)
	if err != nil {
		s.logger.Warn("Failed to get payload",
			"slot", req.Slot,
			"error", err,
			"elapsed", elapsed,
		)
		if errors.Is(err, payloadbuilder.ErrPayloadIDNotFound) || errors.Is(err, engineerrors.ErrUnknownPayload) {
			s.writeError(w, http.StatusNotFound, "payload not available: "+err.Error())
			return ServerResultPayloadNotFound
		}
		s.writeError(w, http.StatusInternalServerError, "internal server error")
		return ServerResultInternalError
	}

	if envelope == nil {
		s.writeError(w, http.StatusNotFound, "payload not available")
		return ServerResultPayloadNotFound
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(NewGetPayloadResponseFromEnvelope(envelope)); err != nil {
		s.logger.Warn("Failed to write payload response", "slot", req.Slot, "error", err)
		return ServerResultResponseWriteError
	}

	s.logger.Info("GetPayloadBySlot completed", "slot", req.Slot, "elapsed", elapsed)
	return ServerResultOK
}

// validateJWT validates the JWT token from the Authorization header and returns
// the validator pubkey associated with the token.
func (s *Server) validateJWT(r *http.Request) (crypto.BLSPubkey, error) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return crypto.BLSPubkey{}, errors.New("missing Authorization header")
	}

	// Expect "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", authHeaderParts)
	if len(parts) != authHeaderParts || strings.ToLower(parts[0]) != "bearer" {
		return crypto.BLSPubkey{}, errors.New("invalid Authorization header format")
	}
	tokenString := parts[1]

	// Try to validate against each validator's JWT secret
	for pubkey, secret := range s.validatorJWTs {
		if s.verifyToken(tokenString, secret) {
			return pubkey, nil
		}
	}

	return crypto.BLSPubkey{}, errors.New("invalid or unknown JWT token")
}

// verifyToken verifies a JWT token against a secret.
func (s *Server) verifyToken(tokenString string, secret *jwt.Secret) bool {
	token, err := gjwt.Parse(tokenString, func(token *gjwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*gjwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret.Bytes(), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Validate iat claim (issued at) - token should be recent
	claims, ok := token.Claims.(gjwt.MapClaims)
	if !ok {
		return false
	}

	iat, err := claims.GetIssuedAt()
	if err != nil || iat == nil {
		return false
	}

	// Check if token was issued within the validity window
	now := time.Now()
	if now.Sub(iat.Time) > jwtValidityWindow || iat.Time.After(now.Add(jwtValidityWindow)) {
		return false
	}

	return true
}

// Handler returns an http.Handler for testing purposes.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(HealthEndpoint, s.handleHealth)
	mux.HandleFunc(PayloadEndpoint, s.handleGetPayload)
	return mux
}

// writeError writes an error response.
func (s *Server) writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Error("Failed to encode error response", "error", err)
	}
}
