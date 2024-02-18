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
 * @dev Interface of the staking module's precompiled contract
 */
contract Staking {
    ////////////////////////////////////////// EVENTS /////////////////////////////////////////////

    /**
     * @dev Emitted by the staking module when `amount` tokens are delegated to
     * `operatorAddress`
     * @param operatorAddress The validator operator address
     * @param amount The amount of tokens delegated
     */
    event Delegate(string operatorAddress, uint256 amount);

    /**
     * @dev Emitted by the staking module when `amount` tokens are unbonded from `validator`
     * @param operatorAddress The validator operator address
     * @param amount The amount of tokens unbonded
     */
    event Undelegate(string operatorAddress, uint256 amount);

    ////////////////////////////////////// WRITE METHODS //////////////////////////////////////////

    /**
     * @dev msg.sender delegates the `amount` of tokens to `operatorAddress`
     * @param operatorAddress The validator operator address
     * @param amount The amount of tokens to delegate
     */
    function delegateFn(string calldata operatorAddress, uint256 amount) external {
        emit Delegate(operatorAddress, amount);
    }

    /**
     * @dev msg.sender undelegates the `amount` of tokens from `operatorAddress`
     * @param operatorAddress The validator operator address
     * @param amount The amount of tokens to undelegate
     */
    function undelegateFn(string calldata operatorAddress, uint256 amount) external {
        emit Undelegate(operatorAddress, amount);
    }
}
