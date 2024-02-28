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

/// @title IDepositContract
/// @author Berachain Team.
/// @dev This contract is the interdface of the deposit contract that the beaconkit uses to handle its
/// delegate proof of stake system. It is derived from the Ethereum 2.0 specification but with an arbitrary
/// staking backend for the chain, hence it is very flexible and can be used with ERC20s or the native token.
interface IDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        EVENTS                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @dev Emitted when a deposit is made.
     * @notice the withdrawCredentials can be left empty if its not the first deposit.
     * @param amount The amount of the deposit.
     * @param pubKey The public key of the validator.
     * @param withdrawalCredentials The withdrawal credentials of the validator.
     * @param signature The signature of the depositor.
     */
    event Deposit(
        uint64 amount,
        bytes pubKey,
        bytes withdrawalCredentials,
        bytes signature
    );

    /**
     * @dev Emitted when a redirection is made.
     * @param srcPub The public key of the source validator.
     * @param dstPub The public key of the destination validator.
     * @param amount The amount to be redirected.
     * @param signature The signature of the redirector with the stake.
     */
    event Redirect(bytes srcPub, bytes dstPub, uint64 amount, bytes signature);

    /**
     * @dev Emitted when a withdrawal is made.
     * @param pubKey The public key of the validator.
     * @param amount The amount to be withdrawn from the validator.
     * @param signature The signature of the withdrawer.
     */
    event Withdraw(bytes pubKey, uint64 amount, bytes signature);

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        WRITES                              */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /**
     * @notice Submit a stake message to the Beaconchain.
     * @param amount The amount of the deposit.
     * @param pubkey A BLS12-381 public key.
     * @param withdrawal_credentials Commitment to a public key for withdrawals.
     * @param signature A BLS12-381 signature from the depositor.
     */
    function deposit(
        uint64 amount, // in Gwei even if an ERC20 token.
        bytes calldata pubkey,
        bytes calldata withdrawal_credentials,
        bytes calldata signature
    )
        external
        payable;

    /**
     * @notice Submit a redirect stake message.
     * @notice This function is only callable by the owner of the stake.
     * @param srcPub A BLS12-381 public key of the source validator.
     * @param dstPub A BLS12-381 public key of the destination validator.
     * @param amount The amount of the deposit.
     * @param signiture A BLS12-381 signature from the redirector.
     */
    function redirect(
        bytes calldata srcPub,
        bytes calldata dstPub,
        uint64 amount,
        bytes calldata signiture
    )
        external
        payable;

    /**
     * @notice Submit a withdrawal message to the Beaconchain.
     * @notice This function is only callable by the owner of the stake.
     * @param pubkey A BLS12-381 public key.
     * @param amount The amount of the deposit to be withdrawn.
     * @param signiture A BLS12-381 signature of the withdrawer.
     */
    function withdraw(
        bytes calldata pubkey,
        uint64 amount,
        bytes calldata signiture
    )
        external
        payable;
}
