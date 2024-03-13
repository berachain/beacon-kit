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

package staking

import (
	"context"

	bls12381 "github.com/berachain/beacon-kit/crypto/bls12-381"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
)

// GetValidatorPubkeyFromConsAddress returns the public
// key of the proposer of the block.
func (k *Keeper) GetValidatorPubkeyFromConsAddress(
	ctx context.Context,
	consAddr []byte,
) ([bls12381.PubKeyLength]byte, error) {
	valAddr, err := k.stakingKeeper.ValidatorByConsensusAddress.Get(
		ctx,
		consAddr,
	)
	if err != nil {
		return [bls12381.PubKeyLength]byte{}, err
	}

	return k.GetValidatorPubkeyFromValAddress(ctx, valAddr)
}

// GetValidatorPubkeyFromValAddress returns the public
// key of the validator with the given validator address.
func (k *Keeper) GetValidatorPubkeyFromValAddress(
	ctx context.Context,
	valAddr []byte,
) ([bls12381.PubKeyLength]byte, error) {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return [bls12381.PubKeyLength]byte{}, err
	}

	var key cmtprotocrypto.PublicKey
	key, err = validator.CmtConsPublicKey()
	if err != nil {
		return [bls12381.PubKeyLength]byte{}, err
	}

	var pubKey [bls12381.PubKeyLength]byte
	copy(pubKey[:], key.GetBls12381())

	return pubKey, nil
}
