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
	"math/big"

	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// HeaderByNumber retrieves the block header by its number.
func (s *EngineClient[
	_, _,
]) HeaderByNumber(
	ctx context.Context,
	number *big.Int,
) (*gethprimitives.Header, error) {
	// Infer the latest height if the number is nil.
	if number == nil {
		latest, err := s.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}
		number = new(big.Int).SetUint64(latest)
	}

	// Check the cache for the header.
	if header, ok := s.engineCache.HeaderByNumber(number.Uint64()); ok {
		return header, nil
	}

	header, err := s.Client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	// Add the header to the cache before returning.
	s.engineCache.AddHeader(header)

	return header, nil
}

// HeaderByHash retrieves the block header by its hash.
func (s *EngineClient[
	_, _,
]) HeaderByHash(
	ctx context.Context,
	hash common.ExecutionHash,
) (*gethprimitives.Header, error) {
	// Check the cache for the header.
	header, ok := s.engineCache.HeaderByHash(hash)
	if ok {
		return header, nil
	}
	header, err := s.Client.HeaderByHash(
		ctx,
		gethprimitives.ExecutionHash(hash),
	)
	if err != nil {
		return nil, err
	}
	s.engineCache.AddHeader(header)
	return header, nil
}
