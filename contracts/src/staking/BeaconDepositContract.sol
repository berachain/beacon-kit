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

/// @title BeaconDepositContract
/// @dev This contract is a modified version fo the BeaconDepositContract as
/// defined in the Ethereum 2.0 specification. It has been extended to also
/// support trigger withdrawals from the consensus layer.
/// @author itsdevbear@berachain.com
/// @author po@berachain.com
/// @author ocnc@berachain.com
contract BeaconDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev MINIMUM_DEPOSIT_IN_GWEI is the minimum size of a deposit in Gwei.
    //       1 ether = 1e9 gwei = 1e18 wei
    uint256 private constant MINIMUM_DEPOSIT_IN_GWEI = 1e9;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev InsufficientDeposit is thrown when the specified deposit amount is
    /// below the minimum.
    error InsufficientDeposit();

    /// @dev Deposit is emitted when a deposit is made to the contract.
    ///
    /// @param validatorPubkey The public key of the validator being deposited
    /// to.
    /// @param withdrawalCredentials The withdrawalCredentials for the deposit
    /// @param amount The amount of the deposit in denominated in Gwei.
    event Deposit(
        bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount
    );

    /// @dev Withdrawal is emitted when a withdrawal is made from the contract.
    ///
    /// @param validatorPubkey The public key of the validator being withdrawn
    /// from.
    /// @param withdrawalCredentials The withdrawalCredentials for the
    /// withdrawal
    /// @param amount The amount of the withdrawal denominated in Gwei.
    event Withdrawal(
        bytes validatorPubkey, bytes withdrawalCredentials, uint64 amount
    );

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                    STAKING FUNCTIONS                       */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev msg.sender deposits the `amount` of tokens to `validatorPubkey`.
    /// @param validatorPubkey The validator's public key.
    /// @param amount The amount of tokens to deposit.
    function deposit(
        bytes calldata validatorPubkey,
        uint64 amount
    )
        external
        // TODO: Remove payable
        payable
    {
        // Ensure the deposit amount is above the minimum.
        if (amount < MINIMUM_DEPOSIT_IN_GWEI) {
            revert InsufficientDeposit();
        }

        // TODO: Properly Handle Token Logic.

        // Emit the deposit event.
        emit Deposit(validatorPubkey, abi.encodePacked(msg.sender), amount);
    }

    /// @dev msg.sender withdraws the `amount` of tokens from `validatorPubkey`.
    /// @param validatorPubkey The validator's public key.
    /// @param amount The amount of tokens to undelegate.
    function withdraw(
        bytes calldata validatorPubkey,
        uint64 amount
    )
        external
        payable
    {
        // TODO: Properly Handle Token Logic.

        emit Withdrawal(validatorPubkey, abi.encodePacked(msg.sender), amount);
    }
}
