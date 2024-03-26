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

pragma solidity 0.8.25;

/// @title IBeaconDepositContract
/// @author Berachain Team.
/// @dev This contract is used to create validator, deposit and withdraw stake from the Beaconchain.
interface IBeaconDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        EVENTS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when a deposit is made, which could mean a new validator or a top up of an existing one.
     * @param pubkey the public key of the validator who is being deposited for if not a new validator.
     * @param credentials the public key of the operator if new validator or the depositor if top up.
     * @param amount the amount of stake being deposited, in Gwei.
     * @param signature the signature of the deposit message, only checked when creating a new validator.
     * @param index the index of the deposit.
     */
    event Deposit(
        bytes pubkey,
        bytes credentials,
        uint64 amount,
        bytes signature,
        uint64 index
    );

    /**
     * @dev Emitted when a withdrawal is made from a validator.
     * @param fromPubkey The public key of the validator that is being withdrawn from.
     * @param credentials The public key of the account that is withdrawing the stake.
     * @param withdrawalCredentials The public key of the account that will receive the withdrawal.
     * @param amount The amount to be withdrawn from the validator, in Gwei.
     * @param index The index of the withdrawal.
     */
    event Withdrawal(
        bytes fromPubkey,
        bytes credentials,
        bytes withdrawalCredentials,
        uint64 amount,
        uint64 index
    );

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        ERRORS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Error thrown when the deposit amount is too small, to prevent dust deposits.
    error InsufficientDeposit();

    /// @dev Error thrown when the deposit amount is not a multiple of Gwei.
    error DepositNotMultipleOfGwei();

    /// @dev Error thrown when the deposit amount is too high, since it is a uint64.
    error DepositValueTooHigh();

    /// @dev Error thrown when the public key length is not 48 bytes.
    error InvalidPubKeyLength();

    /// @dev Error thrown when the withdrawal credentials length is not 32 bytes.
    error InvalidCredentialsLength();

    /// @dev Error thrown when the signature length is not 96 bytes.
    error InvalidSignatureLength();

    /// @dev Error thrown when the withdrawal amount is too small, to prevent dust withdrawals.
    error InsufficientWithdrawalAmount();

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        WRITES                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a deposit message to the Beaconchain.
     * @notice This will be used to create a new validator or to top up an existing one, increasing stake.
     * @param pubkey is the consensus public key of the validator. If subsequent deposit, its ignored.
     * @param credentials is the staking credentials of the validator. If this is the first deposit it is
     * validator operator public key, if subsequent deposit it is the depositors public key.
     * @param amount is the amount of stake native/ERC20 token to be deposited, in Gwei.
     * @param signature is the signature used only on the first deposit.
     */
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        uint64 amount,
        bytes calldata signature
    )
        external
        payable;

    /**
     * @notice Submit a withdrawal message to the Beaconchain.
     * @notice This function is callable by the account with the stake.
     * @param pubkey is the public key of the validator we are withdrawing from.
     * @param withdrawalCredentials is the public key of the account that will receive the withdrawal.
     * @param amount is the amount of stake to be withdrawn, in Gwei. The amount needs to be calculated offchain since
     * validator tokens are not fungible, and their shares -> stake amount can differ if there is a slashing event.
     */
    function withdraw(
        bytes calldata pubkey,
        bytes calldata withdrawalCredentials,
        uint64 amount
    )
        external;
}
