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

import { IBeaconDepositContract } from "./IBeaconDepositContract.sol";
import { IStakeERC20 } from "./IStakeERC20.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits, withdrawals, and redirections of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 */
contract BeaconDepositContract is IBeaconDepositContract {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The address of the native asset as of EIP-7528.
    address public constant NATIVE_ASSET =
        0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE;

    /// @dev The address of the staking asset.
    /// @notice defaults to the native asset but can be changed at genesis at slot!!!.
    /// @notice to allow for the staking asset to be changed, we must store it in the contract storage.
    /// @notice storage layout:
    // | Name        | Type    | Slot | Offset | Bytes | Contract                                |
    // |-------------|---------|------|--------|-------|-----------------------------------------|
    // | STAKE_ASSET | address | 0    | 0      | 20    | src/staking/Deposit.sol:DepositContract |
    //
    // slither-disable-next-line constable-states
    address private STAKE_ASSET = NATIVE_ASSET;

    /// @dev The minimum amount of stake that can be deposited to prevent dust.
    /// @dev This is 32 ether in Gwei since our deposit contract denominates in Gwei. 32e9 * 1e9 = 32e18.
    uint64 private constant MIN_DEPOSIT_AMOUNT = 32e9;

    /// @dev The minimum amount of stake that can be redirected to prevent dust.
    /// leaving the buffer for their deposit to be slashed.
    uint256 private constant MIN_REDIRECT_AMOUNT = MIN_DEPOSIT_AMOUNT / 10;

    /// @dev The minimum amount of stake that can be withdrawn to prevent dust.
    /// leaving the buffer for their deposit to be slashed.
    uint256 private constant MINIMUM_WITHDRAWAL_AMOUNT = MIN_DEPOSIT_AMOUNT / 10;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 private constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 private constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 private constant CREDENTIALS_LENGTH = 32;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconDepositContract
    function deposit(
        bytes calldata validatorPubKey,
        bytes calldata stakingCredentials,
        uint64 amount,
        bytes calldata signature
    )
        external
        payable
    {
        if (validatorPubKey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (stakingCredentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (signature.length != SIGNATURE_LENGTH) {
            revert InvalidSignatureLength();
        }

        if (STAKE_ASSET == NATIVE_ASSET) {
            amount = _depositNative();
        } else {
            _depositERC20(amount);
        }

        // slither-disable-next-line reentrancy-events
        emit Deposit(validatorPubKey, stakingCredentials, amount, signature);
    }

    /// @inheritdoc IBeaconDepositContract
    function redirect(
        bytes calldata fromPubKey,
        bytes calldata toPubKey,
        uint64 amount
    )
        external
    {
        if (
            fromPubKey.length != PUBLIC_KEY_LENGTH
                || toPubKey.length != PUBLIC_KEY_LENGTH
        ) {
            revert InvalidPubKeyLength();
        }

        if (amount < MIN_REDIRECT_AMOUNT) {
            revert InsufficientRedirectAmount();
        }

        emit Redirect(fromPubKey, toPubKey, _toCredentials(msg.sender), amount);
    }

    /// @inheritdoc IBeaconDepositContract
    function withdraw(
        bytes calldata validatorPubKey,
        bytes calldata withdrawalCredentials,
        uint64 amount
    )
        external
    {
        if (validatorPubKey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (withdrawalCredentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (amount < MINIMUM_WITHDRAWAL_AMOUNT) {
            revert InsufficientWithdrawAmount();
        }

        emit Withdraw(validatorPubKey, withdrawalCredentials, amount);
    }

    /**
     * Transform an address into bytes for the credentials appending the 0x01 prefix.
     * @param addr The address to transform.
     * @return The credentials.
     */
    function _toCredentials(address addr) private pure returns (bytes memory) {
        // 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
        return abi.encodePacked(bytes1(0x01), bytes11(0x0), addr);
    }

    /**
     * @notice Validates the deposit amount and sends the native asset to the zero address.
     */
    function _depositNative() private returns (uint64) {
        if (msg.value > type(uint64).max) {
            revert DepositValueTooHigh();
        }

        if (msg.value < MIN_DEPOSIT_AMOUNT) {
            revert InsufficientDeposit();
        }

        if (msg.value % 1 gwei != 0) {
            revert DepositNotMultipleOfGwei();
        }

        _safeTransferETH(address(0), msg.value);

        // Safe since we have already checked that the value is less than uint64.max.
        return uint64(msg.value);
    }

    /*
     * @notice Validates the deposit amount and burns the staking asset from the sender.
     * @param amount The amount of stake to deposit.
     */
    function _depositERC20(uint64 amount) private {
        IStakeERC20(STAKE_ASSET).burn(msg.sender, amount);

        if (amount < MIN_DEPOSIT_AMOUNT) {
            revert InsufficientDeposit();
        }

        if (amount % 1 gwei != 0) {
            revert DepositNotMultipleOfGwei();
        }
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
