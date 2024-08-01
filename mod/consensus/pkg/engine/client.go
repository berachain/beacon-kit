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

package engine

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	abci "github.com/cometbft/cometbft/abci/types"
)

// TODO: decouple from the proto types
type Client interface {
	// InitChain is called when the blockchain is first started
	// It returns the validator set and the app hash.
	InitChain(
		ctx context.Context,
		genesisBz []byte,
	) (transition.ValidatorUpdates, error) // Initialize blockchain w validators/other info from CometBFT

	// PrepareProposal is called when a proposal is made.
	// It returns the txs to be executed in the proposal.
	PrepareProposal(
		ctx context.Context,
		req *abci.PrepareProposalRequest,
	) ([][]byte, error)

	// ProcessProposal is called when a proposal is processed.
	// It returns an error if the proposal is invalid.
	ProcessProposal(
		ctx context.Context,
		req *abci.ProcessProposalRequest,
	) error

	// Deliver the decided block with its txs to the Application
	FinalizeBlock(
		ctx context.Context,
		req *abci.FinalizeBlockRequest,
	) (transition.ValidatorUpdates, error)

	// TODO: snapshot methods
}
