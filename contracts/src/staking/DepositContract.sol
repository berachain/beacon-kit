// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import { ERC165 } from "./IERC165.sol";
import { IDepositContract } from "./IDepositContract.sol";

/**
 * @title DepositContract
 * @author Berachain Team
 * @notice A contract that handles validators deposits.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev This contract does not implement the deposit merkle tree.
 */
contract DepositContract is IDepositContract, ERC165 {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The minimum amount of `BERA` to deposit.
    /// @dev This is 1 ether in Gwei since our deposit contract denominates in Gwei. 1e9 * 1e9 = 1e18.
    uint64 internal constant MIN_DEPOSIT_AMOUNT_IN_GWEI = 1e9;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 internal constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 internal constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 internal constant CREDENTIALS_LENGTH = 32;

    /// @dev 1 day in seconds.
    /// @dev This is the delay before a new operator can accept a change.
    uint96 private constant ONE_DAY = 86_400;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev QueuedOperator is a struct that represents an operator address change request.
    struct QueuedOperator {
        uint96 queuedTimestamp;
        address newOperator;
    }

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    /// @dev The index of the next deposit will use this value.
    uint64 public depositCount;

    /// @dev The hash tree root of the genesis deposits.
    /// @dev Should be set in deployment (predeploy state or constructor).
    // slither-disable-next-line constable-states
    bytes32 public genesisDepositsRoot;

    /// @dev The mapping of public keys to operator addresses.
    mapping(bytes => address) private _operatorByPubKey;

    /// @dev The mapping of public keys to operator change requests.
    mapping(bytes => QueuedOperator) public queuedOperator;

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
    function getOperator(bytes calldata pubkey)
        external
        view
        returns (address)
    {
        return _operatorByPubKey[pubkey];
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            DEPOSIT                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
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
        // zero `_operatorByPubKey[pubkey]` means the pubkey is not registered.
        if (_operatorByPubKey[pubkey] == address(0)) {
            if (operator == address(0)) {
                revert ZeroOperatorOnFirstDeposit();
            }
            _operatorByPubKey[pubkey] = operator;
            emit OperatorUpdated(pubkey, operator, address(0));
        }
        // If not the first deposit, operator address must be 0.
        // This prevents from the front-running of the first deposit to set the operator.
        else if (operator != address(0)) {
            revert OperatorAlreadySet();
        }

        uint64 amountInGwei = _deposit();

        if (amountInGwei < MIN_DEPOSIT_AMOUNT_IN_GWEI) {
            revert InsufficientDeposit();
        }

        // slither-disable-next-line reentrancy-benign,reentrancy-events
        emit Deposit(
            pubkey, credentials, amountInGwei, signature, depositCount++
        );
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        OPERATOR CHANGE                    */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IDepositContract
    function requestOperatorChange(
        bytes calldata pubkey,
        address newOperator
    )
        external
    {
        // Cache the current operator.
        address currentOperator = _operatorByPubKey[pubkey];
        // Only the current operator can request a change.
        // This will also revert if the pubkey is not registered.
        if (msg.sender != currentOperator) {
            revert NotOperator();
        }
        // Revert if the new operator is zero address.
        if (newOperator == address(0)) {
            revert ZeroAddress();
        }
        QueuedOperator storage qO = queuedOperator[pubkey];
        qO.newOperator = newOperator;
        qO.queuedTimestamp = uint96(block.timestamp);
        emit OperatorChangeQueued(
            pubkey, newOperator, currentOperator, block.timestamp
        );
    }

    /// @inheritdoc IDepositContract
    function cancelOperatorChange(bytes calldata pubkey) external {
        // Only the current operator can cancel the change.
        if (msg.sender != _operatorByPubKey[pubkey]) {
            revert NotOperator();
        }
        delete queuedOperator[pubkey];
        emit OperatorChangeCancelled(pubkey);
    }

    /// @inheritdoc IDepositContract
    function acceptOperatorChange(bytes calldata pubkey) external {
        QueuedOperator storage qO = queuedOperator[pubkey];
        (address newOperator, uint96 queuedTimestamp) =
            (qO.newOperator, qO.queuedTimestamp);

        // Only the new operator can accept the change.
        // This will revert if nothing is queued as newOperator will be zero address.
        if (msg.sender != newOperator) {
            revert NotNewOperator();
        }
        // Check if the queue delay has passed.
        if (queuedTimestamp + ONE_DAY > uint96(block.timestamp)) {
            revert NotEnoughTime();
        }
        // Cache the old operator.
        address oldOperator = _operatorByPubKey[pubkey];
        _operatorByPubKey[pubkey] = newOperator;
        delete queuedOperator[pubkey];
        emit OperatorUpdated(pubkey, newOperator, oldOperator);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            INTERNAL                         */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Validates the deposit amount and sends the native asset to the zero address.
    function _deposit() internal virtual returns (uint64) {
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
