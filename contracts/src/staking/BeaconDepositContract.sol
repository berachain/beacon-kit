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
    //
    // slither-disable-next-line locked-ether
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
    //
    // slither-disable-next-line locked-ether
    function withdraw(
        bytes calldata validatorPubkey,
        uint64 amount
    )
        external
        payable
    {
        // TODO: Properly Handle Token Logic.

<<<<<<< HEAD
<<<<<<< HEAD
        emit Withdrawal(validatorPubkey, abi.encodePacked(msg.sender), amount);
=======
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

        emit Withdrawal(validatorPubKey, withdrawalCredentials, amount);
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
>>>>>>> 9d0aa54e (feat(types): fix types for kzg n stuff)
=======
        emit Withdrawal(validatorPubkey, abi.encodePacked(msg.sender), amount);
>>>>>>> 06711e04 (Revert: 6d8a5a0dddb45a46f13f4f746efabe0f73ae3394)
    }
}
