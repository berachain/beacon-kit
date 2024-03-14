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

import { BeaconDepositContract } from "../BeaconDepositContract.sol";
import { IStakeERC20 } from "./IStakeERC20.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits, withdrawals, and redirections of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 */
contract ERC20BeaconDepositContract is BeaconDepositContract {
    /// @notice The ERC20 token contract that is used for staking.
    /// TODO: Change this to the actual ERC20 token contract.
    address public ERC20_DEPOSIT_ASSET;

    /// @dev Validates the deposit amount and burns the staking asset from the sender.
    function _deposit(uint64 amount) internal override returns (uint64) {
        // burn the staking asset from the sender, converting the gwei to wei.
        IStakeERC20(ERC20_DEPOSIT_ASSET).burn(msg.sender, uint256(amount) * 1e9);
        return amount;
    }
}
