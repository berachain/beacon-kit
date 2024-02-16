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

pragma solidity 0.8.24;

/**
 * @dev Interface of the staking contract.
 */
contract Staking {
    //////////////////////////////////////// VARIABLES
    // /////////////////////////////////////////////
    uint64 nonce;

    ////////////////////////////////////////// EVENTS /////////////////////////////////////////////
    /**
     * @dev Emitted by the staking contract when `amount` tokens
     * are delegated to `validatorPubkey`.
     * @param validatorPubkey The validator's public key.
     * @param withdrawalCredentials The withdrawal credentials of the validator.
     * @param amount The amount of tokens delegated.
     * @param nonce The nonce of the delegation.
     */
    event Delegate(
        bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount, uint64 nonce
    );

    /**
     * @dev Emitted by the staking contract when `amount` tokens are unbonded from
     * `validatorPubkey`.
     * @param validatorPubkey The validator's public key.
     * @param amount The amount of tokens unbonded.
     * @param nonce The nonce of the undelegation.
     */
    event Undelegate(bytes validatorPubkey, uint64 amount, uint64 nonce);

    ////////////////////////////////////// WRITE METHODS //////////////////////////////////////////

    /**
     * @dev msg.sender delegates the `amount` of tokens to `validatorPubkey`.
     * @param validatorPubkey The validator's public key.
     * @param amount The amount of tokens to delegate.
     */
    function delegateFn(
        bytes calldata validatorPubkey,
        bytes calldata withdrawalCredentials,
        uint64 amount
    )
        external
    {
        emit Delegate(validatorPubkey, withdrawalCredentials, amount, nonce);
        nonce++;
    }

    /**
     * @dev msg.sender undelegates the `amount` of tokens from `validatorPubkey`.
     * @param validatorPubkey The validator's public key.
     * @param amount The amount of tokens to undelegate.
     */
    function undelegateFn(bytes calldata validatorPubkey, uint64 amount) external {
        emit Undelegate(validatorPubkey, amount, nonce);
        nonce++;
    }
}
