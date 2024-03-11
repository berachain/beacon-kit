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

package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	stakingtypes "cosmossdk.io/x/staking/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The purpose of UnimplementedStakingHooks is to reduce the boilerplate in
// modules that want to implemente hooks.
type UnimplementedStakingHooks struct{}

// UnimplementedStakingHooks implements stakingtypes.StakingHooks.
var _ stakingtypes.StakingHooks = &UnimplementedStakingHooks{}

// AfterValidatorCreated does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterValidatorCreated(
	context.Context,
	sdk.ValAddress,
) error {
	return nil
}

// BeforeValidatorModified does nothing and returns nil.
func (*UnimplementedStakingHooks) BeforeValidatorModified(
	context.Context,
	sdk.ValAddress,
) error {
	return nil
}

// AfterValidatorRemoved does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterValidatorRemoved(
	context.Context,
	sdk.ConsAddress,
	sdk.ValAddress,
) error {
	return nil
}

// BeforeDelegationCreated does nothing and returns nil.
func (*UnimplementedStakingHooks) BeforeDelegationCreated(
	context.Context,
	sdk.AccAddress,
	sdk.ValAddress,
) error {
	return nil
}

// BeforeDelegationSharesModified does nothing and returns nil.
func (*UnimplementedStakingHooks) BeforeDelegationSharesModified(
	context.Context,
	sdk.AccAddress,
	sdk.ValAddress,
) error {
	return nil
}

// AfterDelegationModified does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterDelegationModified(
	context.Context,
	sdk.AccAddress,
	sdk.ValAddress,
) error {
	return nil
}

// BeforeValidatorSlashed does nothing and returns nil.
func (*UnimplementedStakingHooks) BeforeValidatorSlashed(
	context.Context,
	sdk.ValAddress,
	sdkmath.LegacyDec,
) error {
	return nil
}

// AfterValidatorBonded does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterValidatorBonded(
	context.Context,
	sdk.ConsAddress,
	sdk.ValAddress,
) error {
	return nil
}

// AfterValidatorBeginUnbonding does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterValidatorBeginUnbonding(
	context.Context,
	sdk.ConsAddress,
	sdk.ValAddress,
) error {
	return nil
}

// BeforeDelegationRemoved does nothing and returns nil.
func (*UnimplementedStakingHooks) BeforeDelegationRemoved(
	context.Context,
	sdk.AccAddress,
	sdk.ValAddress,
) error {
	return nil
}

// AfterUnbondingInitiated does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterUnbondingInitiated(
	context.Context,
	uint64,
) error {
	return nil
}

// AfterConsensusPubKeyUpdate does nothing and returns nil.
func (*UnimplementedStakingHooks) AfterConsensusPubKeyUpdate(
	context.Context,
	cryptotypes.PubKey,
	cryptotypes.PubKey,
	sdk.Coin,
) error {
	return nil
}
