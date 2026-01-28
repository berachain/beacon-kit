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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
	gjwt "github.com/golang-jwt/jwt/v5"
)

const (
	// jwtValidityWindow is the time window for JWT validity (iat claim).
	jwtValidityWindow = 60 * time.Second

	// serverShutdownTimeout is the timeout for graceful server shutdown.
	serverShutdownTimeout = 5 * time.Second

	// serverReadHeaderTimeout is the timeout for reading request headers.
	serverReadHeaderTimeout = 10 * time.Second

	// authHeaderParts is the expected number of parts in the Authorization header.
	authHeaderParts = 2
)

// PayloadProvider is an interface for retrieving payloads by slot.
type PayloadProvider interface {
	// GetPayloadBySlot returns the payload for the given slot if available.
	GetPayloadBySlot(ctx context.Context, slot math.Slot) (ctypes.BuiltExecutionPayloadEnv, error)
	// GetExpectedProposer returns the expected proposer for the given slot.
	GetExpectedProposer(slot math.Slot) (crypto.BLSPubkey, bool)
}

// Server is the preconf API server that serves GetPayload requests from validators.
type Server struct {
	logger          log.Logger
	validatorJWTs   ValidatorJWTs
	whitelist       Whitelist
	payloadProvider PayloadProvider
	port            int

	mu         sync.RWMutex
	httpServer *http.Server
}

// NewServer creates a new preconf API server.
func NewServer(
	logger log.Logger,
	validatorJWTs ValidatorJWTs,
	whitelist Whitelist,
	payloadProvider PayloadProvider,
	port int,
) *Server {
	return &Server{
		logger:          logger,
		validatorJWTs:   validatorJWTs,
		whitelist:       whitelist,
		payloadProvider: payloadProvider,
		port:            port,
	}
}

// Name returns the name of the service.
func (s *Server) Name() string {
	return "preconf-server"
}

// Start starts the preconf API server.
func (s *Server) Start(_ context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(PayloadEndpoint, s.handleGetPayload)

	addr := fmt.Sprintf(":%d", s.port)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: serverReadHeaderTimeout,
	}

	s.mu.Lock()
	s.httpServer = server
	s.mu.Unlock()

	s.logger.Info("Starting preconf API server", "address", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Preconf API server error", "error", err)
		}
	}()

	return nil
}

// Stop stops the preconf API server.
func (s *Server) Stop() error {
	s.mu.RLock()
	server := s.httpServer
	s.mu.RUnlock()

	if server == nil {
		return nil
	}

	s.logger.Info("Stopping preconf API server")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	return server.Shutdown(ctx)
}

// handleGetPayload handles the GetPayload endpoint.
func (s *Server) handleGetPayload(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Validate JWT and extract validator pubkey
	pubkey, err := s.validateJWT(r)
	if err != nil {
		s.logger.Warn("JWT validation failed", "error", err)
		s.writeError(w, http.StatusUnauthorized, "unauthorized: "+err.Error())
		return
	}

	// Check if validator is whitelisted
	if s.whitelist != nil && !s.whitelist.IsWhitelisted(pubkey) {
		s.logger.Warn("Validator not whitelisted", "pubkey", pubkey)
		s.writeError(w, http.StatusForbidden, "validator not whitelisted")
		return
	}

	// Parse request body
	var req GetPayloadRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Validate that the requesting validator is the expected proposer for this slot
	expectedProposer, found := s.payloadProvider.GetExpectedProposer(req.Slot)
	if !found {
		s.logger.Warn("No expected proposer found for slot", "slot", req.Slot)
		s.writeError(w, http.StatusNotFound, "no payload building in progress for slot")
		return
	}
	if expectedProposer != pubkey {
		s.logger.Warn("Validator is not the expected proposer",
			"slot", req.Slot,
			"expected", expectedProposer.String()[:16]+"...",
			"actual", pubkey.String()[:16]+"...",
		)
		s.writeError(w, http.StatusForbidden, "validator is not the expected proposer for this slot")
		return
	}

	// Get the payload from provider
	ctx := r.Context()
	envelope, err := s.payloadProvider.GetPayloadBySlot(ctx, req.Slot)
	if err != nil {
		s.logger.Warn("Failed to get payload", "slot", req.Slot, "error", err)
		s.writeError(w, http.StatusNotFound, "payload not available: "+err.Error())
		return
	}

	// Convert to response
	resp := NewGetPayloadResponseFromEnvelope(envelope)

	s.logger.Info("Serving payload to validator",
		"slot", req.Slot,
		"validator", pubkey.String()[:16]+"...",
	)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		s.logger.Error("Failed to encode response", "error", err)
	}
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
