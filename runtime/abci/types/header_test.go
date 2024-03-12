// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types_test

import (
	"testing"
	"time"

	"github.com/berachain/beacon-kit/runtime/abci/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmtversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cometbft "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	cmtHeader := cometbft.Header{
		Version: cmtversion.Consensus{Block: 1, App: 2},
		ChainID: "chainId",
		Height:  3,
		Time:    time.Date(2024, 3, 12, 15, 8, 15, 0, time.Local),
		LastBlockID: cometbft.BlockID{
			Hash: tmhash.Sum([]byte("last_block_id")),
			PartSetHeader: cometbft.PartSetHeader{
				Total: 6,
				Hash:  tmhash.Sum([]byte("parts_header_hash")),
			},
		},
		LastCommitHash:     tmhash.Sum([]byte("last_commit_hash")),
		DataHash:           tmhash.Sum([]byte("data_hash")),
		ValidatorsHash:     tmhash.Sum([]byte("validators_hash")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators_hash")),
		ConsensusHash:      tmhash.Sum([]byte("consensus_hash")),
		AppHash:            tmhash.Sum([]byte("app_hash")),
		LastResultsHash:    tmhash.Sum([]byte("last_results_hash")),
		EvidenceHash:       tmhash.Sum([]byte("evidence_hash")),
		ProposerAddress:    crypto.AddressHash([]byte("proposer_address")),
	}
	require.Equal(t,
		"B838F0F31AE754C60FBD7FC48CED2FCAFAE5B389D085D4DF29346EA0A07139CB",
		cmtHeader.Hash().String(),
	)

	header := types.CometBFTHeader{}
	header.FromCometBFT(cmtHeader)

	otherCmtHeader := header.ToCometBFT()
	// Time in otherCmtHeader is in UTC.
	cmtHeader.Time = cmtHeader.Time.UTC()
	require.Equal(t, cmtHeader, otherCmtHeader)
	require.Equal(t, cmtHeader.Hash(), otherCmtHeader.Hash())
}
