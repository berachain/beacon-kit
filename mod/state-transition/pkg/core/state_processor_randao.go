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

package core

// processRandaoReveal processes the randao reveal and
// ensures it matches the local state.
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, ContextT,
	DepositT, ExecutionPayloadT,
	ForkDataT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
]) processRandaoReveal(
	st BeaconStateT,
	blk BeaconBlockT,
	skipVerification bool,
) error {
	return sp.rp.ProcessRandao(st, blk, skipVerification)
}

// processRandaoMixesReset as defined in the Ethereum 2.0 specification.
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#randao-mixes-updates
//
//nolint:lll
func (sp *StateProcessor[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, ContextT,
	DepositT, ExecutionPayloadT,
	ForkDataT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
]) processRandaoMixesReset(
	st BeaconStateT,
) error {
	return sp.rp.ProcessRandaoMixesReset(st)
}
