// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import { IBeaconDepositContract } from "./IBeaconDepositContract.sol";
import { Ownable } from "@solady/src/auth/Ownable.sol";

/**
 * @title BeaconDepositContract
 * @author Berachain Team
 * @notice A contract that handles deposits of stake.
 * @dev Its events are used by the beacon chain to manage the staking process.
 * @dev Its stake asset needs to be of 18 decimals to match the native asset.
 */
contract BeaconDepositContract is IBeaconDepositContract, Ownable {
    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                        CONSTANTS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev The minimum amount of stake that can be deposited to prevent dust.
    /// @dev This is 32 ether in Gwei since our deposit contract denominates in Gwei. 32e9 * 1e9 = 32e18.
    uint64 private constant MIN_DEPOSIT_AMOUNT_IN_GWEI = 32e9;

    /// @dev The length of the public key, PUBLIC_KEY_LENGTH bytes.
    uint8 private constant PUBLIC_KEY_LENGTH = 48;

    /// @dev The length of the signature, SIGNATURE_LENGTH bytes.
    uint8 private constant SIGNATURE_LENGTH = 96;

    /// @dev The length of the credentials, 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
    uint8 private constant CREDENTIALS_LENGTH = 32;

    uint256 private constant TWO_DAYS = 172_800; // 2 days

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           STORAGE                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev QueuedOperator is a struct that represents an operator address change request.
    struct QueuedOperator {
        address newOperator;
        uint256 queuedTimestamp;
    }

    /// @dev A flag to check if the contract has been initialized.
    bool private initialized = false;

    /// @dev depositCount represents the number of deposits that
    /// have been made to the contract.
    uint64 public depositCount;

    /// @dev queuedOperator is a mapping of public keys to operator change requests.
    mapping(bytes => QueuedOperator) private queuedOperator;

    /// @dev _pubkeyToOperator is a mapping of public keys to operators.
    /// @dev It is used in `POL` to control validator operation like setting cutting board, commission, etc.
    mapping(bytes => address) private _pubkeyToOperator;

    /// @dev depositAuth is a mapping of number of deposits an authorized address can make.
    mapping(address => uint64) public depositAuth;

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            WRITES                          */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @dev Initializes the owner of the contract.
    function initializeOwner(address _governance) external {
        require(!initialized, "Already initialized");
        _initializeOwner(_governance);
        initialized = true;
    }

    /// @dev Guard to prevent double initialization of owner.
    function _guardInitializeOwner() internal pure override returns (bool guard) {
        return true;
    }

    /// @inheritdoc IBeaconDepositContract
    function deposit(
        bytes calldata pubkey,
        bytes calldata credentials,
        bytes calldata signature
    )
        external
        payable
    {
        if (depositAuth[msg.sender] == 0) {
            revert UnauthorizedDeposit();
        }

        if (pubkey.length != PUBLIC_KEY_LENGTH) {
            revert InvalidPubKeyLength();
        }

        if (credentials.length != CREDENTIALS_LENGTH) {
            revert InvalidCredentialsLength();
        }

        if (signature.length != SIGNATURE_LENGTH) {
            revert InvalidSignatureLength();
        }

        // Update the pubkey to operator mapping if first deposit.
        if (_pubkeyToOperator[pubkey] == address(0)) {
            _pubkeyToOperator[pubkey] = msg.sender;
            emit OperatorSet(pubkey, msg.sender, address(0));
        }

        uint64 amountInGwei = _deposit();

        if (amountInGwei < MIN_DEPOSIT_AMOUNT_IN_GWEI) {
            revert InsufficientDeposit();
        }

        --depositAuth[msg.sender];

        unchecked {
            // slither-disable-next-line reentrancy-benign,reentrancy-events
            emit Deposit(
                pubkey, credentials, amountInGwei, signature, depositCount++
            );
        }
    }

    /// @inheritdoc IBeaconDepositContract
    function allowDeposit(
        address depositor,
        uint64 number
    )
        external
        onlyOwner
    {
        // If the number is non-zero, set it to zero to avoid front-running by depositors to make more deposits.
        if (depositAuth[depositor] != 0) {
            depositAuth[depositor] = 0;
            return;
        }
        depositAuth[depositor] = number;
    }

    /// @inheritdoc IBeaconDepositContract
    function requestOperatorChange(
        bytes calldata pubkey,
        address newOperator
    )
        external
    {
        // Only the operator can request a change.
        if (msg.sender != _pubkeyToOperator[pubkey]) {
            revert NotOperator();
        }
        QueuedOperator storage qQ = queuedOperator[pubkey];
        qQ.newOperator = newOperator;
        qQ.queuedTimestamp = block.timestamp;
        emit OperatorChangeQueued(pubkey, newOperator);
    }

    /// @inheritdoc IBeaconDepositContract
    function cancelOperatorChange(bytes calldata pubkey) external {
        // Only the operator can cancel the change.
        if (msg.sender != _pubkeyToOperator[pubkey]) {
            revert NotOperator();
        }
        delete queuedOperator[pubkey];
        emit OperatorChangeCancelled(pubkey);
    }

    /// @inheritdoc IBeaconDepositContract
    function acceptOperatorChange(bytes calldata pubkey) external {
        QueuedOperator storage qQ = queuedOperator[pubkey];
        (address newOperator, uint256 queuedTimestamp) =
            (qQ.newOperator, qQ.queuedTimestamp);

        // Only the new operator can accept the change.
        if (msg.sender != newOperator) {
            revert NotNewOperator();
        }
        // 2 days buffer to accept the change.
        if (queuedTimestamp + TWO_DAYS > block.timestamp) {
            revert NotEnoughTimePassed();
        }
        address oldOperator = _pubkeyToOperator[pubkey];
        _pubkeyToOperator[pubkey] = newOperator;
        delete queuedOperator[pubkey];
        emit OperatorSet(pubkey, newOperator, oldOperator);
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                            READS                           */
    /*.•°:°.´+˚.*°.˚:*.´•*.+°.•°:´*.´•*.•°.•°:°.´:•˚°.*°.˚:*.´+°.•*/

    /// @inheritdoc IBeaconDepositContract
    function getOperator(bytes calldata pubkey)
        external
        view
        returns (address)
    {
        return _pubkeyToOperator[pubkey];
    }

    /*´:°•.°+.*•´.*:˚.°*.˚•´.°:°•.°•.*•´.*:˚.°*.˚•´.°:°•.°+.*•´.*:*/
    /*                           INTERNAL                         */
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
