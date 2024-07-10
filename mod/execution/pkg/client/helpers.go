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
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/net/jwt"
	gjwt "github.com/golang-jwt/jwt/v5"
)

// createContextWithTimeout creates a context with a timeout and returns it
// along with the cancel function.
func (s *EngineClient[
	_, _, _, _, _,
]) createContextWithTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	startTime := time.Now()
	dctx, cancel := context.WithTimeoutCause(
		ctx,
		s.cfg.RPCTimeout,
		engineerrors.ErrEngineAPITimeout,
	)
	s.metrics.measureNewPayloadDuration(startTime)
	return dctx, cancel
}

// processPayloadStatusResult processes the payload status result and
// returns the latest valid hash or an error.
func processPayloadStatusResult(
	result *engineprimitives.PayloadStatusV1,
) (*gethprimitives.ExecutionHash, error) {
	switch result.Status {
	case engineprimitives.PayloadStatusAccepted:
		return nil, engineerrors.ErrAcceptedPayloadStatus
	case engineprimitives.PayloadStatusSyncing:
		return nil, engineerrors.ErrSyncingPayloadStatus
	case engineprimitives.PayloadStatusInvalid:
		return result.LatestValidHash, engineerrors.ErrInvalidPayloadStatus
	case engineprimitives.PayloadStatusValid:
		return result.LatestValidHash, nil
	default:
		return nil, engineerrors.ErrUnknownPayloadStatus
	}
}

// buildSignedJWT builds a signed JWT from the provided JWT secret.
func buildSignedJWT(s *jwt.Secret) (string, error) {
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, gjwt.MapClaims{
		"iat": &gjwt.NumericDate{Time: time.Now()},
	})
	str, err := token.SignedString(s[:])
	if err != nil {
		return "", errors.Newf("failed to create JWT token: %w", err)
	}
	return str, nil
}
