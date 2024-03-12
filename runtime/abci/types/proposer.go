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

package types

import (
	"cosmossdk.io/x/staking/keeper"
	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExtractProposalPublicKey returns the public key of the proposer of the block.
func ExtractProposalPublicKey(
	ctx sdk.Context,
	stakingKeeper *keeper.Keeper,
	proposerAddressBz []byte,
) [bls12381.PubKeyLength]byte {
	logger := ctx.Logger()

	// Validate Reveal
	v, err := stakingKeeper.ValidatorByConsensusAddress.Get(
		ctx,
		proposerAddressBz,
	)
	if err != nil {
		logger.Warn("failed to get validator", "error", err)
		panic("failed to get validator by consensus address")
	}

	validator, err := stakingKeeper.GetValidator(ctx, v)
	if err != nil {
		logger.Warn("failed to get validator", "error", err)
		panic("failed to get validator by consensus address")
	}

	key, err := validator.CmtConsPublicKey()
	if err != nil {
		logger.Warn("failed to get validator", "error", err)
		panic("failed to get validator by consensus address")
	}

	var pubKey [bls12381.PubKeyLength]byte
	copy(pubKey[:], key.GetBls12381())

	return pubKey
}
