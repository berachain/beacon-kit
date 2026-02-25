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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
)

const (
	// tokenCacheExpiry is how long to cache JWT tokens before regenerating.
	tokenCacheExpiry = 30 * time.Second
)

var (
	// ErrSequencerUnavailable is returned when the sequencer is not reachable.
	ErrSequencerUnavailable = errors.New("sequencer unavailable")

	// ErrPayloadNotFound is returned when no payload is available for the requested slot.
	ErrPayloadNotFound = errors.New("payload not found")
)

// Client is an HTTP client for fetching payloads from the sequencer.
type Client struct {
	logger       log.Logger
	httpClient   *http.Client
	sequencerURL string
	jwtSecret    *jwt.Secret
	timeout      time.Duration

	// mu protects the JWT token cache
	mu          sync.RWMutex
	cachedToken string
	tokenExpiry time.Time
}

// NewClient creates a new preconf client for fetching payloads from the sequencer.
func NewClient(
	logger log.Logger,
	sequencerURL string,
	jwtSecret *jwt.Secret,
	timeout time.Duration,
) *Client {
	return &Client{
		logger:       logger,
		sequencerURL: sequencerURL,
		jwtSecret:    jwtSecret,
		timeout:      timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetPayloadBySlot fetches a payload from the sequencer for the given slot and parent block root.
func (c *Client) GetPayloadBySlot(
	ctx context.Context,
	slot math.Slot,
	parentBlockRoot common.Root,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	// Build request
	reqBody := GetPayloadRequest{Slot: slot, ParentBlockRoot: parentBlockRoot}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := c.sequencerURL + PayloadEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add JWT authorization
	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	c.logger.Debug("Fetching payload from sequencer", "slot", slot, "url", url)

	resp, err := c.httpClient.Do(req) //#nosec:G704 // sequencer URL from trusted config
	if err != nil {
		return nil, errors.Wrapf(ErrSequencerUnavailable, "request failed: %v", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, errors.New("received nil response from sequencer")
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if jsonErr := json.Unmarshal(body, &errResp); jsonErr == nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, errors.Wrapf(ErrPayloadNotFound, "slot %d: %s", slot, errResp.Message)
			}
			return nil, fmt.Errorf("sequencer error (status %d): %s", resp.StatusCode, errResp.Message)
		}
		return nil, fmt.Errorf("sequencer error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse successful response
	var payloadResp GetPayloadResponse
	if err = json.Unmarshal(body, &payloadResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var blockHash interface{}
	if payloadResp.ExecutionPayload != nil {
		blockHash = payloadResp.ExecutionPayload.GetBlockHash()
	} else {
		blockHash = "nil"
	}

	c.logger.Info("Successfully fetched payload from sequencer",
		"slot", slot,
		"block_hash", blockHash,
	)

	return payloadResp.ToExecutionPayloadEnvelope(), nil
}

// getToken returns a valid JWT token, generating a new one if necessary.
func (c *Client) getToken() (string, error) {
	c.mu.RLock()
	if c.cachedToken != "" && time.Now().Before(c.tokenExpiry) {
		token := c.cachedToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Generate new token
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.cachedToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.cachedToken, nil
	}

	token, err := c.jwtSecret.BuildSignedToken()
	if err != nil {
		return "", err
	}

	// Cache token with shorter expiry than validity window
	c.cachedToken = token
	c.tokenExpiry = time.Now().Add(tokenCacheExpiry)

	return token, nil
}
