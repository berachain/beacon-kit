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

import { IBeaconDepositContract } from "./IBeaconDepositContract.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits and withdrawals of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 */
contract BeaconDepositContract is IBeaconDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The minimum amount of stake that can be deposited to prevent dust.
    /// @dev This is 32 ether in Gwei since our deposit contract denominates in Gwei. 32e9 * 1e9 = 32e18.
    uint64 private constant MIN_DEPOSIT_AMOUNT_IN_GWEI = 32e9;

    /// @dev The minimum amount of stake that can be withdrawn to prevent dust.
    /// leaving the buffer for their deposit to be slashed.
    uint256 private constant MIN_WITHDRAWAL_AMOUNT_IN_GWEI =
        MIN_DEPOSIT_AMOUNT_IN_GWEI / 10;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 private constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 private constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 private constant CREDENTIALS_LENGTH = 32;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    uint64 public depositCount;
    /// @dev withdrawalCount represents the number of withdrawals that
    /// have been requested.
    uint64 public withdrawalCount;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        uint64 amount,
        bytes calldata signature
    )
        external
        payable
    {
        if (pubkey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (credentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (signature.length != SIGNATURE_LENGTH) {
            revert InvalidSignatureLength();
        }

        uint64 amountInGwei = _deposit(amount);

        if (amountInGwei < MIN_DEPOSIT_AMOUNT_IN_GWEI) {
            revert InsufficientDeposit();
        }

        unchecked {
            // slither-disable-next-line reentrancy-benign,reentrancy-events
            emit Deposit(
                pubkey, credentials, amountInGwei, signature, depositCount++
            );
        }
    }

    /// @inheritdoc IBeaconDepositContract
    function withdraw(
        bytes calldata pubkey,
        bytes calldata withdrawalCredentials,
        uint64 amount
    )
        external
    {
        if (pubkey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (withdrawalCredentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (amount < MIN_WITHDRAWAL_AMOUNT_IN_GWEI) {
            revert InsufficientWithdrawalAmount();
        }
        unchecked {
            emit Withdrawal(
                pubkey,
                _toCredentials(msg.sender),
                withdrawalCredentials,
                amount,
                withdrawalCount++
            );
        }
    }

    /**
     * Transform an address into bytes for the credentials appending the 0x01 prefix.
     * @param addr The address to transform.
     * @return credentials The credentials.
     */
    function _toCredentials(address addr)
        private
        pure
        returns (bytes memory credentials)
    {
        // 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
        assembly ("memory-safe") {
            credentials := mload(0x40)
            mstore(credentials, 0x20)
            mstore(add(credentials, 0x20), or(addr, shl(248, 1)))
            mstore(0x40, add(credentials, 0x40))
        }
    }

    /// @dev Validates the deposit amount and sends the native asset to the zero address.
    function _deposit(uint64) internal virtual returns (uint64) {
        if (msg.value % 1 gwei != 0) {
            revert DepositNotMultipleOfGwei();
        }

        uint256 amountInGwei = msg.value / 1 gwei;
        if (amountInGwei > type(uint64).max) {
            revert DepositValueTooHigh();
        }

        _safeTransferETH(address(0), msg.value);

        return uint64(amountInGwei);
    }

    /**
     * @notice Safely transfers ETH to the given address.
     * @dev From the Solady library.
     * @param to The address to transfer the ETH to.
     * @param amount The amount of ETH to transfer.
     */
    function _safeTransferETH(address to, uint256 amount) private {
        /// @solidity memory-safe-assembly
        assembly {
            if iszero(
                call(gas(), to, amount, codesize(), 0x00, codesize(), 0x00)
            ) {
                mstore(0x00, 0xb12d13eb) // `ETHTransferFailed()`.
                revert(0x1c, 0x04)
            }
        }
    }
}
