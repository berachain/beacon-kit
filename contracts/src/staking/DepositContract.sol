// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import { IDepositContract } from "./IDepositContract.sol";
import { ERC165 } from "./IERC165.sol";

/**
 * @title DepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 * @dev This contract does not implement the deposit merkle tree.
 */
abstract contract DepositContract is IDepositContract, ERC165 {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The minimum amount of stake that can be deposited to prevent dust.
    /// @dev This is 32 ether in Gwei since our deposit contract denominates in Gwei. 32e9 * 1e9 = 32e18.
    uint64 internal constant MIN_DEPOSIT_AMOUNT_IN_GWEI = 32e9;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 internal constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 internal constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 internal constant CREDENTIALS_LENGTH = 32;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    uint64 public depositCount;

    /// @dev The operator address for each pubkey.
    mapping(bytes => address) private operatorByPubKey;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            VIEWS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc ERC165
    function supportsInterface(bytes4 interfaceId)
        external
        pure
        override
        returns (bool)
    {
        return interfaceId == type(ERC165).interfaceId
            || interfaceId == type(IDepositContract).interfaceId;
    }

    /// @inheritdoc IDepositContract
    function getOperator(bytes calldata pubkey) public view returns (address) {
        return operatorByPubKey[pubkey];
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        uint64 amount,
        bytes calldata signature,
        address operator
    )
        public
        payable
        virtual
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

        // Set operator on the first deposit.
        // zero `operatorByPubKey[pubkey]` means the pubkey is not registered.
        if (operatorByPubKey[pubkey] == address(0)) {
            if (operator == address(0)) {
                revert ZeroOperatorOnFirstDeposit();
            }
            operatorByPubKey[pubkey] = operator;
            emit OperatorUpdated(pubkey, operator, address(0));
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

    function updateOperator(
        bytes calldata pubkey,
        address newOperator
    )
        external
    {
        // cache the current operator
        address currentOperator = operatorByPubKey[pubkey];
        // This will also revert if the pubkey is not registered.
        if (msg.sender != currentOperator) {
            revert NotCurrentOperator();
        }
        // Revert if the new operator is zero address.
        if (newOperator == address(0)) {
            revert ZeroAddress();
        }
        // update the operator
        operatorByPubKey[pubkey] = newOperator;
        emit OperatorUpdated(pubkey, newOperator, currentOperator);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            INTERNAL                         */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

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
    function _safeTransferETH(address to, uint256 amount) internal {
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
