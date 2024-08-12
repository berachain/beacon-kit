// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package client

import (
	"context"
	"net/http"
	"time"
)

// jwtRefreshLoop refreshes the JWT token for the execution client.
func (s *EngineClient[
	_, _,
]) jwtRefreshLoop(
	ctx context.Context,
) {
	s.logger.Info("Starting JWT refresh loop 🔄")
	ticker := time.NewTicker(s.cfg.RPCJWTRefreshInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			if err := s.dialExecutionRPCClient(ctx); err != nil {
				s.logger.Error(
					"Failed to refresh engine auth token",
					"err",
					err,
				)
			}
		}
	}
}

// buildJWTHeader builds an http.Header that has the JWT token
// attached for authorization.
func (s *EngineClient[
	_, _,
]) buildJWTHeader() (http.Header, error) {
	header := make(http.Header)

	// Build the JWT token.
	token, err := buildSignedJWT(s.jwtSecret)
	if err != nil {
		s.logger.Error("Failed to build JWT token", "err", err)
		return header, err
	}

	// Add the JWT token to the headers.
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+token)
	return header, nil
}
