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
/// @author Berachain Team.
/// @dev This contract is a modified version fo the BeaconDepositContract as defined in the
///      Ethereum 2.0 specification. It has been extended to also support trigger withdrawals
///      from the consensus layer and redelegations of stakes.
interface IDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        EVENTS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @notice A processed deposit event.
    event DepositEvent(
        bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index
    );

    /// @notice A processed redelegation event.
    /// @dev We redelegate the `amount` from `pubkey0` to `pubkey1`.
    event RedelegateEvent(bytes pubkey0, bytes pubkey1, bytes signature, bytes index);

    /// @notice A processed withdrawal event.
    /// @dev We withdraw the total amount of the deposit.
    event WithdrawEvent(bytes pubkey, bytes signature);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        DELEGATE                            */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a Phase 0 DepositData object.
     * @param amount The amount of the deposit in Gwei.
     * @param pubkey A BLS12-381 public key.
     * @param withdrawal_credentials Commitment to a public key for withdrawals.
     * @param signature A BLS12-381 signature.
     * @param deposit_data_root The SHA-256 hash of the SSZ-encoded DepositData object.
     * Used as a protection against malformed input.
     */
    function deposit(
        uint256 amount, // in Gwei even if an ERC20 token.
        bytes calldata pubkey,
        bytes calldata withdrawal_credentials,
        bytes calldata signature,
        bytes32 deposit_data_root
    )
        external
        payable;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        REDELEGATE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a redelegation message.
     * @notice This function is only callable by the owner of the stake.
     * @param pubkey0 A BLS12-381 public key.
     * @param pubkey1 A BLS12-381 public key.
     * @param amount The amount of the deposit in Gwei.
     * @param signiture A BLS12-381 signature.
     */
    function redelegate(
        bytes calldata pubkey0,
        bytes calldata pubkey1,
        uint256 amount,
        bytes calldata signiture
    )
        external
        payable;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        WITHDRAW                            */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a withdrawal message.
     * @notice This function is only callable by the owner of the stake.
     * @param pubkey A BLS12-381 public key.
     * @param amount The amount of the deposit in Gwei, even if stake token is an ERC20 token.
     * @param signiture A BLS12-381 signature.
     */
    function withdraw(
        bytes calldata pubkey,
        uint256 amount,
        bytes calldata signiture
    )
        external
        payable;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        READS                               */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Query the current deposit root hash.
     * @return The deposit root hash.
     */
    function get_deposit_root() external view returns (bytes32);

    /**
     * @notice Query the current deposit count.
     * @return The deposit count encoded as a little endian 64-bit number.
     */
    function get_deposit_count() external view returns (bytes32);
}
